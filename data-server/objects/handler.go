package objects

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"object-storage-go/data-server/model"
	"os"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request)  {

	m := r.Method

	if m == http.MethodPut {
		put(w,r)
		return
	}

	if m == http.MethodGet {
		get(w,r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

//上传文件
func put(w http.ResponseWriter, r *http.Request)  {
	file, err := os.Create(os.Getenv("STORAGE_ROOT") + "/objectStorage/" + strings.Split(r.URL.EscapedPath(), "/")[2])

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	io.Copy(file, r.Body)
}

//下载文件
func get(w http.ResponseWriter, r *http.Request)  {
	file, err := os.Open(os.Getenv("STORAGE_ROOT") + "/objectStorage/" + strings.Split(r.URL.EscapedPath(), "/")[2])

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()
	io.Copy(w, file)
}

func GetObject(context *gin.Context)  {
	log.Debug("start listening get objects")
}

func UploadObject(context *gin.Context)  {
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
