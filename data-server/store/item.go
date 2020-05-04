package store

type HintItemMeta struct {
	Keyhash uint64
	Pos     Position
	Ver     int32
	Vhash   uint16
}

type HTreeItem HintItemMeta

type HintItem struct {
	HintItemMeta
	Key string
}
