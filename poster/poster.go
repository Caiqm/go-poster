package poster

import (
	posterPkg "github.com/Caiqm/go-poster/pkg/poster"
	"github.com/Caiqm/go-poster/pkg/qrcode"
	"github.com/boombuler/barcode/qr"
	"os"
)

func (p *PosterParams) CreatePoster() (string, error) {
	// 二维码参数
	qrc := qrcode.NewQrCode(p.QrCodeUrl, p.QrCodeWidth, p.QrCodeHeight, qr.M, qr.Auto)
	// 海报名称
	posterName := posterPkg.GetPosterFlag() + "-" + qrcode.GetQrCodeFileName(qrc.URL) + qrc.GetQrCodeExt()
	// 海报参数
	poster := posterPkg.NewPoster(posterName, qrc)
	// 参数：背景图路径，海报名称和二维码信息，新绘制背景图大小，二维码坐标，文字坐标
	posterBg := posterPkg.NewPosterBg(
		p.BgUrl,
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
	posterSaveUrl := filePath + posterName
	// 删除文件
	_ = os.Remove(filePath + qrCodeName)
	return posterSaveUrl, nil
}
