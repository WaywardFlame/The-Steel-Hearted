package main

import (
	"encoding/json"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// a struct for holding player mech data and information
type PlayerMech struct {
	Mecha
	MechData
	NumHeals int
}

// for saving and loading player mech, needed cause textures can't be saved
type MechData struct {
	Health       int
	Armor        int
	Ammo         int
	Charge       int
	SupportCores [8]bool
	FrameUpgrade int
}

// a struct for holding general player data, contains a PlayerMech
type PlayerShip struct {
	Mech          PlayerMech
	Sprites       []rl.Texture2D
	CurrentSprite int
	ShipData
}

// for saving and loading player ship, needed cause textures can't be saved
type ShipData struct {
	MechShip      MechData
	Money         int
	CurrentXP     int
	NextLevelXP   int
	Level         int
	Position      rl.Vector2
	CurrentSystem int
	// 4 - bases, 3 - planets, 1 - boss, 1 - trade
	UnitedNationsChecks [9]bool
	GrandRevChecks      [9]bool
	ImperiumChecks      [9]bool
	CollectiveChecks    [9]bool
	GoldenChecks        [9]bool
}

// moves the player ship with camera, for overworld exploration
func (p *PlayerShip) movePlayerShip() {
	dir := rl.NewVector2(0, 0)
	if rl.IsKeyDown(rl.KeyW) {
		dir.Y -= 1
		p.CurrentSprite = 1
	}
	if rl.IsKeyDown(rl.KeyA) {
		dir.X -= 1
		p.CurrentSprite = 3
	}
	if rl.IsKeyDown(rl.KeyS) {
		dir.Y += 1
		p.CurrentSprite = 2
	}
	if rl.IsKeyDown(rl.KeyD) {
		dir.X += 1
		p.CurrentSprite = 0
	}
	p.Position = rl.Vector2Add(p.Position, rl.Vector2Scale(dir, 400*rl.GetFrameTime()))
}

// draws the player ship
func (p *PlayerShip) drawPlayerShip() {
	i := p.CurrentSprite
	sourceRect := rl.NewRectangle(0, 0, float32(p.Sprites[i].Width), float32(p.Sprites[i].Height))
	destRect := rl.NewRectangle(p.Position.X, p.Position.Y, float32(p.Sprites[i].Width)*1, float32(p.Sprites[i].Height)*1)
	origin := rl.Vector2Scale(rl.NewVector2(float32(p.Sprites[i].Width)/2, float32(p.Sprites[i].Height)/2), 1)
	rl.DrawTexturePro(p.Sprites[i], sourceRect, destRect, origin, 0, rl.White)
}

// creates a default player mech without textures
func createEmptyPlayerMechData() MechData {
	p := MechData{
		Health:       100,
		Armor:        100,
		Ammo:         100,
		Charge:       0,
		SupportCores: [8]bool{true, true, true},
		FrameUpgrade: 1,
	}
	return p
}

// creates a default player ship and data without textures
func createEmptyPlayerShipData() ShipData {
	p := ShipData{
		MechShip:      createEmptyPlayerMechData(),
		Money:         0,
		CurrentXP:     0,
		NextLevelXP:   100,
		Level:         1,
		Position:      rl.NewVector2(0, 0),
		CurrentSystem: 1,
	}
	return p
}

func copyDataOver(s ShipData, mechSheets []rl.Texture2D, shipSheet []rl.Texture2D) PlayerShip {
	m := PlayerMech{
		Mecha:    newMecha(rl.NewVector2(0, 0), mechSheets),
		MechData: s.MechShip,
	}
	p := PlayerShip{
		Mech:     m,
		Sprites:  shipSheet,
		ShipData: s,
	}
	return p
}

// function for saving the game
func (player *ShipData) SavePlayer() error {
	data, err := json.MarshalIndent(player, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("data/steelheartedsave.json", data, 0644)
}

// function for loading the game
func (player *ShipData) LoadPlayer() error {
	data, err := os.ReadFile("data/steelheartedsave.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, player)
}
