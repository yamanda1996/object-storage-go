package object

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"object-storage-go/api-server/es"
	"object-storage-go/api-server/heartbeat"
	"object-storage-go/api-server/locate"
	"object-storage-go/api-server/objectstream"
	"object-storage-go/api-server/utils"
	"strconv"
)

func UploadFile(context *gin.Context)  {
	name := context.Param("filename")
	if len(name) <= 0 {
		context.String(http.StatusNotFound, "request parameter filename [%s] invalid", name)
	}

	hash := utils.GetFromHeader(context, "Digest")
	if hash == nil {
		log.Errorf("hash not found")
		context.String(http.StatusBadRequest, "hash not found in header")
	}

	size := utils.GetFromHeader(context, "Size")
	if size == nil {
		log.Errorf("size not found")
		context.String(http.StatusBadRequest, "size not found in header")
	}

	file, _ := context.FormFile("file")
	r, _ := file.Open()
	status, err := storeObject(r, url.PathEscape(hash[0]))
	if err != nil {
		log.Errorf("store [%s] file failed", name)
		context.String(http.StatusInternalServerError, "store file failed")
	}

	s, _ := strconv.ParseInt(size[0], 10, 64)
	err = es.AddVersion(name, hash[0], s)
	if err != nil {
		log.Errorf("add [%s] version failed", name)
		context.String(http.StatusInternalServerError, "add version failed")
	}

	context.String(status, "upload file success")
}

func DownloadFile(context *gin.Context)  {
	filename := context.Param("filename")
	if len(filename) <= 0 {
		context.String(http.StatusNotFound, "request parameter filename [%s] invalid", filename)
	}

	version := context.Query("version")
	var hash string
	if len(version) == 0 {
		metadata, err := es.SearchLatestVersion(filename)
		if err != nil {
			log.Infof("get [%s] latest version failed", filename)
			context.String(http.StatusInternalServerError, "download failed")
		}
		hash = metadata.Hash
	}
	i, _ := strconv.Atoi(version)
	metadata, err := es.GetMetadata(filename, i)
	if err != nil {
		log.Errorf("get [%s] metadata for version [%d] failed", filename, i)
		context.String(http.StatusInternalServerError, "download failed")
	}
	hash = metadata.Hash

	stream, err := getStream(hash)
	if err != nil {
		context.String(http.StatusNotFound, "file [%s] not found", filename)
	}

	b, err := ioutil.ReadAll(stream)
	if err != nil {
		log.Errorf("write to bytes failed")
		context.String(http.StatusInternalServerError, "internal server error")
	}
	context.Data(http.StatusOK, "fileType", b)
}

func DeleteFile(context *gin.Context)  {
	filename := context.Param("filename")
	if len(filename) <= 0 {
		context.String(http.StatusNotFound, "request parameter filename [%s] invalid", filename)
	}
	metadata, err := es.SearchLatestVersion(filename)
	if err != nil {
		log.Errorf("search [%s] latest version failed", filename)
		context.String(http.StatusInternalServerError, "search latest version failed")
	}
	err = es.PutMetadata(filename, metadata.Version+1, 0, "")
	if err != nil {
		log.Errorf("logical delete file [%s] failed", filename)
		context.String(http.StatusInternalServerError, "logical delete file failed")
	}
	context.String(http.StatusOK, "delete file success")
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

