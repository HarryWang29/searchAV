package main

import (
	"flag"
	"fmt"
	"searchAV/search"
)

var proxyAddr string
var keyWord string
var resultFile string

const Url = "https://btsow.pw/search/"

func init() {
	flag.StringVar(&proxyAddr, "proxy", "127.0.0.1:1087", "")
	flag.StringVar(&keyWord, "search", "", "")
	flag.StringVar(&resultFile, "save", "./result.txt", "")
}

func main() {
	flag.Parse()
	if keyWord == "" {
		fmt.Printf("write search\n")
		return
	}
	fmt.Printf("search:%s\nproxy:%s\n", keyWord, proxyAddr)
	search.NewSearch(
		proxyAddr,
		keyWord,
		resultFile,
		Url+keyWord,
	).Run()
}
