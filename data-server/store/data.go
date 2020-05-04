package store

import (
	"fmt"
	"os"
	"sync"
	"time"

	logger "github.com/sirupsen/logrus"
	"object-storage-go/data-server/utils"
)

type dataStore struct {
	bucketID int
	home     string

	sync.Mutex
	flushLock sync.Mutex

	oldHead int // old tail == 0
	newHead int
	newTail int

	chunks        [utils.MAX_NUM_CHUNK]dataChunk
	wbufSize      uint32
	lastFlushTime time.Time
}

func NewDataStore(bucketID int, home string) *dataStore {
	ds := new(dataStore)
	ds.bucketID = bucketID
	ds.home = home
	for i := 0; i < utils.MAX_NUM_CHUNK; i++ {
		ds.chunks[i].chunkid = i
		ds.chunks[i].path = utils.GenDataPath(ds.home, i)
	}
	return ds
}

func (ds *dataStore) ListFiles() (max int, err error) {
	max = -1
	for i := 0; i < utils.MAX_NUM_CHUNK; i++ {
		path := utils.GenDataPath(ds.home, i)
		st, e := os.Stat(path)
		if e != nil {
			pe := e.(*os.PathError)
			if "no such file or directory" == pe.Err.Error() {
				ds.chunks[i].size = 0
			} else {
				logger.Errorf(pe.Err.Error())
				err = pe
				return
			}
		} else {
			sz := uint32(st.Size())
			if (sz & 0xff) != 0 {
				err = fmt.Errorf("file not 256 aligned, size 0x%x: %s ", sz, path)
				return
			}
			ds.chunks[i].size = sz
			ds.chunks[i].writingHead = sz
			max = i
		}
	}
	ds.newHead = max + 1
	return
}