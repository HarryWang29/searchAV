package search

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	cmap "github.com/orcaman/concurrent-map"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type search struct {
	proxy      string
	keyWord    string
	resultFile string
	url        string

	magnetMap cmap.ConcurrentMap
}

type Search interface {
	Run()
	getDoc(link string) (*goquery.Document, error)
	getTitleAndMagnet(doc *goquery.Document)
	getMagnet(link string) string
	setHeader(header *http.Header)
}

func NewSearch(proxy, keyWord, resultFile, url string) Search {
	s := &search{
		proxy:      proxy,
		keyWord:    keyWord,
		resultFile: resultFile,
		url:        url,
	}
	s.magnetMap = cmap.New()
	return s
}

func (s *search) Run() {
	doc, err := s.getDoc(s.url)
	if err != nil {
		fmt.Printf("GetDoc err:%s", err)
		return
	}
	s.getTitleAndMagnet(doc)
}

func (s *search) setHeader(header *http.Header) {
	//增加header选项
	header.Add("Host", "btso.pw")
	header.Add("Connection", "keep-alive")
	header.Add("Pragma", "no-cache")
	header.Add("Cache-Control", "no-cache")
	header.Add("Upgrade-Insecure-Requests", "1")
	header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36")
	header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	header.Add("Referer", "https://btso.pw/search/")
	//header.Add("Accept-Encoding", "gzip, deflate, br")
	header.Add("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8")
	header.Add("Cookie", "__test; _ga=GA1.2.613298755.1516200848; 494668b4c0ef4d25bda4e75c27de2817=326dec85-8387-4503-80a9-9d0c7cc7ee40%3A2%3A1; a=9nbjphmhasplh6ygdxayf4p0xt2u2ycc; _gid=GA1.2.950150230.1516713710; AD_enterTime=1516713710; AD_adca_b_SM_T_728x90=0; AD_jav_b_SM_T_728x90=0; AD_javu_b_SM_T_728x90=0; AD_wav_b_SM_T_728x90=0; AD_wwwp_b_SM_T_728x90=0; AD_adst_b_SM_T_728x90=1; __PPU_SESSION_1_470916_false=1516713723903|1|1516713723903|1|1; AD_exoc_b_SM_T_728x90=1; AD_clic_b_POPUNDER=2")
}

func (s *search) getTitleAndMagnet(doc *goquery.Document) {
	titles := make([]string, 0)
	links := make([]string, 0)
	doc.Find("div[class='row'] a").Each(func(i int, selection *goquery.Selection) {
		title, ok := selection.Attr("title")
		if !ok {
			return
		}
		titles = append(titles, title)

		link, ok := selection.Attr("href")
		if !ok {
			return
		}
		links = append(links, link)
	})
	if len(titles) != len(links) {
		fmt.Printf("get title:%d link:%d\n", len(titles), len(links))
		return
	}
	var wg sync.WaitGroup
	for i, l := range links {
		wg.Add(1)
		go func(title, link string) {
			defer wg.Done()
			magnet := s.getMagnet(link)
			s.magnetMap.Set(title, magnet)
		}(titles[i], l)
	}
	wg.Wait()
	var f *os.File
	var exist = true
	var err error
	if _, err := os.Stat(s.resultFile); os.IsNotExist(err) {
		exist = false
	}
	if exist { //如果文件存在
		f, err = os.OpenFile(s.resultFile, os.O_APPEND, 0666) //打开文件
	} else {
		f, err = os.Create(s.resultFile) //创建文件
	}
	w := bufio.NewWriter(f)

	fmt.Fprintln(w, fmt.Sprintf("search at: %s", time.Now().Format("2006-01-02 15:04:05")))

	for _, t := range titles {
		magnet, _ := s.magnetMap.Get(t)
		fmt.Printf("%s:%s\n", t, magnet.(string))
		if err == nil {
			fmt.Fprintln(w, fmt.Sprintf("%s:%s", t, magnet.(string)))
		}
	}
	w.Flush()
	f.Close()
}

func (s *search) getMagnet(link string) string {
	doc, err := s.getDoc(link)
	if err != nil {
		fmt.Printf("GetDoc err:%s", err)
		return ""
	}
	return doc.Find("#magnetLink").Text()
}

func (s *search) getDoc(link string) (*goquery.Document, error) {
	urli := url.URL{}
	urlproxy, _ := urli.Parse(s.proxy)

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
		},
	}
	//提交请求
	reqest, err := http.NewRequest("GET", link, nil)
	if err != nil {
		fmt.Printf("err:%s", err)
		return nil, err
	}
	s.setHeader(&reqest.Header)
	response, err := client.Do(reqest)
	if err != nil {
		fmt.Printf("client.Do err:%s", err)
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Printf("code:%s\n", response.Status)
		return nil, fmt.Errorf("response.Status:%s", response.Status)
	}
	return goquery.NewDocumentFromReader(response.Body)
}
