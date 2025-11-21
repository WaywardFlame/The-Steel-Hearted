package main

import rl "github.com/gen2brain/raylib-go/raylib"

type ColorTheme struct {
	baseColor   rl.Color
	accentColor rl.Color
	textColor   rl.Color
}

func NewColorTheme(base, accent, text rl.Color) ColorTheme {
	ct := ColorTheme{
		baseColor:   base,
		accentColor: accent,
		textColor:   text,
	}
	return ct
}
