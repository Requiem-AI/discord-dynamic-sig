package services

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/requiem-ai/discord-dynamic-sig/server"
	"golang.org/x/image/font/gofont/gobold"
	"image"
	"net/http"
)

type SignatureService struct {
	DefaultService
}

type Signature struct {
	ShowInviteMessage bool    `json:"show_invite_message"`
	ButtonText        string  `json:"button_text"`
	ButtonColor       string  `json:"button_color"`
	BackgroundColor   string  `json:"background_color"`
	ButtonTextColor   string  `json:"button_text_color"`
	TitleTextColor    string  `json:"title_text_color"`
	InfoTextColor     string  `json:"info_text_color"`
	TitleSize         float64 `json:"title_size"`
	InfoSize          float64 `json:"info_size"`
	ButtonTextSize    float64 `json:"button_text_size"`
}

const SIGNATURE_SVC = "signature_svc"

func (svc SignatureService) Id() string {
	return SIGNATURE_SVC
}

func (svc *SignatureService) Start() error {
	return nil
}

func (svc *SignatureService) Generate(detail *server.DiscordInvite, sig *Signature) ([]byte, error) {
	dc := gg.NewContext(432, 70)

	//font, err := truetype.Parse(goregular.TTF)
	//if err != nil {
	//	return nil, err
	//}

	boldFont, err := truetype.Parse(gobold.TTF)
	if err != nil {
		return nil, err
	}

	face := truetype.NewFace(boldFont, &truetype.Options{Size: sig.TitleSize})
	faceSmall := truetype.NewFace(boldFont, &truetype.Options{Size: sig.InfoSize})
	faceButton := truetype.NewFace(boldFont, &truetype.Options{Size: sig.ButtonTextSize})

	serverImg, err := svc.DownloadImage(detail.ServerAvatar(64))
	if err != nil {
		return nil, err
	}

	dc.SetHexColor(sig.BackgroundColor)
	//dc.SetRGB255(36, 38, 43) //Color
	dc.Clear()

	dc.DrawImage(serverImg, 3, 3) // 3 + 64 + 16 margin3+64

	dc.SetFontFace(face)
	dc.SetHexColor(sig.TitleTextColor) //Color
	dc.DrawString(detail.Guild.Name, 82, 25)

	dc.DrawCircle(85, 45, 4)
	dc.SetHexColor("#3ba55d")
	dc.Fill()

	dc.SetFontFace(faceSmall)
	dc.SetHexColor(sig.InfoTextColor) //Color
	dc.DrawString(fmt.Sprintf("%v Online", detail.ApproximatePresenceCount), 93, 50)

	x := 180.0
	if detail.ApproximatePresenceCount > 1000 {
		x = x + 6 //Bump along a little bit more
	}

	dc.DrawCircle(x, 45, 4)
	dc.SetHexColor("#747f8d")
	dc.Fill()

	dc.SetFontFace(faceSmall)
	dc.SetHexColor(sig.InfoTextColor) //Color
	dc.DrawString(fmt.Sprintf("%v Members", detail.ApproximateMemberCount), x+8, 50)

	dc.DrawRoundedRectangle(350, 15, 64, 40, 5)
	dc.SetHexColor(sig.ButtonColor) //Color
	dc.Fill()

	dc.SetFontFace(faceButton)
	dc.SetHexColor(sig.ButtonTextColor) //Color
	dc.DrawStringAnchored("Join", 382, 32.5, 0.5, 0.5)

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

func (svc *SignatureService) DiscordDefaultSignature() *Signature {
	return &Signature{
		ShowInviteMessage: false,
		ButtonText:        "Join",
		ButtonTextColor:   "#FFFFFF",
		TitleTextColor:    "#FFFFFF",
		InfoTextColor:     "#B9BBBE",
		ButtonColor:       "#3BA55D",
		BackgroundColor:   "#2F3136",
		TitleSize:         16,
		InfoSize:          14,
		ButtonTextSize:    14,
	}
}

//Online dot = 3ba55d
//Member dot = 7474f8d
