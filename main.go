package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/x/explorer"
	"github.com/zawpavel/qrcp-desktop/browser"
	"github.com/zawpavel/qrcp-desktop/filetransfer"
	"github.com/zawpavel/qrcp-desktop/graphic"
)

func main() {
	go func() {
		w := app.NewWindow(
			app.Title("qrcp File transfer"),
			app.Size(unit.Dp(400), unit.Dp(600)),
		)
		if err := draw(w); err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()
	app.Main()
}

var textBody = graphic.InitialTextBody
var qrImage = graphic.EmptyImageWidget

func draw(w *app.Window) error {
	var ops op.Ops
	var sendButton widget.Clickable
	var receiveButton widget.Clickable
	var donateButton widget.Clickable
	expl := explorer.NewExplorer(w)

	for e := range w.Events() {
		expl.ListenEvents(e)
		switch e := e.(type) {
		case system.FrameEvent:
			if sendButton.Clicked() {
				go processSendButton(expl)
			}
			if receiveButton.Clicked() {
				go processReceiveButton(w)
			}
			if donateButton.Clicked() {
				go processDonateButton()
			}
			gtx := layout.NewContext(&ops, e)
			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,
				graphic.LayoutButton(&sendButton, "Send"),
				graphic.LayoutButton(&receiveButton, "Receive"),
				graphic.LayoutText(textBody),
				graphic.LayoutImage(&qrImage),
				graphic.LayoutButton(&donateButton, "Donate"),
			)
			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}

func processSendButton(expl *explorer.Explorer) {
	textBody = graphic.ProcessingTextBody
	qrImage = graphic.EmptyImageWidget
	downloadLink := filetransfer.SendFiles(expl)
	if downloadLink == "" {
		textBody = graphic.InitialTextBody
	}
	if downloadLink != "" {
		textBody = graphic.HintTextBody
		qrImage = graphic.CreateQrImage(downloadLink)
	}
}

func processReceiveButton(w *app.Window) {
	receiveUrlChannel := make(chan string)
	outputDirChannel := make(chan string)
	go filetransfer.ReceiveFiles(receiveUrlChannel, outputDirChannel)
	receiveUrl := <-receiveUrlChannel
	textBody = graphic.HintTextBody
	qrImage = graphic.CreateQrImage(receiveUrl)

	// wait until file is received
	outDir := <-outputDirChannel
	textBody = graphic.CreateOutputDirHint(outDir)
	qrImage = graphic.EmptyImageWidget
	w.Invalidate()
}

func processDonateButton() {
	browser.Open("https://www.buymeacoffee.com/zawpavel")
}
