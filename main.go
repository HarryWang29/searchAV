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
	flag.StringVar(&proxyAddr, "proxy", "", "")
	flag.StringVar(&keyWord, "search", "", "")
	flag.StringVar(&resultFile, "save", "./result.txt", "")
}

func main() {
	flag.Parse()
	if keyWord == "" {
		fmt.Printf("write search\n")
		return
	}
	fmt.Printf("search:%s\n", keyWord)
	if proxyAddr == "" {
		fmt.Printf("proxy:为空，不使用代理，请开启ss/ssr全局模式\n")
	} else {
		fmt.Printf("proxy:%s\n", keyWord)
	}
	search.NewSearch(
		proxyAddr,
		keyWord,
		resultFile,
		Url+keyWord,
	).Run()
}
