package store

import (
	"fmt"
	"sync"

	"object-storage-go/data-server/utils"
)

type hintMgr struct {
	bucketID int
	home     string

	sync.Mutex // protect maxChunkID
	maxChunkID int

	chunks [utils.MAX_NUM_CHUNK]*hintChunk

	maxDumpedHintID HintID

	dumpLock           sync.Mutex
	mergeLock          sync.Mutex
	maxDumpableChunkID int
	merged             *hintFileIndex
	state              int

	collisions *CollisionTable
}

type HintID struct {
	Chunk int
	Split int
}

func (id *HintID) isLarger(ck, sp int) bool {
	return (ck > id.Chunk) || (ck == id.Chunk && sp >= id.Split)
}

func (id *HintID) setIfLarger(ck, sp int) (larger bool) {
	larger = id.isLarger(ck, sp)
	if larger {
		id.Chunk = ck
		id.Split = sp
	}
	return
}

func NewHintMgr(bucketID int, home string) *hintMgr {
	hm := &hintMgr{bucketID: bucketID, home: home}
	for i := 0; i < utils.MAX_NUM_CHUNK; i++ {
		hm.chunks[i] = newHintChunk(i)
	}
	hm.maxDumpableChunkID = utils.MAX_NUM_CHUNK - 1

	hm.collisions = NewCollisionTable()
	return hm
}

type hintChunk struct {
	sync.Mutex
	id       int
	fileLock sync.RWMutex
	splits   []*hintSplit

	// set to 0 : 1. loaded 2. before merge
	lastTS int64
}

func newHintChunk(id int) *hintChunk {
	ck := &hintChunk{id: id}
	ck.rotate()
	return ck
}

func (h *hintMgr) getCollisionPath() string {
	return fmt.Sprintf("%s/collision.yaml", h.home)
}

func (h *hintMgr) loadCollisions() {
	h.collisions.load(h.getCollisionPath())
}
