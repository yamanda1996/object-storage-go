package objectstream

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type PutStream struct {
	writer 		*io.PipeWriter
	c 			chan error
}

func NewPutStream(server, object string) *PutStream {
	//关于io.PipeReader和io.PipeWriter,写入writer中的内容可以从reader中读出来
	reader, writer := io.Pipe()
	c := make(chan error)

	go func() {
		request, err := http.NewRequest("POST", "http://"+server+"/objects/"+object, reader)
		if err != nil {
			log.Errorf("create request to data server [%s] failed", server)
		}
		client := http.Client{}

		response, err := client.Do(request)
		if err != nil || response.StatusCode != http.StatusOK {
			err = fmt.Errorf("dataServer return http code %d", response.StatusCode)
			log.Errorf("data server return failed")
		}
		c <- err
	}()

	return &PutStream{writer, c}
}

func (ps *PutStream) Write(p []byte) (n int, err error) {
	return ps.writer.Write(p)
}

func (ps *PutStream) Close() error {
	ps.writer.Close()
	return <- ps.c
}

type GetStream struct {
	reader io.Reader
}

func NewGetStream(server, object string) (*GetStream, error) {
	if server == "" || object == "" {
		return nil, fmt.Errorf("invalid server %s object %s", server, object)
	}

	return newGetStream("http://" + server + "/objects/" + object)
}

func newGetStream(url string) (*GetStream, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("data server return code %d", resp.StatusCode)
	}
	return &GetStream{resp.Body}, nil
}

func (gs *GetStream) Read(p []byte) (n int, err error) {
	return gs.reader.Read(p)
}



