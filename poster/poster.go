package poster

import (
	"fmt"
	"github.com/Caiqm/go-poster/pkg/file"
	posterPkg "github.com/Caiqm/go-poster/pkg/poster"
	"github.com/Caiqm/go-poster/pkg/qrcode"
	"github.com/boombuler/barcode/qr"
	"os"
	"path/filepath"
	"strings"
)

func (p *PosterParams) CreatePoster() (string, error) {
	var (
		qrc         *qrcode.QrCode
		tmpFileName string
		mold        int64
	)
	dir, _ := os.Getwd()
	// 没有自定义二维码
	if p.CustomQrCodePath == "" {
		// 二维码参数
		qrc = qrcode.NewQrCode(p.QrCodeUrl, p.QrCodeWidth, p.QrCodeHeight, qr.M, qr.Auto)
		tmpFileName = qrcode.GetQrCodeFileName(qrc.URL) + qrc.GetQrCodeExt()
		mold = 1
	} else {
		var err error
		// 压缩二维码图片
		p.CustomQrCodePath, err = file.CreateThumb(p.CustomQrCodePath, p.QrCodeWidth, p.QrCodeHeight)
		if err != nil {
			return "", fmt.Errorf("压缩图片失败：%w", err)
		}
		p.CustomQrCodePath = strings.Replace(p.CustomQrCodePath, dir, "", 1)
		// 自定义二维码路径
		tmpFileName = qrcode.GetQrCodeFileName(p.CustomQrCodePath) + filepath.Ext(p.CustomQrCodePath)
	}
	// 海报名称
	posterName := fmt.Sprintf("%s-%s", posterPkg.GetPosterFlag(), tmpFileName)
	// 海报参数
	poster := posterPkg.NewPoster(posterName)
	// 设置二维码参数
	if mold == 1 {
		poster.SetQrParam(qrc)
	}
	// 参数：背景图路径，海报名称和二维码信息，新绘制背景图大小，二维码坐标，文字坐标
	posterBg := posterPkg.NewPosterBg(
		p.BgUrl,
		p.CustomQrCodePath,
		poster,
		&posterPkg.Pt{
			X: p.QrCodeX,
			Y: p.QrCodeY,
		},
		&posterPkg.DrawText{
			TextMap: p.Text,
		},
		&posterPkg.DrawCover{
			CoverMap: p.Cover,
		},
	)
	qrCodeName, filePath, err := posterBg.Generate()
	if err != nil {
		return "", err
	}
	// 保存海报路径
	posterSaveUrl := filepath.Join(filePath, posterName)
	// 删除二维码图片
	_ = os.Remove(filepath.Join(dir, filePath, qrCodeName))
	return posterSaveUrl, nil
}
