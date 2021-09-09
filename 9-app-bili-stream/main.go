package main

import (
	bili "github.com/JimmyZhangJW/biliStreamClient"
	`log`
)

func main() {
	biliClient := bili.New()
	err := biliClient.Connect(977262)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer biliClient.Disconnect()
	for {
		// 从通过中取数据并根据其类型出发不同的操作
		packBody := <- biliClient.Ch
		switch packBody.Cmd {
		case "DANMU_MSG":
			msg, _ := packBody.ParseDanmu()
			log.Println(msg)
			// 拿到msg之后，可以结合腾讯云提供的接口将弹幕用语音的方式朗读出来
		case "SEND_GIFT":
			log.Println(packBody.ParseGift())
		case "COMBO_SENT":
			log.Println(packBody.ParseGiftCombo())
		default:
			log.Println(packBody.Cmd)
		}
	}
}
