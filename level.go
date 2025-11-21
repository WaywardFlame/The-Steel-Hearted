package main

import (
	"math"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// draws a given planetary system
func drawSystem(system []Body) {
	for i := 0; i < len(system); i++ {
		// get x and y to draw orbital lines, for player navigation
		x := system[i].position.X
		y := system[i].position.Y
		if int32(x) != 0 && int32(y) != 0 && i != 0 { // if distanced diagonally
			nx := int32(x)
			ny := int32(y)
			var radius float32
			if nx == 588 || ny == 588 { // if planet 5 position
				radius = 832
			} else if nx == 678 || ny == 678 { // if planet 6 position
				radius = 960
			} else if nx == 769 || ny == 769 { // if planet 7 position
				radius = 1088
			} else if nx == 859 || ny == 859 { // if planet 8 position
				radius = 1216
			}
			rl.DrawCircleLines(
				0,
				0,
				radius,
				rl.DarkGray,
			)
		} else if int32(x) != 0 && i != 0 { // if distanced horizontally
			x = float32(math.Abs(float64(x)))
			rl.DrawCircleLines(0, 0, x, rl.DarkGray)
		} else if int32(y) != 0 && i != 0 { // if distanced vertically
			y = float32(math.Abs(float64(y)))
			rl.DrawCircleLines(0, 0, y, rl.DarkGray)
		}
		// the body to draw
		system[i].DrawBody()
	}
}

func CheckCollision(system []Body, session *SessionData, gameData *GameData) {
	// more efficient to put inside draw system, but whatever
	for i := 0; i < len(system); i++ {
		if rl.Vector2Distance(system[i].position, session.player.Position) <= (system[i].radius * 2) {
			if system[i].isBase || system[i].isCapital {
				CanInteract = true
				CurrentInteractable = i
				rl.DrawTexture(gameData.MomentarySprites[2], int32(session.player.Position.X-32), int32(session.player.Position.Y-96), rl.White)
			} else {
				CanInteract = false
			}
		}
	}
}

// performs functionality shared by each system
func GenericLevel(session *SessionData, system []Body, gameData *GameData) {
	// set music
	session.NewStream = 1

	// handle player movement
	if !pGamePaused && !isInteracting { // check if game is paused
		session.player.movePlayerShip()
	}
	Camera.Target = rl.NewVector2(float32(int32(session.player.Position.X)), float32(int32(session.player.Position.Y)))
	Camera.Offset = rl.NewVector2(float32(rl.GetScreenWidth()/2), float32(rl.GetScreenHeight()/2))

	// draw player, planets, and star
	rl.BeginMode2D(Camera)
	rl.DrawCircleLines(0, 0, 1224, rl.DarkGray)
	drawSystem(system)
	session.player.drawPlayerShip()
	CheckCollision(system, session, gameData)
	rl.EndMode2D()

	if system[CurrentInteractable].isBase && system[CurrentInteractable].isCleared {
		rl.DrawText("Base has already been cleared.", 0, 140, 20, rl.White)
	} else if system[CurrentInteractable].isBase && !system[CurrentInteractable].isCleared {
		rl.DrawText("Base has not been cleared.", 0, 140, 20, rl.White)
	}

	// change between systems, resets position
	if !pGamePaused && !isInteracting { // check if game is paused
		if rl.IsKeyPressed(rl.KeyOne) { // united nations
			session.CurrentSystem = 1
			session.player.CurrentSystem = 1
			session.player.Position = rl.Vector2Zero()
		}
		if rl.IsKeyPressed(rl.KeyTwo) { // grand rev
			session.CurrentSystem = 2
			session.player.CurrentSystem = 2
			session.player.Position = rl.Vector2Zero()
		}
		if rl.IsKeyPressed(rl.KeyThree) { // imperium
			session.CurrentSystem = 3
			session.player.CurrentSystem = 3
			session.player.Position = rl.Vector2Zero()
		}
		if rl.IsKeyPressed(rl.KeyFour) { // collective
			session.CurrentSystem = 4
			session.player.CurrentSystem = 4
			session.player.Position = rl.Vector2Zero()
		}
		if rl.IsKeyPressed(rl.KeyFive) { // golden empire
			session.CurrentSystem = 5
			session.player.CurrentSystem = 5
			session.player.Position = rl.Vector2Zero()
		}
	}

	// move data over, in case player wants to save
	session.player.ShipData.MechShip = session.player.Mech.MechData
	session.SData = session.player.ShipData

	// handle pause menu
	if rl.IsKeyPressed(rl.KeyP) && !isInteracting { // press P to pause or resume game
		if pGamePaused {
			session.ResumeGame()
		} else {
			pGamePaused = true
		}
	}
	if pGamePaused { // if paused
		if pAreYouSure { // if attempting to return to title or quit without saving
			rl.DrawRectangle(440, 300, 400, 100, rl.White)
			rl.DrawText("You haven't saved, are you sure?", 460, 320, 20, rl.Black)
			rl.DrawText("Y / N", 620, 360, 20, rl.Black)
			if rl.IsKeyPressed(rl.KeyY) {
				if pQuitPressed {
					session.CurrentSystem = -1
				} else {
					session.CurrentSystem = 0
				}
				session.ResumeGame()
				rl.PlaySound(gameData.MenuSounds[0])
			} else if rl.IsKeyPressed(rl.KeyN) {
				pAreYouSure = false
				pQuitPressed = false
				rl.PlaySound(gameData.MenuSounds[0])
			}
		} else { // updates buttons
			pResumeGameButton.UpdateButton(gameData.MenuSounds[0])
			pSaveGameButton.UpdateButton(gameData.MenuSounds[0])
			pReturnToTitleButton.UpdateButton(gameData.MenuSounds[0])
			pQuitGameButton.UpdateButton(gameData.MenuSounds[0])
		}
	}

	CheckForBoss(session)
}

// updates trade menu buttons when interacting with capital
func CapitalInteraction(session *SessionData, gameData *GameData) {
	rMessageButton.DrawButton()
	rFrameUpgradeButton.UpdateButton(gameData.MenuSounds[0])
	rSupportCoreButton.UpdateButton(gameData.MenuSounds[0])
	rExitTradingButton.UpdateButton(gameData.MenuSounds[0])
}

func DrawStats(session *SessionData) {
	rl.DrawText("Money: "+strconv.Itoa(session.player.Money), 0, 60, 20, rl.White)
	rl.DrawText("Level: "+strconv.Itoa(session.player.Level), 0, 80, 20, rl.White)
	rl.DrawText("Frame: "+strconv.Itoa(session.player.ShipData.MechShip.FrameUpgrade), 0, 100, 20, rl.White)
	var numCores int = 0
	for i := 0; i < len(session.player.ShipData.MechShip.SupportCores); i++ {
		if session.player.ShipData.MechShip.SupportCores[i] {
			numCores++
		}
	}
	rl.DrawText("Cores: "+strconv.Itoa(numCores), 0, 120, 20, rl.White)
}

// function for drawing and handling title screen interactions
func TitleScreenLevel(session *SessionData, gameData *GameData) {
	// set music
	session.NewStream = 4

	// draw background
	rl.DrawTexture(gameData.BackgroundImages[3], 0, 0, rl.White)

	// draw buttons
	if tAreYouSure { // if considering starting new game
		rl.DrawRectangle(440, 300, 400, 100, rl.White)
		rl.DrawText("Save already exists, are you sure?", 460, 320, 20, rl.Black)
		rl.DrawText("Y / N", 620, 360, 20, rl.Black)
		if rl.IsKeyPressed(rl.KeyY) {
			session.SaveData()
			tAreYouSure = false
			session.CurrentSystem = 1
			session.player.CurrentSystem = 1
			rl.PlaySound(gameData.MenuSounds[0])
		} else if rl.IsKeyPressed(rl.KeyN) {
			tAreYouSure = false
			rl.PlaySound(gameData.MenuSounds[0])
		}
	} else { // update buttons
		tNewGameButton.UpdateButton(gameData.MenuSounds[0])
		tLoadGameButton.UpdateButton(gameData.MenuSounds[0])
		tQuitButton.UpdateButton(gameData.MenuSounds[0])
		tTitleBox.DrawButton()
	}
}

// united nations specific functionality
func UnitedNationsLevel(session *SessionData, gameData *GameData) {
	GenericLevel(session, UnitedSystem, gameData)
	rl.DrawText("System: United Nations of Earth", 0, 0, 20, rl.White)
	rl.DrawText("System Mech Frame Price: 100", 0, 20, 20, rl.White)
	rl.DrawText("System Support Core Price: 50", 0, 40, 20, rl.White)
	DrawStats(session)

	var numCores int = 0
	for i := 0; i < len(session.player.ShipData.MechShip.SupportCores); i++ {
		if session.player.ShipData.MechShip.SupportCores[i] {
			numCores++
		}
	}

	// handle player interactions
	if CanInteract && rl.IsKeyPressed(rl.KeyE) && !isInteracting && !pGamePaused {
		isInteracting = true
	}

	if isInteracting {
		if CurrentInteractable == 3 { // interacting with capital
			InteractionSystem = 1
			CapitalInteraction(session, gameData)
		} else { // interacting with bases - 1, 2, 4, 5
			Opponent = createBasicEnemy(gameData.PirateKingAnimations, 1)
			session.player.Mech.MechData.Ammo = 100
			session.player.Mech.MechData.Charge = 50
			session.player.Mech.MechData.Armor = 100 * session.player.Mech.FrameUpgrade
			session.player.Mech.MechData.Health = 100 + (10 * numCores)
			if UnitedSystem[CurrentInteractable].isCleared { // determines xp and currency reward
				OldSystem = 1
				session.CurrentSystem = 7
				XPReward = 25
				CurrencyReward = 25
			} else {
				OldSystem = 1
				session.CurrentSystem = 7
				XPReward = 50
				CurrencyReward = 50
			}
			isInteracting = false
		}
	}
}

// grand rev specific functionality
func GrandRevLevel(session *SessionData, gameData *GameData) {
	GenericLevel(session, GrandSystem, gameData)
	rl.DrawText("System: Grand Revolutionaries", 0, 0, 20, rl.White)
	rl.DrawText("System Mech Frame Price: 200", 0, 20, 20, rl.White)
	rl.DrawText("System Support Core Price: 100", 0, 40, 20, rl.White)
	DrawStats(session)

	var numCores int = 0
	for i := 0; i < len(session.player.ShipData.MechShip.SupportCores); i++ {
		if session.player.ShipData.MechShip.SupportCores[i] {
			numCores++
		}
	}

	// handle player interactions
	if CanInteract && rl.IsKeyPressed(rl.KeyE) && !isInteracting && !pGamePaused {
		isInteracting = true
	}

	if isInteracting {
		if CurrentInteractable == 1 { // interacting with capital
			InteractionSystem = 2
			CapitalInteraction(session, gameData)
		} else { // interacting with bases
			Opponent = createBasicEnemy(gameData.JosephAnimations, 2)
			session.player.Mech.MechData.Ammo = 100
			session.player.Mech.MechData.Charge = 50
			session.player.Mech.MechData.Armor = 100 * session.player.Mech.FrameUpgrade
			session.player.Mech.MechData.Health = 100 + (10 * numCores)
			if UnitedSystem[CurrentInteractable].isCleared { // determines xp and currency reward
				OldSystem = 2
				session.CurrentSystem = 7
				XPReward = 50
				CurrencyReward = 50
			} else {
				OldSystem = 2
				session.CurrentSystem = 7
				XPReward = 100
				CurrencyReward = 100
			}
			isInteracting = false
		}
	}
}

// imperium specific functionality
func ImperiumLevel(session *SessionData, gameData *GameData) {
	GenericLevel(session, ImperiumSystem, gameData)
	rl.DrawText("System: The Imperium", 0, 0, 20, rl.White)
	rl.DrawText("System Mech Frame Price: 300", 0, 20, 20, rl.White)
	rl.DrawText("System Support Core Price: 150", 0, 40, 20, rl.White)
	DrawStats(session)

	var numCores int = 0
	for i := 0; i < len(session.player.ShipData.MechShip.SupportCores); i++ {
		if session.player.ShipData.MechShip.SupportCores[i] {
			numCores++
		}
	}

	// handle player interactions
	if CanInteract && rl.IsKeyPressed(rl.KeyE) && !isInteracting && !pGamePaused {
		isInteracting = true
	}

	if isInteracting {
		if CurrentInteractable == 1 { // interacting with capital
			InteractionSystem = 3
			CapitalInteraction(session, gameData)
		} else { // interacting with bases - 1, 2, 4, 5
			Opponent = createBasicEnemy(gameData.ReytAnimations, 3)
			session.player.Mech.MechData.Ammo = 100
			session.player.Mech.MechData.Charge = 50
			session.player.Mech.MechData.Armor = 100 * session.player.Mech.FrameUpgrade
			session.player.Mech.MechData.Health = 100 + (10 * numCores)
			if UnitedSystem[CurrentInteractable].isCleared { // determines xp and currency reward
				OldSystem = 3
				session.CurrentSystem = 7
				XPReward = 75
				CurrencyReward = 75
			} else {
				OldSystem = 3
				session.CurrentSystem = 7
				XPReward = 150
				CurrencyReward = 150
			}
			isInteracting = false
		}
	}
}

// collective specific functionality
func CollectiveLevel(session *SessionData, gameData *GameData) {
	GenericLevel(session, CollectiveSystem, gameData)
	rl.DrawText("System: The Pristine Collective", 0, 0, 20, rl.White)
	rl.DrawText("System Mech Frame Price: 400", 0, 20, 20, rl.White)
	rl.DrawText("System Support Core Price: 200", 0, 40, 20, rl.White)
	DrawStats(session)

	var numCores int = 0
	for i := 0; i < len(session.player.ShipData.MechShip.SupportCores); i++ {
		if session.player.ShipData.MechShip.SupportCores[i] {
			numCores++
		}
	}

	// handle player interactions
	if CanInteract && rl.IsKeyPressed(rl.KeyE) && !isInteracting && !pGamePaused {
		isInteracting = true
	}

	if isInteracting {
		if CurrentInteractable == 1 { // interacting with capital
			InteractionSystem = 4
			CapitalInteraction(session, gameData)
		} else { // interacting with bases - 1, 2, 4, 5
			Opponent = createBasicEnemy(gameData.ElraAnimations, 4)
			session.player.Mech.MechData.Ammo = 100
			session.player.Mech.MechData.Charge = 50
			session.player.Mech.MechData.Armor = 100 * session.player.Mech.FrameUpgrade
			session.player.Mech.MechData.Health = 100 + (10 * numCores)
			if UnitedSystem[CurrentInteractable].isCleared { // determines xp and currency reward
				OldSystem = 4
				session.CurrentSystem = 7
				XPReward = 100
				CurrencyReward = 100
			} else {
				OldSystem = 4
				session.CurrentSystem = 7
				XPReward = 200
				CurrencyReward = 200
			}
			isInteracting = false
		}
	}
}

// golden empire specific functionality
func GoldenLevel(session *SessionData, gameData *GameData) {
	GenericLevel(session, GoldenSystem, gameData)
	rl.DrawText("System: Golden Empire", 0, 0, 20, rl.White)
	rl.DrawText("System Mech Frame Price: 500", 0, 20, 20, rl.White)
	rl.DrawText("System Support Core Price: 250", 0, 40, 20, rl.White)
	DrawStats(session)

	var numCores int = 0
	for i := 0; i < len(session.player.ShipData.MechShip.SupportCores); i++ {
		if session.player.ShipData.MechShip.SupportCores[i] {
			numCores++
		}
	}

	// handle player interactions
	if CanInteract && rl.IsKeyPressed(rl.KeyE) && !isInteracting && !pGamePaused {
		isInteracting = true
	}

	if GoldenSystem[1].isDefeated {
		GoldenSystem[1].isBase = false
	}

	if isInteracting {
		if CurrentInteractable == 1 && !GoldenSystem[1].isBase { // interacting with capital
			InteractionSystem = 5
			CapitalInteraction(session, gameData)
		} else if GoldenSystem[1].isBase { // interacting with bases - 1, 2, 4, 5
			Opponent = createBasicEnemy(gameData.EmperorAnimations, 5)
			session.player.Mech.MechData.Ammo = 100
			session.player.Mech.MechData.Charge = 50
			session.player.Mech.MechData.Armor = 100 * session.player.Mech.FrameUpgrade
			session.player.Mech.MechData.Health = 100 + (10 * numCores)
			if UnitedSystem[CurrentInteractable].isCleared { // determines xp and currency reward
				OldSystem = 5
				session.CurrentSystem = 7
				XPReward = 250
				CurrencyReward = 250
			} else {
				OldSystem = 5
				session.CurrentSystem = 7
				XPReward = 500
				CurrencyReward = 500
			}
			isInteracting = false
		}
	}

	if GoldenSystem[1].isDefeated {
		session.CurrentSystem = 6
	}
}

func CheckForBoss(session *SessionData) {
	cb1 := false
	cb2 := false
	cb3 := false
	cb4 := false
	if session.player.UnitedNationsChecks[0] && session.player.UnitedNationsChecks[1] && session.player.UnitedNationsChecks[2] && session.player.UnitedNationsChecks[3] {
		cb1 = true
	}
	if session.player.GrandRevChecks[0] && session.player.GrandRevChecks[1] && session.player.GrandRevChecks[2] && session.player.GrandRevChecks[3] {
		cb2 = true
	}
	if session.player.ImperiumChecks[0] && session.player.ImperiumChecks[1] && session.player.ImperiumChecks[2] && session.player.ImperiumChecks[3] {
		cb3 = true
	}
	if session.player.CollectiveChecks[0] && session.player.CollectiveChecks[1] && session.player.CollectiveChecks[2] && session.player.CollectiveChecks[3] {
		cb4 = true
	}
	if cb1 && cb2 && cb3 && cb4 {
		GoldenSystem[1].isBase = true
	}
}

func EndScreen(session *SessionData, gameData *GameData) {
	session.NewStream = 4
	rl.DrawTexture(gameData.BackgroundImages[3], 0, 0, rl.White)

	endScreenMessage.DrawButton()

	pHasSaved = true
	endScreenTitle.UpdateButton(gameData.MenuSounds[0])
	pHasSaved = false
}
