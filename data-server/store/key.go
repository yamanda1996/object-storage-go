package store

type KeyInfo struct {
	KeyHash   uint64
	KeyIsPath bool
	Key       []byte
	StringKey string
	KeyPos
}