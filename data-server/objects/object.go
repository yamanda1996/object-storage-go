package objects

import (
	log "github.com/sirupsen/logrus"
	"go-common/app/service/ops/log-agent/pkg/bufio"
	"io"
	"object-storage-go/common/common_grpc/objectpb"
	"object-storage-go/data-server/model"
	"os"
)

type Object struct {

}

func (o *Object) UploadObject(stream objectpb.Object_UploadObjectServer) error {
	log.Info("data server start to store object")

	var length int64 = 0
	count := 0
	for {
		chunk, err := stream.Recv()
		log.Info("receive chunk from api server")
		// store chunk
		err = storeChunk(chunk, &length)
		if err != nil {
			log.Error("store %d chunk failed")
			return stream.SendAndClose(&objectpb.UploadReply{
				Success:false,
				Length:-1,
			})
		}
		count++
		if err == io.EOF {
			log.Infof("store object success, chunk %d", count)
			return stream.SendAndClose(&objectpb.UploadReply{
				Success:true,
				Length:length,
			})
		}
	}
}

func storeChunk(chunk *objectpb.ObjectChunk, length *int64) error {
	filename := model.Config.DataServerConfig.StorageRoot + "/objects/" +
		chunk.Filename

	var file *os.File
	if IsExist(filename) {
		file, _ = os.OpenFile(filename, os.O_APPEND, 0644)
	} else {
		file, _ = os.Create(filename)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	len, _ := w.Write(chunk.Buffer)
	w.Flush()

	*length = *length + int64(len)
	return nil
}

func IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}