package store

import (
	"io"
	"object-storage-go/data-server/model"
	"object-storage-go/data-server/utils"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	logger "github.com/sirupsen/logrus"
)

type BucketStat struct {
	// pre open init
	State int

	// init in open
	ID   int
	Home string

	TreeID      HintID
	NextGCChunk int
}

type BucketInfo struct {
	BucketStat

	// tmp
	Pos             Position
	LastGC          *GCState
	HintState       int
	MaxDumpedHintID HintID
	DU              int64
	NumSameVhash    int64
	SizeSameVhash   int64
	SizeVhashKey    string
	NumSet          int64
	NumGet          int64
}

type Bucket struct {
	writeLock sync.Mutex
	BucketInfo

	htree     *HTree
	hints     *hintMgr
	datas     *dataStore
	GCHistory []GCState
}

func (bkt *Bucket) open(bucketID int, home string) (err error) {
	st := time.Now()
	// load HTree
	bkt.ID = bucketID
	bkt.Home = home
	bkt.datas = NewDataStore(bucketID, home)
	bkt.hints = NewHintMgr(bucketID, home)
	bkt.hints.loadCollisions()
	htree := NewHTree(model.Conf.DataServerConfig.TreeDepth, bucketID, model.Conf.DataServerConfig.TreeHeight)

	bkt.TreeID = HintID{0, -1}

	maxdata, err := bkt.datas.ListFiles()
	if err != nil {
		return err
	}
	htrees, ids := bkt.getAllIndex(utils.HTREE_SUFFIX)
	for i := len(htrees) - 1; i >= 0; i-- {
		treepath := htrees[i]
		id := ids[i]
		if id.Chunk > maxdata {
			logger.Errorf("htree beyond data: htree=%s, maxdata=%d", treepath, maxdata)
			utils.Remove(treepath)
		} else {
			if bkt.TreeID.isLarger(id.Chunk, id.Split) {
				err := htree.load(treepath)
				if err != nil {
					bkt.TreeID = HintID{0, -1}
					htree = NewHTree(model.Conf.DataServerConfig.TreeDepth, bucketID,
						model.Conf.DataServerConfig.TreeHeight)
				} else {
					bkt.TreeID = id
				}
			} else {
				logger.Errorf("found old htree: htree=%s, currenct_htree_id=%v", treepath, bkt.TreeID)
				utils.Remove(treepath)
			}
		}
	}
	bkt.htree = htree

	bkt.hints.maxDumpedHintID = bkt.TreeID
	for i := bkt.TreeID.Chunk; i < utils.MAX_NUM_CHUNK; i++ {
		startsp := 0
		if i == bkt.TreeID.Chunk {
			startsp = bkt.TreeID.Split + 1
		}
		e := bkt.checkHintWithData(i)
		if e != nil {
			err = e
			logger.Fatalf("fail to start for bad data, bkt:%02x, chuck:%d, err: %s", bkt.ID, i, e.Error())
		}
		splits := bkt.hints.chunks[i].splits
		numhintfile := len(splits) - 1
		if startsp >= numhintfile { // rebuilt
			continue
		}
		for j, sp := range splits[:numhintfile] {
			_, e := bkt.updateHtreeFromHint(i, sp.file.path)
			if e != nil {
				err = e
				return
			}
			bkt.hints.maxDumpedHintID = HintID{i, startsp + j}
		}
	}
	go func() {
		for i := 0; i < bkt.TreeID.Chunk; i++ {
			bkt.checkHintWithData(i)
		}
	}()

	if bkt.checkForDump(Conf.TreeDump) {
		bkt.dumpHtree()
	}

	bkt.loadGCHistroy()
	logger.Infof("bucket %x opened, max rss = %d, use time %s",
		bucketID, utils.GetMaxRSS(), time.Since(st))
	return nil
}

func (bkt *Bucket) loadGCHistroy() (err error) {
	fd, err := os.Open(bkt.getGCHistoryPath())
	if err != nil {
		logger.Infof("%v", err)
		return
	}
	defer fd.Close()
	buf := make([]byte, 10)
	n, e := fd.Read(buf)
	if e != nil && e != io.EOF {
		err = e
		logger.Errorf("%v", err)
		return
	}
	s := string(buf[:n])
	s = strings.TrimSpace(s)
	n, err = strconv.Atoi(s)
	if err == nil {
		bkt.NextGCChunk = n
		logger.Infof("bucket %d load nextgc %d", bkt.ID, n)
	}
	return
}

func (bkt *Bucket) dumpHtree() {
	hintID := bkt.hints.maxDumpedHintID
	if bkt.TreeID.isLarger(hintID.Chunk, hintID.Split) {
		bkt.removeHtree()
		bkt.TreeID = hintID
		bkt.htree.dump(bkt.getHtreePath(bkt.TreeID.Chunk, bkt.TreeID.Split))
	}
}

func (bkt *Bucket) checkForDump(dumpthreshold int) bool {
	maxdata, err := bkt.datas.ListFiles()
	if err != nil {
		return false
	}
	logger.Infof("maxdata %d", maxdata)
	if maxdata < 0 {
		return false
	}
	htrees, ids := bkt.getAllIndex(HTREE_SUFFIX)
	for i := len(htrees) - 1; i >= 0; i-- {
		id := ids[i]
		if maxdata > id.Chunk+dumpthreshold {
			return false
		}
	}

	if len(ids) > 0 {
		return false
	}
	return true
}

func (bkt *Bucket) updateHtreeFromHint(chunkID int, path string) (maxoffset uint32, err error) {
	logger.Infof("updateHtreeFromHint chunk %d, %s", chunkID, path)
	meta := Meta{}
	tree := bkt.htree
	var pos Position
	r := newHintFileReader(path, chunkID, 1<<20)
	r.open()
	maxoffset = r.datasize
	defer r.close()
	for {
		item, e := r.next()
		if e != nil {
			err = e
			return
		}
		if item == nil {
			return
		}
		ki := NewKeyInfoFromBytes([]byte(item.Key), item.Keyhash, false)
		ki.Prepare()
		meta.ValueHash = item.Vhash
		meta.Ver = item.Ver
		pos.Offset = item.Pos.Offset
		if item.Ver > 0 {
			pos.ChunkID = chunkID
			tree.set(ki, &meta, pos)
		} else {
			pos.ChunkID = -1
			tree.remove(ki, pos)
		}
	}
	return
}

func (bkt *Bucket) checkHintWithData(chunkID int) (err error) {
	size := bkt.datas.chunks[chunkID].size
	if size == 0 {
		bkt.hints.RemoveHintfilesByChunk(chunkID)
		return
	}
	hintDataSize := bkt.hints.loadHintsByChunk(chunkID)
	if hintDataSize < size {
		err = bkt.buildHintFromData(chunkID, hintDataSize)
	}
	return
}


func (bkt *Bucket) getAllIndex(suffix string) (paths []string, ids []HintID) {
	pattern := utils.GetIndexPath(bkt.Home, -1, -1, suffix)
	paths0, _ := filepath.Glob(pattern)
	sort.Sort(sort.StringSlice(paths0))
	for _, p := range paths0 {
		id, ok := utils.ParseIDFromPath(p)
		if !ok {
			logger.Errorf("find index file with wrong name %s", p)
		} else {
			paths = append(paths, p)
			ids = append(ids, id)
		}
	}
	return
}