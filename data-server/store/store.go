package store

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"object-storage-go/data-server/model"
	"object-storage-go/data-server/utils"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

var mergeChan chan int

type HStore struct {
	buckets   []*Bucket
	gcMgr     *GCMgr
	htree     *HTree
	htreeLock sync.Mutex
}

func NewHStore() (store *HStore, err error) {
	if err := os.MkdirAll(model.Conf.DataServerConfig.StorageRoot, os.ModePerm); err != nil {
		logger.Fatalf("fail to init home %s", model.Conf.DataServerConfig.StorageRoot)
	}
	mergeChan = nil
	//cmem.DBRL.ResetAll()
	store = new(HStore)
	store.gcMgr = &GCMgr{stat: make(map[*Bucket]*GCState)}
	store.buckets = make([]*Bucket, model.Conf.DataServerConfig.BucketNumber)
	for i := 0; i < model.Conf.DataServerConfig.BucketNumber; i++ {
		store.buckets[i] = &Bucket{}
		store.buckets[i].ID = i
	}
	err = store.scanBuckets()
	if err != nil {
		return
	}

	for i := 0; i < model.Conf.DataServerConfig.BucketNumber; i++ {
		if store.buckets[i].BucketStat.State < utils.BUCKET_STAT_NOT_EMPTY {
			err = store.allocBucket(i)
			if err != nil {
				return
			}

		}
		store.buckets[i].State = utils.BUCKET_STAT_READY
	}

	var n int32
	var wg = sync.WaitGroup{}
	wg.Add(model.Conf.DataServerConfig.BucketNumber)
	errs := make(chan error, model.Conf.DataServerConfig.BucketNumber)
	for i := 0; i < model.Conf.DataServerConfig.BucketNumber; i++ {
		go func(id int) {
			defer wg.Done()
			bkt := store.buckets[id]
			if store.buckets[id].State > utils.BUCKET_STAT_EMPTY {
				err = bkt.open(id, utils.GetBucketPath(id))
				if err != nil {
					logger.Errorf("Error in bkt open %s", err.Error())
					errs <- err
				} else {
					atomic.AddInt32(&n, 1)
				}
			}
		}(i)
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		if e != nil {
			err = e
			return
		}
	}
	if model.Conf.DataServerConfig.TreeDepth > 0 {
		store.htree = NewHTree(0, 0, model.Conf.DataServerConfig.TreeDepth + 1)
	}
	return
}

func (store *HStore) allocBucket(bucketID int) (err error) {
	dirpath := utils.GetBucketPath(bucketID)
	if _, err = os.Stat(dirpath); err != nil {
		err = os.MkdirAll(dirpath, 0755)
	}
	logger.Infof("allocBucket %s", dirpath)
	return
}

func (s *HStore) scanBuckets() (err error) {
	for id := 0; id < model.Conf.DataServerConfig.BucketNumber; id++ {
		path := utils.GetBucketPath(id)
		fi, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			logger.Infof("%s", err.Error())
			return err
		}
		if !fi.IsDir() {
			err = fmt.Errorf("%s is not dir", path)
			logger.Errorf("%s", err.Error())
			return err
		}

		datas, err := filepath.Glob(filepath.Join(path, "*.data"))
		if err != nil {
			logger.Errorf("%s", err.Error())
			return err
		}
		if len(datas) == 0 {
			if model.Conf.DataServerConfig.BucketNumber > 1 {
				logger.Warnf("remove empty bucket dir %s", path)
				if err = os.RemoveAll(path); err != nil {
					logger.Errorf("fail to delete empty bucket %s", path)
				}
			}
		} else {
			logger.Infof("found bucket %x", id)
			s.buckets[id].State = utils.BUCKET_STAT_NOT_EMPTY
		}
	}
	return nil
}