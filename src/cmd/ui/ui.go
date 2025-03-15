package ui

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func Start() error {
	// ebiten.SetWindowDecorated(false)
	ebiten.SetWindowSize(900, 800)
	ebiten.SetWindowTitle("Build Downloader")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	root := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x22, 0x22, 0x22, 0xff})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Direction(widget.DirectionVertical))),
	)

	eui := &ebitenui.UI{
		Container: root,
	}
	s, err := text.NewGoTextFaceSource(bytes.NewReader(Geist))
	if err != nil {
		log.Fatal(err)
	}
	s2, err := text.NewGoTextFaceSource(bytes.NewReader(Binance))
	if err != nil {
		log.Fatal(err)
	}
	fontFace := &text.GoTextFace{
		Source: s,
		Size:   24,
	}
	fontFace2 := &text.GoTextFace{
		Source: s2,
		Size:   14,
	}
	helloWorldLabel := widget.NewText(
		widget.TextOpts.Text("Fortnite Archive", fontFace, color.NRGBA{0xd4, 0xd4, 0xd4, 0xff}),
	)
	root.AddChild(helloWorldLabel)
	helloWorldLabel2 := widget.NewText(
		widget.TextOpts.Text("Hello World!", fontFace2, color.NRGBA{0xd4, 0xd4, 0xd4, 0xff}),
	)
	root.AddChild(helloWorldLabel2)

	render := &renderer{
		ui: eui,
	}

	if err := ebiten.RunGame(render); err != nil {
		return fmt.Errorf("ebiten.RunGame: %w", err)
	}

	return nil
}