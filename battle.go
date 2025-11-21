package main

import (
	"math/rand/v2"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// player animation checks
var pIsIdle bool = true
var pIsGuard bool = false
var pIsMelee bool = false
var pIsLaser bool = false
var pIsMissile bool = false
var pCoreSwitch bool = false
var pIsMoving bool = false
var pCalculated bool = false

// enemy animation checks
var eIsIdle bool = true
var eIsGuard bool = false
var eIsMelee bool = false
var eIsLaser bool = false
var eIsMissile bool = false
var eIsParrying bool = false // not actually parrying, just to indicate next move
var eIsMoving bool = false
var eCalculated bool = false

// turns
var playerTurn bool = true
var enemyTurn bool = false

func reset() {
	// player animation checks
	pIsIdle = true
	pIsGuard = false
	pIsMelee = false
	pIsLaser = false
	pIsMissile = false
	pCoreSwitch = false
	pIsMoving = false

	// enemy animation checks
	eIsIdle = true
	eIsGuard = false
	eIsMelee = false
	eIsLaser = false
	eIsMissile = false
	eIsParrying = false // not actually parrying, just to indicate next move
	eIsMoving = false

	// turns
	playerTurn = true
	enemyTurn = false

	// battle globals
	OldSystem = 1
	CurrencyReward = 0
	XPReward = 0
	CreatingEnemy = false
	PlayerCore = 0
}

func Battle(player *PlayerMech, session *SessionData, gameData *GameData) {
	// set music
	if OldSystem == 5 {
		session.NewStream = 0
		rl.DrawTexture(gameData.BackgroundImages[1], 0, 0, rl.White)
	} else {
		session.NewStream = 3
		rl.DrawTexture(gameData.BackgroundImages[0], 0, 0, rl.White)
	}

	// player ui
	rl.DrawRectangle(100, 160, 200, 200, rl.White)
	rl.DrawText("Health: "+strconv.Itoa(player.Health), 120, 180, 20, rl.Black)
	rl.DrawText("Armor: "+strconv.Itoa(player.Armor), 120, 200, 20, rl.Black)
	rl.DrawText("Ammo: "+strconv.Itoa(player.Ammo), 120, 220, 20, rl.Black)
	rl.DrawText("Charge: "+strconv.Itoa(player.Charge), 120, 240, 20, rl.Black)
	rl.DrawText("Core:"+strconv.Itoa(PlayerCore), 120, 260, 20, rl.Black)

	// enemy ui
	rl.DrawRectangle(980, 160, 200, 200, rl.White)
	rl.DrawText("Health: "+strconv.Itoa(Opponent.Health), 1000, 180, 20, rl.Black)
	rl.DrawText("Armor: "+strconv.Itoa(Opponent.Armor), 1000, 200, 20, rl.Black)
	rl.DrawText("Level: "+strconv.Itoa(Opponent.Level), 1000, 220, 20, rl.Black)

	// set opponent position
	Opponent.Pos = rl.NewVector2(1080, 525)

	if playerTurn {
		// get player input
		playerInput(player, gameData)
	}
	if !pIsIdle && playerTurn {
		// perform calculations based on player action
		if !pCalculated {
			calculateImpactOnEnemy(player)
		}

		// do victory check
		victoryCheck(player, session, gameData)

		// change turns
		playerTurn = false
		enemyTurn = true
		player.Charge += 20
		if player.Charge > 100 {
			player.Charge = 100
		} else if player.Charge < 0 {
			player.Charge = 20
		}
	}
	playerAnimations(player)

	if enemyTurn {
		// perform enemy animations
		enemyInput(gameData)
	}
	if !eIsIdle && enemyTurn {
		// perform calculations based on enemy action
		if !eCalculated {
			calculateImpactOnPlayer(player)
		}

		// do victory check
		victoryCheck(player, session, gameData)

		enemyTurn = false
		playerTurn = true
	}
	enemyAnimations()

	drawCurrentMoves()
}

func drawCurrentMoves() {
	// player move
	if pIsGuard {
		rl.DrawText("Player using Guard", 400, 180, 40, rl.White)
	} else if pIsLaser {
		rl.DrawText("Player using Laser", 400, 180, 40, rl.White)
	} else if pIsMelee {
		rl.DrawText("Player using Melee", 400, 180, 40, rl.White)
	} else if pIsMissile {
		rl.DrawText("Player using Missile", 400, 180, 40, rl.White)
	} else if pCoreSwitch {
		rl.DrawText("Player switching Core", 400, 180, 40, rl.White)
	}

	// enemy move
	if eIsGuard {
		rl.DrawText("Enemy using Guard", 400, 240, 40, rl.White)
	} else if eIsLaser {
		rl.DrawText("Enemy using Laser", 400, 240, 40, rl.White)
	} else if eIsMelee {
		rl.DrawText("Enemy using Melee", 400, 240, 40, rl.White)
	} else if eIsMissile {
		rl.DrawText("Enemy using Missile", 400, 240, 40, rl.White)
	}
}

func victoryCheck(player *PlayerMech, session *SessionData, gameData *GameData) {
	if Opponent.Health <= 0 || Opponent.Armor <= 0 { // player has won battle
		// check for system, then check for interactable
		if OldSystem == 1 { // united nations
			session.player.Money += CurrencyReward
			if CurrentInteractable == 1 {
				session.player.UnitedNationsChecks[0] = true
			} else if CurrentInteractable == 2 {
				session.player.UnitedNationsChecks[1] = true
			} else if CurrentInteractable == 4 {
				session.player.UnitedNationsChecks[2] = true
			} else if CurrentInteractable == 5 {
				session.player.UnitedNationsChecks[3] = true
			}
			adjustUnitedSystem(session.player.UnitedNationsChecks)
		} else if OldSystem == 2 { // grand rev
			session.player.Money += CurrencyReward
			if CurrentInteractable == 2 {
				session.player.GrandRevChecks[0] = true
			} else if CurrentInteractable == 3 {
				session.player.GrandRevChecks[1] = true
			} else if CurrentInteractable == 4 {
				session.player.GrandRevChecks[2] = true
			} else if CurrentInteractable == 5 {
				session.player.GrandRevChecks[3] = true
			}
			adjustGrandSystem(session.player.GrandRevChecks)
		} else if OldSystem == 3 { // imperium system
			session.player.Money += CurrencyReward
			if CurrentInteractable == 2 {
				session.player.ImperiumChecks[0] = true
			} else if CurrentInteractable == 3 {
				session.player.ImperiumChecks[1] = true
			} else if CurrentInteractable == 4 {
				session.player.ImperiumChecks[2] = true
			} else if CurrentInteractable == 5 {
				session.player.ImperiumChecks[3] = true
			}
			adjustImperiumSystem(session.player.ImperiumChecks)
		} else if OldSystem == 4 { // collective system
			session.player.Money += CurrencyReward
			if CurrentInteractable == 2 {
				session.player.CollectiveChecks[0] = true
			} else if CurrentInteractable == 3 {
				session.player.CollectiveChecks[1] = true
			} else if CurrentInteractable == 4 {
				session.player.CollectiveChecks[2] = true
			} else if CurrentInteractable == 5 {
				session.player.CollectiveChecks[3] = true
			}
			adjustCollectiveSystem(session.player.CollectiveChecks)
		} else if OldSystem == 5 { // golden system
			session.player.Money += CurrencyReward
			if CurrentInteractable == 1 {
				session.player.GoldenChecks[7] = true
			}
			adjustGoldenSystem(session.player.GoldenChecks)
		}
		session.CurrentSystem = session.player.CurrentSystem
		reset() // reset necessary globals
		rl.PlaySound(gameData.BattleSounds[0])
	} else if player.Health <= 0 || player.Armor <= 0 { // player lost, return to last save
		session.LoadData()
		reset() // reset necessary globals
		rl.PlaySound(gameData.BattleSounds[0])
	}
}

func calculateImpactOnEnemy(player *PlayerMech) {
	if eIsGuard {
		if pIsLaser {
			pCalculated = true
			if PlayerCore == 6 { // decrease damage, decrease charge cost
				Opponent.Health -= 20
				Opponent.Armor -= (20 * player.FrameUpgrade) / 2
				player.Charge -= 20
				return
			} else if PlayerCore == 7 { // increase damage, increase charge cost
				Opponent.Health -= 40
				Opponent.Armor -= (40 * player.FrameUpgrade)
				player.Charge -= 40
				return
			} else { // normal damage and charge cost
				Opponent.Health -= 30
				Opponent.Armor -= (30 * player.FrameUpgrade)
				player.Charge -= 30
				return
			}
		} else if pIsMissile {
			pCalculated = true
			if PlayerCore == 4 { // decrease ammo cost by small amount
				player.Ammo -= 10
			} else if PlayerCore == 5 { // decrease ammo cost by large amount
				player.Ammo -= 5
			} else { // decrease ammo by normal amount
				player.Ammo -= 15
			}
			Opponent.Health -= 5
			Opponent.Armor -= (50 * player.FrameUpgrade)
			return
		} else {
			return
		}
	} else { // if enemy is not guarding
		if pIsLaser {
			pCalculated = true
			if PlayerCore == 6 { // decrease damage, decrease charge cost
				Opponent.Health -= 20
				Opponent.Armor -= (20 * player.FrameUpgrade) / 2
				player.Charge -= 20
				return
			} else if PlayerCore == 7 { // increase damage, increase charge cost
				Opponent.Health -= 40
				Opponent.Armor -= (40 * player.FrameUpgrade)
				player.Charge -= 40
				return
			} else { // normal damage and charge cost
				Opponent.Health -= 30
				Opponent.Armor -= (30 * player.FrameUpgrade)
				player.Charge -= 30
				return
			}
		} else if pIsMissile {
			pCalculated = true
			if PlayerCore == 4 { // decrease ammo cost by small amount
				player.Ammo -= 10
			} else if PlayerCore == 5 { // decrease ammo cost by large amount
				player.Ammo -= 5
			} else { // decrease ammo by normal amount
				player.Ammo -= 15
			}
			Opponent.Health -= 5
			Opponent.Armor -= (50 * player.FrameUpgrade)
			return
		} else if pIsMelee {
			pCalculated = true
			if PlayerCore == 3 { // charge steal, light armor damage
				Opponent.Armor -= (5 * player.FrameUpgrade)
				player.Charge += (5 * player.FrameUpgrade)
			} else if PlayerCore == 2 { // heavy damage to armor, 1 health damage
				Opponent.Health -= 1
				Opponent.Armor -= (25 * player.FrameUpgrade)
			} else if PlayerCore == 1 { // medium armor damage, little health damage
				Opponent.Armor -= (15 * player.FrameUpgrade)
				Opponent.Health -= 5
			} else if PlayerCore == 0 { // low armor damage, some health damage
				Opponent.Armor -= (5 * player.FrameUpgrade)
				Opponent.Health -= 10
			}
			return
		} else if pIsGuard || pCoreSwitch {
			pCalculated = true
			return
		}
	}
}

func calculateImpactOnPlayer(player *PlayerMech) {
	if eIsGuard || eIsParrying {
		eCalculated = true
		return
	} else {
		if pIsGuard {
			if eIsLaser {
				eCalculated = true
				player.Health -= (1 * Opponent.Level)
				player.Armor -= (5 * Opponent.Level)
				return
			} else if eIsMissile {
				eCalculated = true
				player.Health -= (1 * Opponent.Level)
				player.Armor -= (10 * Opponent.Level)
				return
			} else {
				return
			}
		} else {
			eCalculated = true
			if eIsLaser {
				player.Health -= (2 * Opponent.Level)
				player.Armor -= (6 * Opponent.Level)
				return
			} else if eIsMissile {
				player.Health -= (2 * Opponent.Level)
				player.Armor -= (11 * Opponent.Level)
				return
			} else if eIsMelee {
				player.Health -= (1 * Opponent.Level)
				player.Armor -= (15 * Opponent.Level)
				return
			}
		}
	}
}

func playerInput(player *PlayerMech, gameData *GameData) {
	player.Pos = rl.NewVector2(200, 525)

	if !pIsMoving {
		if rl.IsKeyPressed(rl.KeyOne) { // melee
			pIsMelee = true
			pIsIdle = false
			rl.PlaySound(gameData.BattleSounds[6])
		} else if rl.IsKeyPressed(rl.KeyTwo) { // guard
			pIsGuard = true
			pIsIdle = false
			rl.PlaySound(gameData.BattleSounds[4])
		} else if rl.IsKeyPressed(rl.KeyThree) { // laser
			pIsLaser = true
			pIsIdle = false
			rl.PlaySound(gameData.BattleSounds[5])
		} else if rl.IsKeyPressed(rl.KeyFour) { // missile
			pIsMissile = true
			pIsIdle = false
			rl.PlaySound(gameData.BattleSounds[8])
		} else if rl.IsKeyDown(rl.KeyFive) { // cores
			if rl.IsKeyPressed(rl.KeyKp1) {
				pCoreSwitch = true
				pIsIdle = false
				PlayerCore = 0
				rl.PlaySound(gameData.BattleSounds[2])
			} else if rl.IsKeyPressed(rl.KeyKp2) {
				pCoreSwitch = true
				pIsIdle = false
				PlayerCore = 1
				rl.PlaySound(gameData.BattleSounds[2])
			} else if rl.IsKeyPressed(rl.KeyKp3) {
				pCoreSwitch = true
				pIsIdle = false
				PlayerCore = 2
				rl.PlaySound(gameData.BattleSounds[2])
			} else if rl.IsKeyPressed(rl.KeyKp4) && player.SupportCores[3] {
				pCoreSwitch = true
				pIsIdle = false
				PlayerCore = 3
				rl.PlaySound(gameData.BattleSounds[2])
			} else if rl.IsKeyPressed(rl.KeyKp5) && player.SupportCores[4] {
				pCoreSwitch = true
				pIsIdle = false
				PlayerCore = 4
				rl.PlaySound(gameData.BattleSounds[2])
			} else if rl.IsKeyPressed(rl.KeyKp6) && player.SupportCores[5] {
				pCoreSwitch = true
				pIsIdle = false
				PlayerCore = 5
				rl.PlaySound(gameData.BattleSounds[2])
			} else if rl.IsKeyPressed(rl.KeyKp7) && player.SupportCores[6] {
				pCoreSwitch = true
				pIsIdle = false
				PlayerCore = 6
				rl.PlaySound(gameData.BattleSounds[2])
			} else if rl.IsKeyPressed(rl.KeyKp8) && player.SupportCores[7] {
				pCoreSwitch = true
				pIsIdle = false
				PlayerCore = 7
				rl.PlaySound(gameData.BattleSounds[2])
			}
		}
	}
}

func playerAnimations(player *PlayerMech) {
	// player animations
	if pIsIdle {
		pIsMoving = false
		player.AnimationStateMachine.ChangeState(IDLESTATE)
	} else if pIsMelee {
		pIsMoving = true
		player.AnimationStateMachine.ChangeState(MELEESTATE)
	} else if pIsGuard {
		pIsMoving = true
		player.AnimationStateMachine.ChangeState(GUARDSTATE)
	} else if pIsLaser {
		pIsMoving = true
		player.AnimationStateMachine.ChangeState(LASERSTATE)
	} else if pIsMissile {
		pIsMoving = true
		player.AnimationStateMachine.ChangeState(MISSILESTATE)
	} else if pCoreSwitch {
		pIsMoving = true
		player.AnimationStateMachine.ChangeState(CORESTATE)
	}
	player.AnimationStateMachine.Tick()
}

func enemyInput(gameData *GameData) {
	if OldSystem == 1 {
		choosePattern(3)
		if attack == 0 {
			enemyAttack(une1[cursor], gameData)
		} else if attack == 1 {
			enemyAttack(une2[cursor], gameData)
		} else if attack == 2 {
			enemyAttack(une3[cursor], gameData)
		}
		cursor += 1
	} else if OldSystem == 2 {
		choosePattern(3)
		if attack == 0 {
			enemyAttack(gre1[cursor], gameData)
		} else if attack == 1 {
			enemyAttack(gre2[cursor], gameData)
		} else if attack == 2 {
			enemyAttack(gre3[cursor], gameData)
		}
		cursor += 1
	} else if OldSystem == 3 {
		choosePattern(3)
		if attack == 0 {
			enemyAttack(ime1[cursor], gameData)
		} else if attack == 1 {
			enemyAttack(ime2[cursor], gameData)
		} else if attack == 2 {
			enemyAttack(ime3[cursor], gameData)
		}
		cursor += 1
	} else if OldSystem == 4 {
		choosePattern(3)
		if attack == 0 {
			enemyAttack(coe1[cursor], gameData)
		} else if attack == 1 {
			enemyAttack(coe2[cursor], gameData)
		} else if attack == 2 {
			enemyAttack(coe3[cursor], gameData)
		}
		cursor += 1
	} else if OldSystem == 5 {
		choosePattern(5)
		if attack == 0 {
			enemyAttack(boss1[cursor], gameData)
		} else if attack == 1 {
			enemyAttack(boss2[cursor], gameData)
		} else if attack == 2 {
			enemyAttack(boss3[cursor], gameData)
		} else if attack == 3 {
			enemyAttack(boss4[cursor], gameData)
		} else if attack == 4 {
			enemyAttack(boss5[cursor], gameData)
		}
		cursor += 1
	}
}

func choosePattern(num int) {
	pattern := rand.IntN(num)
	if cursor >= 3 {
		attack = -1
		cursor = 0
	}
	if attack == -1 {
		attack = pattern
	}
}

func enemyAttack(move string, gameData *GameData) {
	if move == MELEESTATE {
		eIsMelee = true
		eIsIdle = false
		rl.PlaySound(gameData.BattleSounds[6])
	} else if move == GUARDSTATE {
		eIsGuard = true
		eIsIdle = false
		rl.PlaySound(gameData.BattleSounds[4])
	} else if move == PARRYSTATE {
		eIsParrying = true
		eIsIdle = false
		rl.PlaySound(gameData.BattleSounds[9])
	} else if move == LASERSTATE {
		eIsLaser = true
		eIsIdle = false
		rl.PlaySound(gameData.BattleSounds[5])
	} else if move == MISSILESTATE {
		eIsMissile = true
		eIsIdle = false
		rl.PlaySound(gameData.BattleSounds[8])
	}
}

func enemyAnimations() {
	// enemy animations
	if eIsIdle {
		eIsMoving = false
		Opponent.AnimationStateMachine.ChangeState(IDLESTATE)
	} else if eIsMelee {
		eIsMoving = true
		Opponent.AnimationStateMachine.ChangeState(MELEESTATE)
	} else if eIsGuard {
		eIsMoving = true
		Opponent.AnimationStateMachine.ChangeState(GUARDSTATE)
	} else if eIsLaser {
		eIsMoving = true
		Opponent.AnimationStateMachine.ChangeState(LASERSTATE)
	} else if eIsMissile {
		eIsMoving = true
		Opponent.AnimationStateMachine.ChangeState(MISSILESTATE)
	} else if eIsParrying {
		eIsMoving = true
		Opponent.AnimationStateMachine.ChangeState(PARRYSTATE)
	}
	Opponent.AnimationStateMachine.Tick()
}
