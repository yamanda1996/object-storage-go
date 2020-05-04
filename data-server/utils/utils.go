package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"object-storage-go/data-server/model"
	"object-storage-go/data-server/store"
	"os"
	"path/filepath"
	"strconv"
)

var DataServerUrl 	string
var DataServerPrefix string

// func init() {
// 	err := InitLog()
// 	if err != nil {
// 		fmt.Println("init log failed")
// 		os.Exit(1)
// 	}

// 	err = InitConfig()
// 	if err != nil {
// 		log.Error("init config file failed")
// 		os.Exit(1)
// 	}

// 	DataServerPrefix = model.Config.DataServerConfig.DataServerEtcdPrefix + model.Config.DataServerConfig.DataServerIndex
// 	DataServerUrl = model.Config.DataServerConfig.DataServerAddress + ":" +
// 		model.Config.DataServerConfig.DataServerPort
// 	err = RegisterToEtcdServer(DataServerPrefix, DataServerUrl)
// 	if err != nil {
// 		log.Error("register data server to etcd server failed")
// 		os.Exit(1)
// 	}
// }

func InitLog() error {
	log.SetFormatter(&log.TextFormatter{})
	file, err := os.OpenFile(model.Conf.LogConfig.LogFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("init log failed")
		return fmt.Errorf("init log failed")
	}
	log.SetOutput(file)
	log.SetLevel(log.DebugLevel)
	return nil
}

func GetBucketDir(numBucket, bucketID int) string {
	if numBucket == 1 {
		return ""
	} else if numBucket == 16 {
		return fmt.Sprintf("%x", bucketID)
	} else if numBucket == 256 {
		return fmt.Sprintf("%x/%x", bucketID/16, bucketID%16)
	}
	panic(fmt.Sprintf("wrong numBucket: %d", numBucket))
}

func GetBucketPath(bucketID int) string {
	return filepath.Join(model.Conf.DataServerConfig.StorageRoot,
		GetBucketDir(model.Conf.DataServerConfig.BucketNumber, bucketID))
}

func GenDataPath(home string, chunkID int) string {
	return fmt.Sprintf("%s/%03d.data", home, chunkID)
}

func GetIndexPath(home string, chunkID, splitID int, suffix string) string {
	return fmt.Sprintf("%s/%s.%s.idx.%s", home, idToStr(chunkID), idToStr(splitID), suffix)
}

func idToStr(id int) string {
	if id < 0 {
		return "*"
	}
	return fmt.Sprintf("%03d", id)
}

func ParseIDFromPath(path string) (id store.HintID, ok bool) {
	return parseIDFromName(filepath.Base(path))
}

func parseIDFromName(name string) (id store.HintID, ok bool) {
	ck, err1 := parseChunkIDFromName(name)
	sp, err2 := parseSplitIDFromName(name)
	if err1 == nil && err2 == nil && ck < MAX_NUM_CHUNK && sp < MAX_NUM_SPLIT {
		return store.HintID{ck, sp}, true
	}
	return
}

func parseChunkIDFromName(name string) (int, error) {
	return strconv.Atoi(name[:3])
}

func parseSplitIDFromName(name string) (int, error) {
	return strconv.Atoi(name[4:7])
}

func Remove(path string) error {
	return os.Remove(path)
}

// func InitConfig() error {

// 	err := gcfg.ReadFileInto(&model.Config, constant.DATA_SERVER_CONFIG_FILEPATH)
// 	if err != nil {
// 		log.Errorf("read config file [%s] failed", constant.DATA_SERVER_CONFIG_FILEPATH)
// 		os.Exit(1)
// 	}
// 	return nil
// }

// func RegisterToEtcdServer(k,v string) error {

// 	endpoints := []string{model.Config.EtcdServerConfig.EtcdServerAddress + ":" +
// 		model.Config.EtcdServerConfig.EtcdServerPort}

// 	var err error
// 	Register, err = etcd.NewServiceRegister(endpoints, common_constant.ETCD_DIAL_TIMEOUT)
// 	if err != nil {
// 		log.Error(err.Error())
// 		return err
// 	}

// 	err = Register.RegisterService(k, v)
// 	if err != nil {
// 		log.Error(err.Error())
// 		return err
// 	}
// 	return nil
// }