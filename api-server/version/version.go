package version

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"object-storage-go/api-server/es"
)

func GetVersion(context *gin.Context)  {
	filename := context.Param("filename")
	from := 0
	size := 500
	for {
		metadatas, err := es.SearchAllVersions(filename, from, size)
		if err != nil {
			log.Errorf("search all [%s] metadata failed", filename)
			context.String(http.StatusInternalServerError, "search all metadata failed")
		}
		r, w := io.Pipe()
		for _, metadata := range metadatas {
			result, _ := json.Marshal(metadata)
			w.Write(result)
			w.Write([]byte("\n"))
		}
		if len(metadatas) != size {
			bytes, _ := ioutil.ReadAll(r)
			context.Data(http.StatusOK, "json", bytes)
		}
		from += size
	}
}
