package store

import "sync"

type HTree struct {
	sync.Mutex

	/*
	 *            Root of hstore.tree
	 *               /  | ... \
	 *              /   |      \
	 * depth -->  ht1  ht2     htN     # root of bucket trees (also are the leafs of hstore tree)
	 *             ^
	 *             |
	 *            pos
	 *
	 * #bucket (number of buckets) = 16 ^ bucket_tree.depth
	 */

	// depth is level (0-based) of root Node (of this htree) in hstore.tree
	depth int

	// bucketID is position (offset) of this htree in the list of htrees at same level.
	bucketID int

	/* runtime */

	// level[0][0] is root of a htree,
	// levels[i] is a list of nodes at same level `i` of htree,
	// Node stores the summary info of its childs.
	// Height of htree = len(levels)
	levels [][]Node

	// leafs is the place to store key related info (e.g. keyhash, version, vhash etc.)
	leafs []SliceHeader

	// tmp, to avoid alloc
	ni NodeInfo
}