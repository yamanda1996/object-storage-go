package store

import "sync"

type dataChunk struct {
	sync.Mutex

	chunkid int
	path    string
	size    uint32

	writingHead uint32
	wbuf        []*WriteRecord

	rewriting bool
	gcbufsize uint32
	gcWriter  *DataStreamWriter
}
