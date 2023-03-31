package main

import (
	"github.com/Caiqm/go-poster/poster"
)

func main() {
	posterClient := poster.NewPosterClient("bg.jpg", "https://baidu.com", 300, 300, 250, 300)
	posterClient.SetText("这是文字内容", 30, 50, 28)
	posterClient.SetText("这是第二行文字内容", 30, 150, 18)
	posterClient.SetCover("1.png", 300, 300, 100, 300)
	posterClient.CreatePoster()
}
