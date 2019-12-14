package main

import (
	"fmt"
	"net/url"
)

func main()  {
	name := "%yamanda&&"
	e := url.PathEscape(name)
	fmt.Println(e)
}
