package services

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/requiem-ai/discord-dynamic-sig/server"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"net/http"
)

type SignatureService struct {
	DefaultService
}

const SIGNATURE_SVC = "signature_svc"

func (svc SignatureService) Id() string {
	return SIGNATURE_SVC
}

func (svc *SignatureService) Start() error {
	return nil
}

func (svc *SignatureService) Generate(detail *server.DiscordInvite) ([]byte, error) {
	dc := gg.NewContext(430, 70)

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 32})
	faceSmall := truetype.NewFace(font, &truetype.Options{Size: 18})

	serverImg, err := svc.DownloadImage(detail.ServerAvatar(64))
	if err != nil {
		return nil, err
	}

	dc.SetRGB255(36, 38, 43) //Color
	dc.Clear()

	dc.DrawImage(serverImg, 3, 3)

	dc.SetFontFace(faceSmall)
	dc.SetRGB255(255, 255, 255) //Color
	dc.DrawString(detail.Guild.Name, 85, 25)
	dc.DrawString(fmt.Sprintf("%v Online", detail.ApproximatePresenceCount), 85, 50)
	dc.DrawString(fmt.Sprintf("%v Members", detail.ApproximateMemberCount), 190, 50)

	dc.DrawRoundedRectangle(320, 7, 100, 50, 5)
	dc.SetRGB255(44, 140, 105) //Color
	dc.Fill()

	dc.SetFontFace(face)
	dc.SetRGB255(255, 255, 255) //Color
	dc.DrawString("Join", 340, 42)

	var buf bytes.Buffer
	err = dc.EncodePNG(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (svc *SignatureService) DownloadImage(imgPath string) (image.Image, error) {
	resp, err := http.Get(imgPath)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}
