package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Chengxufeng1994/go-crawler/collect"
	"github.com/Chengxufeng1994/go-crawler/proxy"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	proxyUrls := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8889"}
	p, _ := proxy.NewRoundRobinProxySwitcher(proxyUrls...)

	url := "https://google.com"
	f := collect.NewBaseFetch(500*time.Millisecond, p)
	body, err := f.Get(url)
	if err != nil {
		fmt.Printf("read content failed:%v", err)
		return
	}

	fmt.Println(string(body))

	// 加载HTML文档
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("read content failed:%v", err)
		return
	}

	doc.Find("div.news_li h2 a[target=_blank]").Each(func(i int, s *goquery.Selection) {
		// 获取匹配元素的文本
		title := s.Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
}
