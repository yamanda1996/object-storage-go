package main

import (
	"object-storage-go/data-serverage-go/data-server/heartbeat"
	"object-storage-go/data-serverage-go/data-server/locate"
	"object-storage-go/data-serverage-go/data-server/objects"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
	"os"
)

func main() {

	go heartbeat.HeartBeat()
	go locate.StartLocate()

	http.HandleFunc("/object/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))

}


