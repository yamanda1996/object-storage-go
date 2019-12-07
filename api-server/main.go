package main

import (
	"object-storage-go/api-serverrage-go/api-server/heartbeat"
	"object-storage-go/api-serverrage-go/api-server/locate"
	"object-storage-go/api-serverrage-go/api-server/object"
	"log"
	"net/http"
	"os"
)

func main()  {
	go heartbeat.ListenHeartBeat()
	http.HandleFunc("/objects/", object.Handler)
	http.HandleFunc("/locate/", locate.Handler)

	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
