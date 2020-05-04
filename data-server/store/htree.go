package store

import (
	"bufio"
	"encoding/binary"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"io"
	"object-storage-go/data-server/model"
	"object-storage-go/data-server/utils"
	"os"
	"sync"
)

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

type Node struct {
	// count is the number of keys (with version > 0) under this node.
	count uint32

	// hash is the summary of it's child nodes.
	hash uint16

	// isHashUpdated is true iff the hash value of node is updated.
	isHashUpdated bool
}

type NodeInfo struct {
	node   *Node
	level  int
	offset int
	path   []int
}

func NewHTree(depth, bucketID, height int) *HTree {
	if depth + height > utils.MAX_DEPTH {
		panic("HTree too high")
	}
	tree := new(HTree)
	tree.depth = depth
	tree.bucketID = bucketID
	tree.levels = make([][]Node, height)
	size := 1
	for i := 0; i < height; i++ {
		tree.levels[i] = make([]Node, size)
		size *= 16
	}
	size /= 16
	leafnodes := tree.levels[height - 1]
	for i := 0; i < size; i++ {
		leafnodes[i].isHashUpdated = true
	}
	tree.leafs = make([]SliceHeader, size)
	return tree
}

func (tree *HTree) load(path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		logger.Errorf("fail to load htree %s", err.Error())
		return
	}
	defer f.Close()
	logger.Infof("loading htree %s", path)
	reader := bufio.NewReader(f)
	buf := make([]byte, 6)
	leafnodes := tree.levels[model.Conf.DataServerConfig.TreeHeight - 1]
	size := len(leafnodes)
	for i := 0; i < size; i++ {
		if _, err = io.ReadFull(reader, buf); err != nil {
			logger.Errorf("load htree err %s %v", path, err)
			return
		}
		leafnodes[i].count = binary.LittleEndian.Uint32(buf[0:4])
		leafnodes[i].hash = binary.LittleEndian.Uint16(buf[4:6])
	}
	for i := 0; i < size; i++ {
		if _, err = io.ReadFull(reader, buf[:4]); err != nil {
			logger.Errorf("load htree err %s %v", path, err)
			return
		}
		l := int(binary.LittleEndian.Uint32(buf[:4]))
		if l > 0 {
			tree.leafs[i].enlarge(int(l))
			if _, err = io.ReadFull(reader, tree.leafs[i].ToBytes()); err != nil {
				logger.Errorf("load htree err %s %v", path, err)
				return
			}
		}
	}

	tree.ListTop()
	return nil
}

func (tree *HTree) ListTop() {
	path := fmt.Sprintf("%x", tree.bucketID)
	ki := &KeyInfo{
		StringKey: path,
		Key:       []byte(path),
		KeyIsPath: true}
	ki.Prepare()
	data, _ := tree.ListDir(ki)
	logger.Infof("listing %s:\n%s", path, string(data))
	//items, nodes := tree.listDir(ki)
	//logger.Infof("%s %#v %#v", path, items, nodes)
}