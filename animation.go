package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	IDLESTATE    = "idle"
	GUARDSTATE   = "guard"
	MELEESTATE   = "melee"
	PARRYSTATE   = "parry"
	LASERSTATE   = "laser"
	MISSILESTATE = "missile"
	CORESTATE    = "core"
)

type Animation struct {
	*Transform
	SpriteSheet  rl.Texture2D
	MaxIndex     int
	CurrentIndex int
	Timer        float32
	SwitchTime   float32
	Name         string
}

func (a *Animation) TickState() {
	a.UpdateTime()
	a.DrawAnimation()
}

func (a *Animation) GetName() string {
	return a.Name
}

func (a *Animation) ResetTime() {
	a.Timer = 0
}

func NewAnimation(newTransform *Transform, newSheet rl.Texture2D, newTime float32, newName string) Animation {
	spriteDimension := newSheet.Height
	frames := int(newSheet.Width / spriteDimension)
	newAnimation := Animation{
		Transform:    newTransform,
		SpriteSheet:  newSheet,
		MaxIndex:     frames - 1,
		CurrentIndex: 0,
		Timer:        0,
		SwitchTime:   newTime,
		Name:         newName,
	}
	return newAnimation
}

func (a *Animation) UpdateTime() {
	a.Timer += rl.GetFrameTime()
	if a.Timer > a.SwitchTime {
		a.Timer = 0
		a.CurrentIndex++
	}

	if a.CurrentIndex > a.MaxIndex {
		a.CurrentIndex = 0
		if pIsMelee || eIsMelee {
			pIsMelee = false
			eIsMelee = false
		}
		if pIsGuard || eIsGuard {
			pIsGuard = false
			eIsGuard = false
		}
		if pIsLaser || eIsLaser {
			pIsLaser = false
			eIsLaser = false
		}
		if pIsMissile || eIsMissile {
			pIsMissile = false
			eIsMissile = false
		}
		if pCoreSwitch || eIsParrying {
			pCoreSwitch = false
			eIsParrying = false
		}
		if pIsMoving || eIsMoving {
			pIsMoving = false
			eIsMoving = false
			pIsIdle = true
			eIsIdle = true
			pCalculated = false
			eCalculated = false
		}
	}
}

func (a Animation) DrawAnimation() {
	sourceRect := rl.NewRectangle(float32(128*a.CurrentIndex), 0, 128, 128)
	destRect := rl.NewRectangle(a.Pos.X, a.Pos.Y, 128*float32(a.Scale), 128*float32(a.Scale))
	origin := rl.Vector2Scale(rl.NewVector2(float32(a.SpriteSheet.Height)/2, float32(a.SpriteSheet.Height)/2), float32(a.Scale))
	rl.DrawTexturePro(a.SpriteSheet, sourceRect, destRect, origin, 0, rl.White)
}
