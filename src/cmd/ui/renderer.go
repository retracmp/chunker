package ui

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type renderer struct {
	ui *ebitenui.UI

	btn *widget.Button
}

func (g *renderer) Update() error {
	g.ui.Update()
	return nil
}

func (g *renderer) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)
}

func (g *renderer) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}