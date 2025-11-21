package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ClickBound interface {
	OnClick()
}

type Button struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
	ColorTheme
	text         string
	textSize     int32
	mouseOver    bool
	mouseDown    bool
	mouseClicked bool
	onClickFuncs []func()
}

func NewButton(x, y, width, height int32, newTheme ColorTheme) Button {
	nb := Button{X: x, Y: y, Width: width, Height: height, ColorTheme: newTheme}
	nb.onClickFuncs = make([]func(), 0)
	nb.mouseDown = false
	nb.mouseClicked = false
	nb.mouseDown = false
	return nb
}

func (b *Button) SetText(newText string, newTextSize int32) {
	b.text = newText
	b.textSize = newTextSize
}

func (b *Button) DrawButton() {
	currentColor := b.baseColor

	if b.mouseOver {
		currentColor = rl.NewColor(b.baseColor.R/2, b.baseColor.G/2, b.baseColor.B/2, 255)
	}
	if b.mouseDown {
		currentColor = b.accentColor
	}

	rl.DrawRectangle(b.X, b.Y, b.Width, b.Height, currentColor)

	textWidth := rl.MeasureText(b.text, int32(b.textSize))
	rl.DrawText(
		b.text,
		b.X+(b.Width/2)-(textWidth/2),
		b.Y+(b.Height/2)-(b.textSize/2),
		int32(b.textSize),
		b.textColor,
	)
}

func (b *Button) CenterButtonX() {
	b.X = int32(rl.GetScreenWidth()/2) - b.Width/2
}

func (b *Button) CenterButtonY() {
	b.Y = int32(rl.GetScreenHeight()/2) - b.Height/2
}

func (b *Button) CenterButton() {
	b.CenterButtonX()
	b.CenterButtonY()
}

func (b *Button) CheckMouseOver() {
	mousePos := rl.GetMousePosition()
	b.mouseOver = false
	if int32(mousePos.X) < b.X || int32(mousePos.X) > b.X+b.Width {
		return
	}
	if int32(mousePos.Y) < b.Y || int32(mousePos.Y) > b.Y+b.Height {
		return
	}
	b.mouseOver = true
}

func (b *Button) CheckMouseClick(sound rl.Sound) {

	b.mouseDown = false

	if b.mouseOver && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		b.mouseClicked = true
		for i := 0; i < len(b.onClickFuncs); i++ {
			b.onClickFuncs[i]() //call the funcs
			rl.PlaySound(sound)
		}
	}

	if b.mouseClicked && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		b.mouseDown = true
	}

	if !b.mouseOver || !rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		b.mouseClicked = false
	}
}

func (b *Button) AddOnClickFunc(newFunc func()) {
	b.onClickFuncs = append(b.onClickFuncs, newFunc)
}

func (b *Button) UpdateButton(sound rl.Sound) {

	b.CheckMouseOver()

	b.DrawButton()

	b.CheckMouseClick(sound)

}
