package main

import (
	"github.com/Caiqm/go-poster/poster"
)

func main() {
	posterClient := poster.NewPosterClient("bg.jpg", "https://baidu.com", 220, 220, 30, 1000)
	posterClient.SetText("这是第二行文字内容", 360, 1135, 26)
	posterClient.SetCover("https://static.golangjob.cn/haoimg/wechat.jpg", 300, 300, 60, 250)
	posterClient.SetCover("https://static.golangjob.cn/haoimg/wechat.jpg", 300, 300, 270, 1130)
	posterClient.CreatePoster()
}
