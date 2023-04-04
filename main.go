package main

import (
	"fmt"
	"github.com/Caiqm/go-poster/pkg/file"
	"github.com/Caiqm/go-poster/poster"
)

func main() {
	posterClient := poster.NewPosterClient("bg.jpg", "https://baidu.com", 220, 220, 30, 1000)
	posterClient.SetText("Caiqm", 360, 1135, 26)
	//posterClient.SetCover("20230402105515_6428ee93b18dd.jpg", 630, 630, 60, 250, 0)
	//posterClient.SetCover("user_avatar_o9LGt6oJBtSfxv7FlE2jl1Uvzhik.jpeg", 70, 70, 270, 1130, 1)
	filePath, err := posterClient.CreatePoster()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(filePath)
	f, _ := file.CreateThumb(filePath, 750, 1334)
	fmt.Println(f)
}
