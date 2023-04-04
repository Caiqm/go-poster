package circle

import (
	"github.com/Caiqm/go-poster/pkg/file"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

// 图片画圆
func ThumbToCircle(filePath string, width, height int) (string, error) {
	// Load the image
	imgFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return "", err
	}
	jpg := image.NewNRGBA(image.Rect(0, 0, width, height))

	radius := width / 2
	draw.DrawMask(jpg, jpg.Bounds(), img, image.Point{}, &circle{p: image.Point{X: radius, Y: radius}, r: radius}, image.Point{}, draw.Over)

	dir, _ := os.Getwd()
	thumbPath := dir + "/runtime/thumb"
	err = file.IsNotExistMkDir(thumbPath)
	if err != nil {
		return "", err
	}
	fileName := file.Md5Str(filePath+"output") + ".png"
	newFilePath := thumbPath + "/" + fileName
	// Save the new image to a file
	out, err := os.Create(newFilePath)
	if err != nil {
		return "", err
	}
	defer out.Close()
	// 生成图片
	err = png.Encode(out, jpg)
	if err != nil {
		return "", err
	}
	return newFilePath, nil
}
