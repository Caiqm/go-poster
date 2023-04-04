# go-poster

用golang来生成海报

#### 使用说明

下载代码库

```
go get github.com/Caiqm/go-poster
```

普通二维码海报

```golang
package main

import (
	"fmt"
	"github.com/Caiqm/go-poster/poster"
)

func main() {
	// 背景图路径(根目录起)、二维码跳转链接、二维码宽度、二维码高度、相对背景图x坐标、相对背景图y坐标
	posterClient := poster.NewPosterClient("bg.jpg", "https://baidu.com", 220, 220, 30, 1000)
	// 生成海报
	filePath, err := posterClient.CreatePoster()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 输出保留路径
	fmt.Println(filePath)
}
```

带文字二维码海报

```golang
package main

import (
	"fmt"
	"github.com/Caiqm/go-poster/poster"
)

func main() {
	// 背景图路径(根目录起)、二维码跳转链接、二维码宽度、二维码高度、相对背景图x坐标、相对背景图y坐标
	posterClient := poster.NewPosterClient("bg.jpg", "https://baidu.com", 220, 220, 30, 1000)
	// 文字内容、相对背景图x坐标、相对背景图y坐标、文字大小
	posterClient.SetText("xxx", 360, 1135, 26)
	posterClient.SetText("yyy", 360, 1160, 26)
	// 生成海报
	filePath, err := posterClient.CreatePoster()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 输出保留路径
	fmt.Println(filePath)
}
```

带水印文字二维码海报

```golang
package main

import (
	"fmt"
	"github.com/Caiqm/go-poster/poster"
)

func main() {
	// 背景图路径(根目录起)、二维码跳转链接、二维码宽度、二维码高度、相对背景图x坐标、相对背景图y坐标
	posterClient := poster.NewPosterClient("bg.jpg", "https://baidu.com", 220, 220, 30, 1000)
	// 文字内容、相对背景图x坐标、相对背景图y坐标、文字大小
	posterClient.SetText("xxx", 360, 1135, 26)
	posterClient.SetText("yyy", 360, 1160, 26)
	// 水印路径(根目录起)、宽度、高度、相对背景图x坐标、相对背景图y坐标、是否截取圆形
	posterClient.SetCover("avatar.jpeg", 70, 70, 270, 1130, 1)
	// 生成海报
	filePath, err := posterClient.CreatePoster()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 输出保留路径
	fmt.Println(filePath)
}
```

压缩图片

```golang
package main

import (
	"fmt"
	"github.com/Caiqm/go-poster/pkg/file"
)

func main() {
	// 图片路径(根目录起)、宽度、高度
	f, _ := file.CreateThumb("p.jpg", 750, 1334)
	fmt.Println(f)
}
```