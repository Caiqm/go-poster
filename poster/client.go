package poster

import "github.com/Caiqm/go-poster/pkg/poster"

type PosterParams struct {
	QrCodeWidth      int
	QrCodeHeight     int
	QrCodeX          int
	QrCodeY          int
	QrCodeUrl        string
	BgUrl            string
	CustomQrCodePath string
	Text             []poster.TextMain
	Cover            []poster.Cover
}

func NewPosterClient(bgUrl string) *PosterParams {
	return &PosterParams{
		BgUrl: bgUrl,
	}
}

// 设置自定义二维码路径
func (p *PosterParams) SetCustomQrCodePath(qrcodePath string, qrWidth, qrHeight, qrX, qrY int) *PosterParams {
	p.CustomQrCodePath = qrcodePath
	p.QrCodeWidth = qrWidth
	p.QrCodeHeight = qrHeight
	p.QrCodeX = qrX
	p.QrCodeY = qrY
	return p
}

// 设置二维码参数
func (p *PosterParams) SetQrcodeParam(qrUrl string, qrWidth, qrHeight, qrX, qrY int) *PosterParams {
	p.QrCodeUrl = qrUrl
	p.QrCodeWidth = qrWidth
	p.QrCodeHeight = qrHeight
	p.QrCodeX = qrX
	p.QrCodeY = qrY
	return p
}

// 设置文字
func (p *PosterParams) SetText(txt string, x, y, size int) *PosterParams {
	var t poster.TextMain
	t.Text = txt
	t.TextX = x
	t.TextY = y
	t.TextSize = size
	p.Text = append(p.Text, t)
	return p
}

// 设置封面
func (p *PosterParams) SetCover(url string, width, height, x, y, circle int) *PosterParams {
	var c poster.Cover
	c.CoverX = x
	c.CoverY = y
	c.CoverWidth = width
	c.CoverHeight = height
	c.CoverUrl = url
	c.Circle = circle
	p.Cover = append(p.Cover, c)
	return p
}
