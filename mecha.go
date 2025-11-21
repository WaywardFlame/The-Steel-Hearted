package main

import rl "github.com/gen2brain/raylib-go/raylib"

// a mecha struct
// to be embedded in player, enemies, and bosses
type Mecha struct {
	*Transform
	AnimationStateMachine StateMachine
}

// generic mecha constructor, for player and regular enemies
func newMecha(newPos rl.Vector2, spriteSheets []rl.Texture2D) Mecha {
	// create transform
	newTransform := NewTransform(newPos)

	// create animations here
	newIdleAnimation := NewAnimation(newTransform, spriteSheets[0], 0.25, IDLESTATE)
	newGuardAnimation := NewAnimation(newTransform, spriteSheets[1], 0.25, GUARDSTATE)
	newMeleeAnimation := NewAnimation(newTransform, spriteSheets[2], 0.25, MELEESTATE)
	newParryAnimation := NewAnimation(newTransform, spriteSheets[3], 0.25, PARRYSTATE)
	newLaserAnimation := NewAnimation(newTransform, spriteSheets[4], 0.25, LASERSTATE)
	newMissileAnimation := NewAnimation(newTransform, spriteSheets[5], 0.25, MISSILESTATE)

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
