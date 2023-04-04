package poster

import "github.com/Caiqm/go-poster/pkg/poster"

type PosterParams struct {
	QrCodeUrl    string
	QrCodeWidth  int
	QrCodeHeight int
	QrCodeX      int
	QrCodeY      int
	BgUrl        string
	Text         []poster.TextMain
	Cover        []poster.Cover
}

func NewPosterClient(bgUrl, qrUrl string, qrWidth, qrHeight, qrX, qrY int) *PosterParams {
	return &PosterParams{
		QrCodeUrl:    qrUrl,
		QrCodeWidth:  qrWidth,
		QrCodeHeight: qrHeight,
		QrCodeX:      qrX,
		QrCodeY:      qrY,
		BgUrl:        bgUrl,
	}
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
