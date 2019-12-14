package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func Locate(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func GetUrl(url string) {
	resp, _ := http.Get(url)

	bytes, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(bytes))
}


func main()  {
	//res := Locate("/Users/yangmoda/code/golang/beego/src/object-storage-go/data-server/data/objects/hello.cpp")
	//
	//fmt.Println(res)

	GetUrl("http://localhost:8792/objects/hello.cpp")

}
