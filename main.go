package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Chengxufeng1994/go-crawler/collect"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://book.douban.com/subject/1007305/"
	f := collect.NewBaseFetch(5 * time.Second)
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
	}

	doc.Find("div.review-list h2 a").Each(func(i int, s *goquery.Selection) {
		// 获取匹配元素的文本
		title := s.Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
}
