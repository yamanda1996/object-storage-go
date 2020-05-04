package store

import "sync"

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

