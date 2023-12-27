package poster

import (
	"errors"
	"fmt"
	"github.com/Caiqm/go-poster/pkg/circle"
	"github.com/Caiqm/go-poster/pkg/file"
	"github.com/Caiqm/go-poster/pkg/qrcode"
	"github.com/golang/freetype"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Poster struct {
	PosterName string
	Qr         *qrcode.QrCode
}

type PosterBg struct {
	Name             string
	CustomQrCodePath string
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
	Circle      int
	CoverUrl    string
}

// 初始化海报参数
func NewPoster(posterName string) *Poster {
	return &Poster{
		PosterName: posterName,
	}
}

// 初始化生成海报参数
func NewPosterBg(name, cusQrPath string, ap *Poster, pt *Pt, drawText *DrawText, drawCover *DrawCover) *PosterBg {
	return &PosterBg{
		Name:             name,
		CustomQrCodePath: cusQrPath,
		Poster:           ap,
		Pt:               pt,
		DrawText:         drawText,
		DrawCover:        drawCover,
	}
}

// 海报前缀
func GetPosterFlag() string {
	return "poster"
}

// 设置二维码参数
func (a *Poster) SetQrParam(qr *qrcode.QrCode) *Poster {
	a.Qr = qr
	return a
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

// 获取文件本地地址
func (a *Poster) GetRealFilePath(filePath string) (string, error) {
	fileRealPath := filePath
	if file.CheckIsHttpResource(filePath) {
		fileRealPathTmp, err := file.DownloadFile(filePath)
		if err != nil {
			return "", err
		}
		fileRealPath = fileRealPathTmp
	}
	return fileRealPath, nil
}

// 生成海报主方法
func (a *PosterBg) Generate() (string, string, error) {
	// 获取背景图本地地址
	filePath, err := a.GetRealFilePath(a.Name)
	if err != nil {
		return "", "", err
	}
	var path, fileName string
	// 生成二维码图像
	if a.CustomQrCodePath != "" {
		path = filepath.Dir(a.CustomQrCodePath)
		fileName = filepath.Base(a.CustomQrCodePath)
	} else {
		// 获取二维码存储路径
		fullPath := qrcode.GetQrCodeFullPath()
		fileName, path, err = a.Qr.Encode(fullPath)
		if err != nil {
			return "", "", err
		}
	}
	// 检查合并后图像（指的是存放合并后的海报）是否存在
	if !a.CheckMergedImage(path) {
		// 若不存在，则生成待合并的图像 mergedF
		mergedF, err2 := a.OpenMergedImage(path)
		if err2 != nil {
			return "", "", err2
		}
		defer mergedF.Close()
		// 打开事先存放的背景图 bgF
		bgF, err2 := file.Open(filePath, os.O_RDWR, 0666)
		if err2 != nil {
			return "", "", err2
		}
		defer bgF.Close()
		// 打开生成的二维码图像 qrF
		qrF, err2 := file.MustOpen(fileName, path)
		if err2 != nil {
			return "", "", err2
		}
		defer qrF.Close()
		// 解码 bgF 和 qrF 返回 image.Image
		bgImage, _, err2 := image.Decode(bgF)
		if err2 != nil {
			return "", "", err2
		}
		qrImage, err2 := jpeg.Decode(qrF)
		if err2 != nil {
			return "", "", err2
		}
		bgSize := bgImage.Bounds().Size()
		// 创建一个新的 RGBA 图像
		jpg := image.NewRGBA(image.Rect(0, 0, bgSize.X, bgSize.Y))
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
		err = png.Encode(mergedF, jpg)
		if err != nil {
			return "", "", err
		}
	}
	return fileName, path, nil
}

// 画入图片
func (a *PosterBg) DrawPosterCover(dp *DrawCover) error {
	if len(dp.CoverMap) <= 0 {
		return nil
	}
	var (
		err       error
		coverPath []string
	)
	// 删除文件
	defer func() {
		if len(coverPath) > 0 {
			for _, fp := range coverPath {
				err = file.RemoveFile(fp)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()
	for _, v := range dp.CoverMap {
		filePath := v.CoverUrl
		// 网络链接
		if file.CheckIsHttpResource(filePath) {
			filePath, err = file.DownloadFile(filePath)
			if err != nil {
				return err
			}
		} else {
			dir, _ := os.Getwd()
			filePath = filepath.Join(dir, filePath)
		}
		var (
			orgFilePath  string
			orgThumbPath string
		)
		// 普通路径
		if file.CheckNotExist(filePath) {
			return errors.New("file path not exist")
		}
		// 文件路径加入切片
		coverPath = append(coverPath, filePath)
		// 生成缩略图
		if v.CoverWidth > 0 {
			orgFilePath = filePath
			filePath, err = file.CreateThumb(filePath, v.CoverWidth, v.CoverHeight)
			if err != nil {
				return err
			}
		}
		// 生成圆形缩略图
		if v.Circle > 0 {
			orgThumbPath = filePath
			filePath, err = circle.ThumbToCircle(filePath, v.CoverWidth, v.CoverHeight)
			if err != nil {
				return err
			}
		}
		// 打开图
		coverF, err1 := file.Open(filePath, os.O_APPEND|os.O_RDWR, 0666)
		if err1 != nil {
			return err1
		}
		defer coverF.Close()
		var coverImage image.Image
		// 解码返回 image.Image
		coverImage, _, err = image.Decode(coverF)
		if err != nil {
			return err
		}
		// 在 RGBA 图像上绘制 背景图（bgF）
		draw.Draw(dp.JPG, dp.JPG.Bounds(), coverImage, coverImage.Bounds().Min.Sub(image.Pt(v.CoverX, v.CoverY)), draw.Over)
		// 删除原文件
		_ = RemoveOrgFile(orgFilePath)
		_ = RemoveOrgFile(orgThumbPath)
	}
	return nil
}

// 写入文字
func (a *PosterBg) DrawPosterText(d *DrawText, fontName string) error {
	if len(d.TextMap) <= 0 {
		return nil
	}
	// 字体文件路径
	fontSource := "fonts/" + fontName
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

// 删除文件
func RemoveOrgFile(orgFilePath string) error {
	// 删除原文件
	if orgFilePath == "" {
		return errors.New("file not exist")
	}
	err := file.RemoveFile(orgFilePath)
	if err != nil {
		return err
	}
	return nil
}
