package main

import (
	"fmt"
	"github.com/Caiqm/go-poster/poster"
	"time"
)

func main() {
	start := time.Now()
	posterClient := poster.NewPosterClient("bg.jpg")
	posterClient.SetCustomQrCodePath("7_qrcode.jpg", 350, 350, 335, 590)
	//posterClient.SetQrcodeParam("https://baidu.com", 220, 220, 30, 1000)
	//posterClient.SetText("Caiqm", 360, 1135, 26)
	//posterClient.SetCover("20230526173155_64707c8b2a5f3.jpg", 630, 630, 60, 250, 0)
	//posterClient.SetCover("user_avatar_o9LGt6h2-5TIpNLLzNLARu_GRe1c.png", 70, 70, 270, 1130, 1)
	filePath, err := posterClient.CreatePoster()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(filePath)
	//f, _ := file.CreateThumb(filePath, 750, 1334)
	//fmt.Println(f)
	fmt.Println(time.Since(start).Seconds())
}
