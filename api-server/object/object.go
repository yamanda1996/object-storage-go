package object

import (
	"github.com/gin-gonic/gin"
	"object-storage-go/api-server/heartbeat"
	"object-storage-go/api-server/locate"
	"object-storage-go/api-server/objectstream"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request)  {
	m := r.Method

	if m == http.MethodGet {
		get(w,r)
		return
	}

	if m == http.MethodPut {
		put(w,r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func get(w http.ResponseWriter, r *http.Request)  {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, err := getStream(object)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, stream)
}

func put(w http.ResponseWriter, r *http.Request)  {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	respCode, err := storeObject(r.Body, object)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(respCode)
}

func storeObject(r io.Reader, object string) (int ,error) {
	putStream, err := putStream(object)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	io.Copy(putStream, r)

	err = putStream.Close()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func putStream(object string) (*objectstream.PutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataserver")
	}

	return objectstream.NewPutStream(server, object), nil
}

func getStream(object string) (*objectstream.GetStream, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate failed", object)
	}
	return objectstream.NewGetStream(server, object)
}

func GetFile(conetxt *gin.Context)  {

}

