package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Enemy struct {
	Mecha
	Health int
	Armor  int
	Ammo   int
	Charge int
	Level  int
}

// enemy attack patterns below
var cursor int = 0 // for iterating through attack pattern
var attack int = -1

// united nations enemies
var une1 = [3]string{PARRYSTATE, MELEESTATE, GUARDSTATE}
var une2 = [3]string{MELEESTATE, MELEESTATE, MELEESTATE}
var une3 = [3]string{MELEESTATE, GUARDSTATE, MELEESTATE}

// grand rev enemies
var gre1 = [3]string{MELEESTATE, PARRYSTATE, GUARDSTATE}
var gre2 = [3]string{PARRYSTATE, MISSILESTATE, GUARDSTATE}
var gre3 = [3]string{GUARDSTATE, GUARDSTATE, MELEESTATE}

// imperium enemies
var ime1 = [3]string{GUARDSTATE, PARRYSTATE, MISSILESTATE}
var ime2 = [3]string{MELEESTATE, MELEESTATE, GUARDSTATE}
var ime3 = [3]string{MELEESTATE, PARRYSTATE, GUARDSTATE}

// collective enemies
var coe1 = [3]string{GUARDSTATE, PARRYSTATE, LASERSTATE}
var coe2 = [3]string{MELEESTATE, PARRYSTATE, MELEESTATE}
var coe3 = [3]string{MELEESTATE, GUARDSTATE, GUARDSTATE}

// empire boss
var boss1 = [3]string{GUARDSTATE, PARRYSTATE, LASERSTATE}
var boss2 = [3]string{MELEESTATE, PARRYSTATE, MISSILESTATE}
var boss3 = [3]string{MELEESTATE, MELEESTATE, MELEESTATE}
var boss4 = [3]string{PARRYSTATE, GUARDSTATE, LASERSTATE}
var boss5 = [3]string{PARRYSTATE, MELEESTATE, MISSILESTATE}

func createBasicEnemy(mechSheets []rl.Texture2D, level int) Enemy {
	CreatingEnemy = true
	m := newEnemyMecha(rl.Vector2Zero(), mechSheets)
	CreatingEnemy = false
	e := Enemy{
		Mecha:  m,
		Health: 100,
		Armor:  100 * level,
		Ammo:   50,
		Charge: 50,
		Level:  level,
	}
	return e
}

// generic mecha constructor, for player and regular enemies
func newEnemyMecha(newPos rl.Vector2, spriteSheets []rl.Texture2D) Mecha {
	// create transform
	newTransform := NewTransform(newPos)

	// create animations here
	newIdleAnimation := NewAnimation(newTransform, spriteSheets[0], 0.25, IDLESTATE)
	newGuardAnimation := NewAnimation(newTransform, spriteSheets[1], 0.25, GUARDSTATE)
	newMeleeAnimation := NewAnimation(newTransform, spriteSheets[2], 0.25, LASERSTATE)
	newParryAnimation := NewAnimation(newTransform, spriteSheets[3], 0.25, MELEESTATE)
	newLaserAnimation := NewAnimation(newTransform, spriteSheets[4], 0.25, MISSILESTATE)
	newMissileAnimation := NewAnimation(newTransform, spriteSheets[5], 0.25, PARRYSTATE)

	// create mecha struct variable
	newMecha := Mecha{
		Transform:             newTransform,
		AnimationStateMachine: NewStateMachine(&newIdleAnimation),
	}

	// add animations here
	newMecha.AnimationStateMachine.AddState(&newGuardAnimation)
	newMecha.AnimationStateMachine.AddState(&newMeleeAnimation)
	newMecha.AnimationStateMachine.AddState(&newParryAnimation)
	newMecha.AnimationStateMachine.AddState(&newLaserAnimation)
	newMecha.AnimationStateMachine.AddState(&newMissileAnimation)
	if !CreatingEnemy {
		newCoreSwitchAnimation := NewAnimation(newTransform, spriteSheets[6], 0.25, CORESTATE) // won't be used by enemies
		newMecha.AnimationStateMachine.AddState(&newCoreSwitchAnimation)                       // won't be used by enemies
	}
	// return mecha
	return newMecha
}
