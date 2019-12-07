package init

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gcfg.v1"
	"object-storage-go/data-server/constant"
	"os"
)

func InitLog() error {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile(constant.LOG_FILEPATH, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("init log failed")
		return fmt.Errorf("init log failed")
	}
	log.SetOutput(file)
	log.SetLevel(log.DebugLevel)
	return nil
}

func InitConfig() error {
	gcfg.ReadFileInto()
}
