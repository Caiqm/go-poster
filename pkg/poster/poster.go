package poster

import (
	"errors"
	"github.com/Caiqm/go-poster/pkg/file"
	"github.com/Caiqm/go-poster/pkg/qrcode"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/freetype"
)

type Poster struct {
	PosterName string
	Qr         *qrcode.QrCode
}

// 初始化海报参数
func NewPoster(posterName string, qr *qrcode.QrCode) *Poster {
	return &Poster{
		PosterName: posterName,
		Qr:         qr,
	}
}

func GetPosterFlag() string {
	return "poster"
}

// 检查合并后图像（指的是存放合并后的海报）是否存在
func (a *Poster) CheckMergedImage(path string) bool {
	if file.CheckNotExist(path+a.PosterName) == true {
		return false
	}

	return true
}

// 若不存在，则生成待合并的图像 mergedF
func (a *Poster) OpenMergedImage(path string) (*os.File, error) {
	f, err := file.MustOpen(a.PosterName, path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

type PosterBg struct {
	Name string
	*Poster
	*Rect
	*Pt
	*DrawText
	*DrawCover
}

type Rect struct {
	Name string
	X0   int
	Y0   int
	X1   int
	Y1   int
}

type Pt struct {
	X int
	Y int
}

type DrawText struct {
	JPG    draw.Image
	Merged *os.File

	TextMap []TextMain
}

type DrawCover struct {
	JPG    draw.Image
	Merged *os.File

	CoverMap []Cover
}

type TextMain struct {
	Text     string
	TextX    int
	TextY    int
	TextSize int
}

type Cover struct {
	CoverX      int
	CoverY      int
	CoverWidth  int
	CoverHeight int
	CoverUrl    string
}

func NewPosterBg(name string, ap *Poster, rect *Rect, pt *Pt, drawText *DrawText, drawCover *DrawCover) *PosterBg {
	return &PosterBg{
		Name:      name,
		Poster:    ap,
		Rect:      rect,
		Pt:        pt,
		DrawText:  drawText,
		DrawCover: drawCover,
	}
}

func (a *PosterBg) Generate() (string, string, error) {
	// 获取二维码存储路径
	fullPath := qrcode.GetQrCodeFullPath()
	// 生成二维码图像
	fileName, path, err := a.Qr.Encode(fullPath)
	if err != nil {
		return "", "", err
	}
	// 检查合并后图像（指的是存放合并后的海报）是否存在
	if !a.CheckMergedImage(path) {
		// 若不存在，则生成待合并的图像 mergedF
		mergedF, err := a.OpenMergedImage(path)
		if err != nil {
			return "", "", err
		}
		defer mergedF.Close()
		// 打开事先存放的背景图 bgF
		bgF, err := file.MustOpen(a.Name, path)
		if err != nil {
			return "", "", err
		}
		defer bgF.Close()
		// 打开生成的二维码图像 qrF
		qrF, err := file.MustOpen(fileName, path)
		if err != nil {
			return "", "", err
		}
		defer qrF.Close()
		// 解码 bgF 和 qrF 返回 image.Image
		bgImage, err := jpeg.Decode(bgF)
		if err != nil {
			return "", "", err
		}
		qrImage, err := jpeg.Decode(qrF)
		if err != nil {
			return "", "", err
		}
		// 创建一个新的 RGBA 图像
		jpg := image.NewRGBA(image.Rect(a.Rect.X0, a.Rect.Y0, a.Rect.X1, a.Rect.Y1))
		// 在 RGBA 图像上绘制 背景图（bgF）
		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
		// 在已绘制背景图的 RGBA 图像上，在指定 Point 上绘制二维码图像（qrF）
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X, a.Pt.Y)), draw.Over)
		// 写入图片
		a.DrawCover.JPG = jpg
		a.DrawCover.Merged = mergedF
		err = a.DrawPosterCover(a.DrawCover)
		if err != nil {
			return "", "", err
		}
		// 写入文字
		a.DrawText.JPG = jpg
		a.DrawText.Merged = mergedF
		err = a.DrawPosterText(a.DrawText, "simhei.ttf")
		if err != nil {
			return "", "", err
		}
		// 将绘制好的 RGBA 图像以 JPEG 4：2：0 基线格式写入合并后的图像文件（mergedF）
		jpeg.Encode(mergedF, jpg, nil)
	}

	return fileName, path, nil
}

// 画入图片
func (a *PosterBg) DrawPosterCover(dp *DrawCover) error {
	if len(dp.CoverMap) <= 0 {
		return nil
	}
	var err error
	for _, v := range dp.CoverMap {
		filePath := v.CoverUrl
		// 网络链接
		if strings.HasPrefix(filePath, "https://") || strings.HasPrefix(filePath, "http://") {
			filePath, err = file.DownloadFile(filePath)
			if err != nil {
				return err
			}
		} else {
			dir, _ := os.Getwd()
			filePath = dir + "/" + filePath
		}
		// 普通路径
		if file.CheckNotExist(filePath) {
			return errors.New("file path not exist")
		}
		// 打开图
		coverF, err := file.Open(filePath, os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			return err
		}
		defer coverF.Close()
		var coverImage image.Image
		// 解码返回 image.Image
		if strings.HasSuffix(filePath, ".png") {
			coverImage, err = png.Decode(coverF)
		} else {
			coverImage, err = jpeg.Decode(coverF)
		}
		if err != nil {
			return err
		}
		// 在 RGBA 图像上绘制 背景图（bgF）
		draw.Draw(dp.JPG, dp.JPG.Bounds(), coverImage, coverImage.Bounds().Min.Sub(image.Pt(v.CoverX, v.CoverY)), draw.Over)
		//err = jpeg.Encode(dp.Merged, dp.JPG, nil)
		//if err != nil {
		//	return err
		//}
	}
	return nil
}

// 写入文字
func (a *PosterBg) DrawPosterText(d *DrawText, fontName string) error {
	if len(d.TextMap) <= 0 {
		return nil
	}
	// 字体文件路径
	fontSource := "runtime/fonts/" + fontName
	fontSourceBytes, err := ioutil.ReadFile(fontSource)
	if err != nil {
		return err
	}

	trueTypeFont, err := freetype.ParseFont(fontSourceBytes)
	if err != nil {
		return err
	}
	// 创建一个新的 Context，会对其设置一些默认值
	fc := freetype.NewContext()
	// 设置屏幕每英寸的分辨率
	fc.SetDPI(72)
	// 设置用于绘制文本的字体
	fc.SetFont(trueTypeFont)
	// 设置剪裁矩形以进行绘制
	fc.SetClip(d.JPG.Bounds())
	// 设置目标图像
	fc.SetDst(d.JPG)
	// 设置绘制操作的源图像
	fc.SetSrc(image.Black)
	// 以磅为单位设置字体大小
	for _, v := range d.TextMap {
		fc.SetFontSize(float64(v.TextSize))
		// 根据 Pt 的坐标值绘制给定的文本内容
		_, err = fc.DrawString(v.Text, freetype.Pt(v.TextX, v.TextY))
		if err != nil {
			return err
		}
	}

	err = jpeg.Encode(d.Merged, d.JPG, nil)
	if err != nil {
		return err
	}

	return nil
}
