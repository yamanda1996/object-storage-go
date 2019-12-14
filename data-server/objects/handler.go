package objects

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"object-storage-go/data-server/model"
	"os"
)

func GetObject(context *gin.Context)  {
	filename := context.Param("filename")
	log.Infof("start to get object file [%s]", filename)

	f, err := ioutil.ReadFile(model.Config.DataServerConfig.StorageRoot + "/objects/" + filename)
	if err != nil {
		log.Errorf("read file [%s] failed", filename)
		context.String(http.StatusNotFound, "get file [%s] failed", filename)
	}
	context.Data(http.StatusOK, "fileType", f)
}

func PutObject(context *gin.Context)  {
	name := context.Param("filename")
	log.Debugf("upload file [%s]", name)
	file, err := os.Create(model.Config.DataServerConfig.StorageRoot + "/objects/" + name)
	if err != nil {
		log.Errorf("data server create file failed")
		context.String(http.StatusServiceUnavailable, "create file [%s] failed", name)
	}

	defer file.Close()

	io.Copy(file, context.Request.Body)
}
