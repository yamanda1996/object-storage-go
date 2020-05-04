package store

import (
	"sync"
	"time"
)

type GCMgr struct {
	mu   sync.RWMutex
	stat map[*Bucket]*GCState // map[bucketID]*GCState
}

type GCState struct {
	BeginTS time.Time
	EndTS   time.Time

	// Begin and End are chunckIDs, they determine the range of GC.
	Begin int
	End   int

	// Src and Dst are chunkIDs, they are tmp variables used in gc process.
	Src int
	Dst int

	// For beansdbadmin check status.
	Running bool

	Err        error
	CancelFlag bool
	// sum
	GCFileState
}