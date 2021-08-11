package main        // 入口

import (
	`1-get-bilibili/downloader`
	myfmt "1-get-bilibili/fmt"
	"fmt"
)

func main() {
	fmt.Println("hello")
	myfmt.Logger.Println("hello")

	request := downloader.InfoRequest{Bvids: []string{"BV1Ff4y187q9", "BV1DV411s7ij"}}
	response, err := downloader.BatchDownloadVideoInfo(request); if err != nil {
		panic(err)
	}

	for _, info := range response.Infos {
		myfmt.Logger.Printf("title: %s\ndesc: %s\n\n", info.Data.Title, info.Data.Desc)
	}
}
