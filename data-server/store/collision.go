package store

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"sync"

	logger "github.com/sirupsen/logrus"
)

type CollisionTable struct {
	sync.Mutex `yaml:"-"`
	HintID
	Items map[uint64]map[string]HintItem
}

func NewCollisionTable() *CollisionTable {
	t := &CollisionTable{}
	t.Items = make(map[uint64]map[string]HintItem)
	return t
}

func (table *CollisionTable) load(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			logger.Errorf("read yaml failed %s: %s", path, err.Error())
		}
		return
	}
	table.Lock()
	if err := yaml.Unmarshal(content, table); err != nil {
		logger.Errorf("unmarshal yaml faild %s %s", path, err.Error())
	}
	table.Unlock()
}
