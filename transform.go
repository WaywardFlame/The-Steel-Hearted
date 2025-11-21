package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Transform struct {
	Pos   rl.Vector2
	Scale int
}

func NewTransform(newPos rl.Vector2) *Transform {
	return &Transform{Pos: newPos, Scale: 1}
}
