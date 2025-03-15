package ui

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

func Start() error {
	ebiten.SetWindowSize(900, 800)
	ebiten.SetWindowTitle("Build Downloader")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	root := widget.NewContainer()
	eui := &ebitenui.UI{
		Container: root,
	}
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}
	fontFace := &text.GoTextFace{
		Source: s,
		Size:   32,
	}
	helloWorldLabel := widget.NewText(
		widget.TextOpts.Text("Hello World!", fontFace, color.White),
	)
	root.AddChild(helloWorldLabel)

	render := &renderer{
		ui: eui,
	}

	if err := ebiten.RunGame(render); err != nil {
		return fmt.Errorf("ebiten.RunGame: %w", err)
	}

	return nil
}