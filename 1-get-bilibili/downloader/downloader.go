package downloader

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// InfoRequest 自定义结构体，存放要查询的bvid
type InfoRequest struct {
	Bvids []string
}

type VideInfo struct {
	Code int `json:"code"`
	Data struct {
		Bvid  string `json:"bvid"` // 使用tag
		Title string `json:"title"`
		Desc  string `json:"desc"`
	} `json:"data"`
}

type InfoResponse struct {
	Infos []VideInfo
}

// BatchDownloadVideoInfo 自定义的函数，标准写法是将error作为输出之一
func BatchDownloadVideoInfo(request InfoRequest) (InfoResponse, error) {
	var response InfoResponse
	for _, bvid := range request.Bvids {
		var videoInfo VideInfo
		resp, err := http.Get("https://api.bilibili.com/x/web-interface//view?bvid=" + bvid)
		if err != nil {
			return InfoResponse{}, err
		}
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return InfoResponse{}, err
		}

		// TODO：从网络中读取的是编码后的字节数组，将其解码为json格式
		if err = json.Unmarshal(respBytes, &videoInfo); err != nil {
			return InfoResponse{}, err
		}
		response.Infos = append(response.Infos, videoInfo)

		// 最后，所有的流，在使用完毕之后，都应该关闭
		if err = resp.Body.Close(); err != nil {
			return InfoResponse{}, err
		}
	}
	return response, nil
}
