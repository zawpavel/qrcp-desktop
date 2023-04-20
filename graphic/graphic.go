package graphic

import (
	"image"
	"log"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/skip2/go-qrcode"
)

type C = layout.Context
type D = layout.Dimensions

var EmptyImageWidget = widget.Image{
	Src:      paint.NewImageOp(image.NewRGBA(image.Rect(0, 0, 256, 256))),
	Scale:    1,
	Position: layout.Center,
}

var buttonMargins = layout.Inset{
	Top:    unit.Dp(25),
	Bottom: unit.Dp(25),
	Right:  unit.Dp(35),
	Left:   unit.Dp(35),
}

var theme = material.NewTheme(gofont.Collection())

var InitialTextBody = createTextBody("  This device and the mobile device must be connected to the same WiFi")
var HintTextBody = createTextBody("  Scan the QR code with your mobile device\n")
var ProcessingTextBody = createTextBody("  processing...\n")

func LayoutButton(buttonWidget *widget.Clickable, buttonText string) layout.FlexChild {
	return layout.Rigid(
		func(gtx C) D {
			return buttonMargins.Layout(gtx,
				func(gtx C) D {
					btn := material.Button(theme, buttonWidget, buttonText)
					return btn.Layout(gtx)
				},
			)
		},
	)
}

func LayoutImage(qrImage *widget.Image) layout.FlexChild {
	return layout.Rigid(
		func(gtx C) D {
			return qrImage.Layout(gtx)
		},
	)
}

func LayoutText(textBody material.LabelStyle) layout.FlexChild {
	return layout.Rigid(
		func(gtx C) D {
			return textBody.Layout(gtx)
		},
	)
}

func CreateQrImage(s string) widget.Image {
	q, err := qrcode.New(s, qrcode.Medium)
	if err != nil {
		log.Fatal(err)
	}
	img := q.Image(256)
	return widget.Image{
		Src:      paint.NewImageOp(img),
		Scale:    1,
		Position: layout.Center,
	}
}

func createTextBody(s string) material.LabelStyle {
	body := material.H6(theme, s)
	body.Alignment = text.Middle
	body.MaxLines = 2
	return body
}

func CreateOutputDirHint(s string) material.LabelStyle {
	return createTextBody("File uploaded to " + s)
}
