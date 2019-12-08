package object

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"object-storage-go/api-server/heartbeat"
	"object-storage-go/api-server/locate"
	"object-storage-go/api-server/objectstream"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func UploadFile(context *gin.Context)  {
	name := context.Param("filename")
	if len(name) <= 0 {
		context.String(http.StatusNotFound, "request parameter filename [%s] invalid", name)
	}

	file, _ := context.FormFile("file")
	r, _ := file.Open()
	status, err := storeObject(r, name)
	context.String(status, err.Error())
}

func GetFile(conetxt *gin.Context)  {

}

func storeObject(r io.Reader, name string) (int, error) {
	stream, err := putStream(name)
	if err != nil {
		log.Errorf("create put stream failed")
		return http.StatusServiceUnavailable, err
	}
	io.Copy(stream, r)
	err = stream.Close()
	if err != nil {
		log.Errorf("closer put stream failed")
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func putStream(object string) (*objectstream.PutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		log.Errorf("select data server failed")
		return nil, fmt.Errorf("cannot find any dataserver")
	}
	log.Debugf("select data server [%s]", server)
	return objectstream.NewPutStream(server, object), nil
}

func getStream(object string) (*objectstream.GetStream, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate failed", object)
	}
	return objectstream.NewGetStream(server, object)
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




