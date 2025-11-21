package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// struct for planetary bodies and stars
type Body struct {
	position   rl.Vector2
	radius     float32 // for collision purposes
	texture    rl.Texture2D
	isBase     bool
	isCleared  bool
	isCapital  bool
	isDefeated bool
	isTrade    bool
}

// draws a body
func (b *Body) DrawBody() {
	sourceRect := rl.NewRectangle(0, 0, float32(b.texture.Width), float32(b.texture.Height))
	destRect := rl.NewRectangle(b.position.X, b.position.Y, float32(b.texture.Width)*1, float32(b.texture.Height)*1)
	origin := rl.Vector2Scale(rl.NewVector2(float32(b.texture.Width)/2, float32(b.texture.Height)/2), 1)
	rl.DrawTexturePro(b.texture, sourceRect, destRect, origin, 0, rl.White)
}

// player globals
var PlayerLoaded bool = false
var PlayerPaused bool = false
var Camera rl.Camera2D = rl.NewCamera2D(
	rl.NewVector2(0, 0),
	rl.NewVector2(0, 0),
	0,
	1,
)
var CanInteract bool = false
var CurrentInteractable int = 0
var isInteracting bool = false
var InteractionSystem int

// battle globals
var OldSystem int = 1
var CurrencyReward int = 0
var XPReward int = 0
var Opponent Enemy
var CreatingEnemy bool = false
var PlayerCore int = 0

// pause menu globals
var pResumeGameButton Button
var pSaveGameButton Button
var pReturnToTitleButton Button
var pQuitGameButton Button
var pGamePaused bool = false
var pHasSaved bool = false
var pAreYouSure bool = false
var pQuitPressed bool = false

// trade menu globals
var rMessageButton Button
var rFrameUpgradeButton Button
var rSupportCoreButton Button
var rExitTradingButton Button

// title screen globals
var tNewGameButton Button
var tLoadGameButton Button
var tQuitButton Button
var tTitleBox Button
var tAreYouSure bool = false

// body positions
var StarPosition rl.Vector2 = rl.NewVector2(0, 0) // center of system
var Planet1 rl.Vector2 = rl.NewVector2(320, 0)    // to the right of star
var Planet2 rl.Vector2 = rl.NewVector2(0, 448)    // to the bottom of star
var Planet3 rl.Vector2 = rl.NewVector2(0, -576)   // to the top of star
var Planet4 rl.Vector2 = rl.NewVector2(-704, 0)   // to the left of star
var Planet5 rl.Vector2 = rl.NewVector2(           // to the bottom-right of star
	float32((0 + 832*math.Cos((135*math.Pi)/180))),
	float32((0 + 832*math.Sin((135*math.Pi)/180))),
)
var Planet6 rl.Vector2 = rl.NewVector2( // to the upper-left of star
	float32((0 + -960*math.Cos((225*math.Pi)/180))),
	float32((0 + -960*math.Cos((225*math.Pi)/180))),
)
var Planet7 rl.Vector2 = rl.NewVector2( // to upper-right of star
	float32((0 + 1088*math.Cos((45*math.Pi)/180))),
	float32((0 + -1088*math.Sin((45*math.Pi)/180))),
)
var Planet8 rl.Vector2 = rl.NewVector2( // to the bottom-left of star
	float32((0 + -1216*math.Cos((315*math.Pi)/180))),
	float32((0 + 1216*math.Sin((315*math.Pi)/180))),
)

// united nations of earth globals
var UnitedSystem []Body = make([]Body, 0, 9)

func adjustUnitedSystem(checks [9]bool) {
	for i := 0; i < len(UnitedSystem); i++ {
		UnitedSystem[i].isBase = false
		UnitedSystem[i].isCleared = false
		UnitedSystem[i].isCapital = false
		UnitedSystem[i].isDefeated = false
		UnitedSystem[i].isTrade = false
	}

	// mercury - base 1
	UnitedSystem[1].isBase = true
	UnitedSystem[1].isCleared = checks[0]

	// venus - base 2
	UnitedSystem[2].isBase = true
	UnitedSystem[2].isCleared = checks[1]

	// earth - capital
	UnitedSystem[3].isCapital = true
	UnitedSystem[3].isDefeated = checks[7]
	UnitedSystem[3].isTrade = checks[8]

	// mars - base 3
	UnitedSystem[4].isBase = true
	UnitedSystem[4].isCleared = checks[2]

	// jupiter - base 4
	UnitedSystem[5].isBase = true
	UnitedSystem[5].isCleared = checks[3]
}

// grand revolutionary system
var GrandSystem []Body = make([]Body, 0, 9)

func adjustGrandSystem(checks [9]bool) {
	for i := 0; i < len(GrandSystem); i++ {
		GrandSystem[i].isBase = false
		GrandSystem[i].isCleared = false
		GrandSystem[i].isCapital = false
		GrandSystem[i].isDefeated = false
		GrandSystem[i].isTrade = false
	}

	// capital
	GrandSystem[1].isCapital = true
	GrandSystem[1].isDefeated = checks[7]
	GrandSystem[1].isTrade = checks[8]

	// base 1
	GrandSystem[2].isBase = true
	GrandSystem[2].isCleared = checks[0]

	// base 2
	GrandSystem[3].isBase = true
	GrandSystem[3].isCleared = checks[1]

	// base 3
	GrandSystem[4].isBase = true
	GrandSystem[4].isCleared = checks[2]

	// base 4
	GrandSystem[5].isBase = true
	GrandSystem[5].isCleared = checks[3]
}

// imperium system
var ImperiumSystem []Body = make([]Body, 0, 9)

func adjustImperiumSystem(checks [9]bool) {
	for i := 0; i < len(ImperiumSystem); i++ {
		ImperiumSystem[i].isBase = false
		ImperiumSystem[i].isCleared = false
		ImperiumSystem[i].isCapital = false
		ImperiumSystem[i].isDefeated = false
		ImperiumSystem[i].isTrade = false
	}

	// capital
	ImperiumSystem[1].isCapital = true
	ImperiumSystem[1].isDefeated = checks[7]
	ImperiumSystem[1].isTrade = checks[8]

	// base 1
	ImperiumSystem[2].isBase = true
	ImperiumSystem[2].isCleared = checks[0]

	// base 2
	ImperiumSystem[3].isBase = true
	ImperiumSystem[3].isCleared = checks[1]

	// base 3
	ImperiumSystem[4].isBase = true
	ImperiumSystem[4].isCleared = checks[2]

	// base 4
	ImperiumSystem[5].isBase = true
	ImperiumSystem[5].isCleared = checks[3]
}

// collective system
var CollectiveSystem []Body = make([]Body, 0, 9)

func adjustCollectiveSystem(checks [9]bool) {
	for i := 0; i < len(CollectiveSystem); i++ {
		CollectiveSystem[i].isBase = false
		CollectiveSystem[i].isCleared = false
		CollectiveSystem[i].isCapital = false
		CollectiveSystem[i].isDefeated = false
		CollectiveSystem[i].isTrade = false
	}

	// capital
	CollectiveSystem[1].isCapital = true
	CollectiveSystem[1].isDefeated = checks[7]
	CollectiveSystem[1].isTrade = checks[8]

	// base 1
	CollectiveSystem[2].isBase = true
	CollectiveSystem[2].isCleared = checks[0]

	// base 2
	CollectiveSystem[3].isBase = true
	CollectiveSystem[3].isCleared = checks[1]

	// base 3
	CollectiveSystem[4].isBase = true
	CollectiveSystem[4].isCleared = checks[2]

	// base 4
	CollectiveSystem[5].isBase = true
	CollectiveSystem[5].isCleared = checks[3]
}

// golden system
var GoldenSystem []Body = make([]Body, 0, 9)

func adjustGoldenSystem(checks [9]bool) {
	for i := 0; i < len(GoldenSystem); i++ {
		GoldenSystem[i].isBase = false
		GoldenSystem[i].isCleared = false
		GoldenSystem[i].isCapital = false
		GoldenSystem[i].isDefeated = false
		GoldenSystem[i].isTrade = false
	}

	// capital
	GoldenSystem[1].isCapital = true
	GoldenSystem[1].isDefeated = checks[7]
	GoldenSystem[1].isTrade = true
}

// end screen variables
var endScreenTitle Button
var endScreenMessage Button
