package qrcode

import (
	"github.com/Caiqm/go-poster/pkg/file"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/jpeg"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type QrCode struct {
	URL    string
	Width  int
	Height int
	Ext    string
	Level  qr.ErrorCorrectionLevel
	Mode   qr.Encoding
}

const (
	EXT_JPG = ".jpg"
)

func NewQrCode(url string, width, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		URL:    url,
		Width:  width,
		Height: height,
		Level:  level,
		Mode:   mode,
		Ext:    EXT_JPG,
	}
}

// 获取二维码路径
func GetQrCodePath() string {
	return "qrcode/"
}

// 获取完整二维码路径
func GetQrCodeFullPath() string {
	return "runtime/" + GetQrCodePath()
}

// 获取二维码访问路径
//func GetQrCodeFullUrl(name string) string {
//	return setting.AppSetting.PrefixUrl + "/" + GetQrCodePath() + name
//}

// 获取二维码文件名称
func GetQrCodeFileName(value string) string {
	timeStr := strconv.Itoa(int(time.Now().UnixMicro()))
	return file.Md5Str(value + timeStr)
}

// 获取二维码扩展名
func (q *QrCode) GetQrCodeExt() string {
	return q.Ext
}

// 检测是否存在
func (q *QrCode) CheckEncode(path string) bool {
	src := path + GetQrCodeFileName(q.URL) + q.GetQrCodeExt()
	if file.CheckNotExist(src) == true {
		return false
	}

	return true
}

// 生成二维码
func (q *QrCode) Encode(path string) (string, string, error) {
	name := GetQrCodeFileName(q.URL) + q.GetQrCodeExt()
	src := path + name
	if file.CheckNotExist(src) == true {
		code, err := qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			return "", "", err
		}

		code, err = barcode.Scale(code, q.Width, q.Height)
		if err != nil {
			return "", "", err
		}

		f, err := file.MustOpen(name, path)
		if err != nil {
			return "", "", err
		}
		defer f.Close()

		err = jpeg.Encode(f, code, nil)
		if err != nil {
			return "", "", err
		}
	}

	return name, path, nil
}

// 删除二维码
func (q *QrCode) RmQrcode(fileName, filePath string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	qrSrc := filepath.Join(dir, filePath, fileName)
	err = os.Remove(qrSrc)
	return err
}
