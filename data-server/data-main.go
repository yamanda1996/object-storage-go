package data_server

import (
	"fmt"
	"object-storage-go/data-server/cron"
	"object-storage-go/data-server/utils"
)

func main() {

	err := utils.InitLog()
	if err != nil {
		fmt.Println("Init log failed")
		return
	}

	cron.Cron()

}


