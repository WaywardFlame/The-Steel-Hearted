package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameData struct {
	RegularMechAnimations    []rl.Texture2D
	PlayerShipSprites        []rl.Texture2D
	MomentarySprites         []rl.Texture2D
	SolarSystemSprites       []rl.Texture2D
	BackgroundImages         []rl.Texture2D
	ImperiumSystemSprites    []rl.Texture2D
	GrandRevSystemSprites    []rl.Texture2D
	GoldenSystemSprites      []rl.Texture2D
	CollectiveSystemSprites  []rl.Texture2D
	ElraAnimations           []rl.Texture2D
	EmperorAnimations        []rl.Texture2D
	JosephAnimations         []rl.Texture2D
	PirateKingAnimations     []rl.Texture2D
	ReytAnimations           []rl.Texture2D
	BattleSounds             []rl.Sound
	MenuSounds               []rl.Sound
	MusicSounds              []rl.Music
	PauseMenuInitialized     bool
	TradeMenuInitialized     bool
	TitleScreenInitialized   bool
	UnitedNationsInitialized bool
	GrandRevInitialized      bool
	ImperiumInitialized      bool
	CollectiveInitialized    bool
	GoldenEmpireInitialized  bool
	EndScreenInitialized     bool
}

func main() {
	rl.InitWindow(1280, 720, "The Steel Hearted")
	rl.InitAudioDevice()
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	gameData := LoadGameData()
	session := SessionData{
		CurrentSystem: 0,
		SData:         createEmptyPlayerShipData(),
	}
	InitializeLevels(&session, &gameData)
	var defaultChecks [9]bool // default value is false, so no need to set array
	adjustUnitedSystem(defaultChecks)
	adjustGrandSystem(defaultChecks)
	adjustImperiumSystem(defaultChecks)
	adjustCollectiveSystem(defaultChecks)
	adjustGoldenSystem(defaultChecks)
	session.player.Sprites = gameData.PlayerShipSprites
	session.player = copyDataOver(session.SData, gameData.RegularMechAnimations, gameData.PlayerShipSprites)

	for !rl.WindowShouldClose() {
		UpdateCurrentMusic(&session, &gameData)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// don't want to assign data until actually in game
		if PlayerLoaded {
			session.player = copyDataOver(session.SData, gameData.RegularMechAnimations, gameData.PlayerShipSprites)
			adjustUnitedSystem(session.player.UnitedNationsChecks)
			adjustGrandSystem(session.player.GrandRevChecks)
			adjustImperiumSystem(session.player.ImperiumChecks)
			adjustCollectiveSystem(session.player.CollectiveChecks)
			adjustGoldenSystem(session.player.GoldenChecks)
			session.CurrentSystem = session.player.CurrentSystem
			PlayerLoaded = false
		}

		if session.CurrentSystem == -1 {
			rl.EndDrawing()
			break
		} else if session.CurrentSystem == 0 { // title screen level
			TitleScreenLevel(&session, &gameData)
		} else if session.CurrentSystem == 1 { // Solar System
			UnitedNationsLevel(&session, &gameData)
		} else if session.CurrentSystem == 2 { // Grand Revolutionary System
			GrandRevLevel(&session, &gameData)
		} else if session.CurrentSystem == 3 { // Imperium System
			ImperiumLevel(&session, &gameData)
		} else if session.CurrentSystem == 4 { // Collective System
			CollectiveLevel(&session, &gameData)
		} else if session.CurrentSystem == 5 { // Golden Empire System
			GoldenLevel(&session, &gameData)
		} else if session.CurrentSystem == 6 { // end screen, after beating all systems
			EndScreen(&session, &gameData)
		} else if session.CurrentSystem == 7 { // for turn based battles
			Battle(&session.player.Mech, &session, &gameData)
		}

		rl.EndDrawing()
	}
}

func UpdateCurrentMusic(session *SessionData, gameData *GameData) {
	if session.CurrentStream != session.NewStream {
		session.CurrentStream = session.NewStream
		rl.PlayMusicStream(gameData.MusicSounds[session.CurrentStream])
	}
	rl.UpdateMusicStream(gameData.MusicSounds[session.CurrentStream])
}

func InitializeLevels(session *SessionData, gameData *GameData) {
	// create color theme
	colorTheme := NewColorTheme(
		rl.NewColor(255, 255, 255, 255),
		rl.NewColor(128, 255, 255, 255),
		rl.NewColor(0, 0, 0, 255),
	)

	// initialize pause menu
	if !gameData.PauseMenuInitialized {
		// resume game button
		newButton1 := NewButton(0, 0, 300, 50, colorTheme)
		newButton1.SetText("Resume Game", 20)
		newButton1.CenterButtonX()
		newButton1.Y = 160
		newButton1.AddOnClickFunc(session.ResumeGame)

		// save game button
		newButton2 := NewButton(0, 0, 300, 50, colorTheme)
		newButton2.SetText("Save Game", 20)
		newButton2.CenterButtonX()
		newButton2.Y = 260
		newButton2.AddOnClickFunc(session.SaveData)

		// title screen button
		newButton3 := NewButton(0, 0, 300, 50, colorTheme)
		newButton3.SetText("Return To Title", 20)
		newButton3.CenterButtonX()
		newButton3.Y = 360
		newButton3.AddOnClickFunc(session.ReturnToTitle)

		// quit game game
		newButton4 := NewButton(0, 0, 300, 50, colorTheme)
		newButton4.SetText("Quit Game", 20)
		newButton4.CenterButtonX()
		newButton4.Y = 460
		newButton4.AddOnClickFunc(session.QuitInGame)

		// assign buttons
		pResumeGameButton = newButton1
		pSaveGameButton = newButton2
		pReturnToTitleButton = newButton3
		pQuitGameButton = newButton4

		gameData.PauseMenuInitialized = true
	}

	// initialize trade menu
	if !gameData.TradeMenuInitialized {
		// button displaying message
		newButton1 := NewButton(0, 0, 500, 100, colorTheme)
		newButton1.SetText("Purchase whatever you like!", 20)
		newButton1.CenterButtonX()
		newButton1.Y = 100

		// purchase frame upgrade button
		newButton2 := NewButton(0, 0, 300, 50, colorTheme)
		newButton2.SetText("Frame Upgrade", 20)
		newButton2.CenterButtonX()
		newButton2.Y = 260
		newButton2.AddOnClickFunc(session.UpgradeTheFrame)

		// purchase support core button
		newButton3 := NewButton(0, 0, 300, 50, colorTheme)
		newButton3.SetText("New Support Core", 20)
		newButton3.CenterButtonX()
		newButton3.Y = 360
		newButton3.AddOnClickFunc(session.PurchaseCore)

		// exit trade button
		newButton4 := NewButton(0, 0, 300, 50, colorTheme)
		newButton4.SetText("Exit Trading", 20)
		newButton4.CenterButtonX()
		newButton4.Y = 460
		newButton4.AddOnClickFunc(session.QuitInteraction)

		// assign buttons
		rMessageButton = newButton1
		rFrameUpgradeButton = newButton2
		rSupportCoreButton = newButton3
		rExitTradingButton = newButton4
	}

	// initialized title screen
	if !gameData.TitleScreenInitialized {
		// new game button
		newButton1 := NewButton(0, 0, 300, 50, colorTheme)
		newButton1.SetText("New Game", 20)
		newButton1.CenterButtonX()
		newButton1.Y = 360
		newButton1.AddOnClickFunc(session.CreateSave)

		// load game button
		newButton2 := NewButton(0, 0, 300, 50, colorTheme)
		newButton2.SetText("Load Game", 20)
		newButton2.CenterButtonX()
		newButton2.Y = 460
		newButton2.AddOnClickFunc(session.LoadData)

		// quit game button
		newButton3 := NewButton(0, 0, 300, 50, colorTheme)
		newButton3.SetText("Quit Game", 20)
		newButton3.CenterButtonX()
		newButton3.Y = 560
		newButton3.AddOnClickFunc(session.QuitGame)

		// title box, button for ease of creation, won't be updated
		newButton4 := NewButton(0, 0, 500, 200, colorTheme)
		newButton4.SetText("The Steel Hearted\nby WaywardFlame", 40)
		newButton4.CenterButtonX()
		newButton4.Y = 100

		// assign to globals in level.go
		tNewGameButton = newButton1
		tLoadGameButton = newButton2
		tQuitButton = newButton3
		tTitleBox = newButton4

		// set initialized to true
		gameData.TitleScreenInitialized = true
	}

	// initialize the united nations of earth level
	if !gameData.UnitedNationsInitialized {
		UnitedSystem = append(UnitedSystem, Body{position: StarPosition, radius: 128, texture: gameData.SolarSystemSprites[6]}) // sun
		UnitedSystem = append(UnitedSystem, Body{position: Planet1, radius: 64, texture: gameData.SolarSystemSprites[3]})       // mercury - base 1
		UnitedSystem = append(UnitedSystem, Body{position: Planet2, radius: 64, texture: gameData.SolarSystemSprites[8]})       // venus - base 2
		UnitedSystem = append(UnitedSystem, Body{position: Planet3, radius: 64, texture: gameData.SolarSystemSprites[0]})       // earth - capital
		UnitedSystem = append(UnitedSystem, Body{position: Planet4, radius: 64, texture: gameData.SolarSystemSprites[2]})       // mars - base 3
		UnitedSystem = append(UnitedSystem, Body{position: Planet5, radius: 64, texture: gameData.SolarSystemSprites[1]})       // jupiter - base 4
		UnitedSystem = append(UnitedSystem, Body{position: Planet6, radius: 64, texture: gameData.SolarSystemSprites[5]})       // saturn
		UnitedSystem = append(UnitedSystem, Body{position: Planet7, radius: 64, texture: gameData.SolarSystemSprites[7]})       // uranus
		UnitedSystem = append(UnitedSystem, Body{position: Planet8, radius: 64, texture: gameData.SolarSystemSprites[4]})       // neptune

		gameData.UnitedNationsInitialized = true
	}

	// initialize the grand revoluationaries level
	if !gameData.GrandRevInitialized {
		GrandSystem = append(GrandSystem, Body{position: StarPosition, radius: 128, texture: gameData.GrandRevSystemSprites[3]}) // grand star
		GrandSystem = append(GrandSystem, Body{position: Planet3, radius: 64, texture: gameData.GrandRevSystemSprites[4]})       // grand capital
		GrandSystem = append(GrandSystem, Body{position: Planet1, radius: 64, texture: gameData.GrandRevSystemSprites[5]})       // grand base 1
		GrandSystem = append(GrandSystem, Body{position: Planet2, radius: 64, texture: gameData.GrandRevSystemSprites[6]})       // grand base 2
		GrandSystem = append(GrandSystem, Body{position: Planet4, radius: 64, texture: gameData.GrandRevSystemSprites[7]})       // grand base 3
		GrandSystem = append(GrandSystem, Body{position: Planet7, radius: 64, texture: gameData.GrandRevSystemSprites[8]})       // grand base 4
		GrandSystem = append(GrandSystem, Body{position: Planet5, radius: 64, texture: gameData.GrandRevSystemSprites[0]})       // grand planet 1
		GrandSystem = append(GrandSystem, Body{position: Planet6, radius: 64, texture: gameData.GrandRevSystemSprites[1]})       // grand planet 2
		GrandSystem = append(GrandSystem, Body{position: Planet8, radius: 64, texture: gameData.GrandRevSystemSprites[2]})       // grand planet 3

		gameData.GrandRevInitialized = true
	}

	// initialize the imperium level
	if !gameData.ImperiumInitialized {
		ImperiumSystem = append(ImperiumSystem, Body{position: StarPosition, radius: 128, texture: gameData.ImperiumSystemSprites[8]}) // imperium star
		ImperiumSystem = append(ImperiumSystem, Body{position: Planet5, radius: 64, texture: gameData.ImperiumSystemSprites[7]})       // imperium capital
		ImperiumSystem = append(ImperiumSystem, Body{position: Planet1, radius: 64, texture: gameData.ImperiumSystemSprites[0]})       // imperium base 1
		ImperiumSystem = append(ImperiumSystem, Body{position: Planet2, radius: 64, texture: gameData.ImperiumSystemSprites[1]})       // imperium base 2
		ImperiumSystem = append(ImperiumSystem, Body{position: Planet7, radius: 64, texture: gameData.ImperiumSystemSprites[2]})       // imperium base 3
		ImperiumSystem = append(ImperiumSystem, Body{position: Planet8, radius: 64, texture: gameData.ImperiumSystemSprites[3]})       // imperium base 4
		ImperiumSystem = append(ImperiumSystem, Body{position: Planet3, radius: 64, texture: gameData.ImperiumSystemSprites[4]})       // imperium planet 1
		ImperiumSystem = append(ImperiumSystem, Body{position: Planet4, radius: 64, texture: gameData.ImperiumSystemSprites[5]})       // imperium planet 2
		ImperiumSystem = append(ImperiumSystem, Body{position: Planet6, radius: 64, texture: gameData.ImperiumSystemSprites[6]})       // imperium planet 3

		gameData.ImperiumInitialized = true
	}

	// initialize the collective level
	if !gameData.CollectiveInitialized {
		CollectiveSystem = append(CollectiveSystem, Body{position: StarPosition, radius: 128, texture: gameData.CollectiveSystemSprites[8]}) // collective star
		CollectiveSystem = append(CollectiveSystem, Body{position: Planet8, radius: 64, texture: gameData.CollectiveSystemSprites[7]})       // collective capital
		CollectiveSystem = append(CollectiveSystem, Body{position: Planet1, radius: 64, texture: gameData.CollectiveSystemSprites[3]})       // collective base 1
		CollectiveSystem = append(CollectiveSystem, Body{position: Planet2, radius: 64, texture: gameData.CollectiveSystemSprites[4]})       // collective base 2
		CollectiveSystem = append(CollectiveSystem, Body{position: Planet7, radius: 64, texture: gameData.CollectiveSystemSprites[5]})       // collective base 3
		CollectiveSystem = append(CollectiveSystem, Body{position: Planet5, radius: 64, texture: gameData.CollectiveSystemSprites[6]})       // collective base 4
		CollectiveSystem = append(CollectiveSystem, Body{position: Planet3, radius: 64, texture: gameData.CollectiveSystemSprites[0]})       // collective planet 1
		CollectiveSystem = append(CollectiveSystem, Body{position: Planet4, radius: 64, texture: gameData.CollectiveSystemSprites[1]})       // collective planet 2
		CollectiveSystem = append(CollectiveSystem, Body{position: Planet6, radius: 64, texture: gameData.CollectiveSystemSprites[2]})       // collective planet 3

		gameData.CollectiveInitialized = true
	}

	if !gameData.GoldenEmpireInitialized {
		GoldenSystem = append(GoldenSystem, Body{position: StarPosition, radius: 128, texture: gameData.GoldenSystemSprites[8]}) // golden star
		GoldenSystem = append(GoldenSystem, Body{position: Planet1, radius: 64, texture: gameData.GoldenSystemSprites[7]})       // golden capital
		GoldenSystem = append(GoldenSystem, Body{position: Planet2, radius: 64, texture: gameData.GoldenSystemSprites[0]})       // golden planet 1
		GoldenSystem = append(GoldenSystem, Body{position: Planet3, radius: 64, texture: gameData.GoldenSystemSprites[1]})       // golden planet 2
		GoldenSystem = append(GoldenSystem, Body{position: Planet4, radius: 64, texture: gameData.GoldenSystemSprites[2]})       // golden planet 3
		GoldenSystem = append(GoldenSystem, Body{position: Planet5, radius: 64, texture: gameData.GoldenSystemSprites[3]})       // golden planet 4
		GoldenSystem = append(GoldenSystem, Body{position: Planet6, radius: 64, texture: gameData.GoldenSystemSprites[4]})       // golden planet 5
		GoldenSystem = append(GoldenSystem, Body{position: Planet7, radius: 64, texture: gameData.GoldenSystemSprites[5]})       // golden planet 6
		GoldenSystem = append(GoldenSystem, Body{position: Planet8, radius: 64, texture: gameData.GoldenSystemSprites[6]})       // golden planet 7

		gameData.GoldenEmpireInitialized = true
	}

	if !gameData.EndScreenInitialized {
		newButton1 := NewButton(0, 0, 400, 200, colorTheme)
		newButton1.SetText("You beat the game. Congrats!", 20)
		newButton1.CenterButtonX()
		newButton1.CenterButtonY()

		newButton2 := NewButton(0, 0, 200, 100, colorTheme)
		newButton2.SetText("Return to Title", 20)
		newButton2.CenterButtonX()
		newButton2.CenterButtonY()
		newButton2.Y += 200
		newButton2.AddOnClickFunc(session.ReturnToTitle)

		endScreenTitle = newButton2
		endScreenMessage = newButton1
	}
}

func LoadGameData() GameData {
	// load in textures
	regularMechAnimations := []rl.Texture2D{
		rl.LoadTexture("textures/regular_mech_animations/mech_idle.png"),
		rl.LoadTexture("textures/regular_mech_animations/mech_guard.png"),
		rl.LoadTexture("textures/regular_mech_animations/mech_melee.png"),
		rl.LoadTexture("textures/regular_mech_animations/mech_parry.png"),
		rl.LoadTexture("textures/regular_mech_animations/mech_laser.png"),
		rl.LoadTexture("textures/regular_mech_animations/mech_missile.png"),
		rl.LoadTexture("textures/regular_mech_animations/mech_core_switch.png"),
	}
	playerShipSprites := []rl.Texture2D{
		rl.LoadTexture("textures/player_ship_animations/ship_east.png"),
		rl.LoadTexture("textures/player_ship_animations/ship_north.png"),
		rl.LoadTexture("textures/player_ship_animations/ship_south.png"),
		rl.LoadTexture("textures/player_ship_animations/ship_west.png"),
	}
	momentarySprites := []rl.Texture2D{
		rl.LoadTexture("textures/momentary_sprites/laser_bullet.png"),
		rl.LoadTexture("textures/momentary_sprites/missile_bullet.png"),
		rl.LoadTexture("textures/momentary_sprites/parry_alert.png"),
	}
	solarSystemSprites := []rl.Texture2D{
		rl.LoadTexture("textures/solar_system/earth.png"),
		rl.LoadTexture("textures/solar_system/jupiter.png"),
		rl.LoadTexture("textures/solar_system/mars.png"),
		rl.LoadTexture("textures/solar_system/mercury.png"),
		rl.LoadTexture("textures/solar_system/neptune.png"),
		rl.LoadTexture("textures/solar_system/saturn.png"),
		rl.LoadTexture("textures/solar_system/sun.png"),
		rl.LoadTexture("textures/solar_system/uranus.png"),
		rl.LoadTexture("textures/solar_system/venus.png"),
	}
	backgroundImages := []rl.Texture2D{
		rl.LoadTexture("textures/boss_battle_background.png"),
		rl.LoadTexture("textures/emperor_battle_background.png"),
		rl.LoadTexture("textures/regular_battle_background.png"),
		rl.LoadTexture("textures/title_screen_background.png"),
	}
	imperiumSystemSprites := []rl.Texture2D{
		rl.LoadTexture("textures/imperium_system/bp_1.png"),
		rl.LoadTexture("textures/imperium_system/bp_2.png"),
		rl.LoadTexture("textures/imperium_system/bp_3.png"),
		rl.LoadTexture("textures/imperium_system/bp_4.png"),
		rl.LoadTexture("textures/imperium_system/gp_1.png"),
		rl.LoadTexture("textures/imperium_system/gp_2.png"),
		rl.LoadTexture("textures/imperium_system/gp_3.png"),
		rl.LoadTexture("textures/imperium_system/imperium_capital.png"),
		rl.LoadTexture("textures/imperium_system/purple_star.png"),
	}
	grandRevSystemSprites := []rl.Texture2D{
		rl.LoadTexture("textures/grand_rev_system/pp_1.png"),
		rl.LoadTexture("textures/grand_rev_system/pp_2.png"),
		rl.LoadTexture("textures/grand_rev_system/pp_3.png"),
		rl.LoadTexture("textures/grand_rev_system/red_star.png"),
		rl.LoadTexture("textures/grand_rev_system/rev_capital.png"),
		rl.LoadTexture("textures/grand_rev_system/rp_1.png"),
		rl.LoadTexture("textures/grand_rev_system/rp_2.png"),
		rl.LoadTexture("textures/grand_rev_system/rp_3.png"),
		rl.LoadTexture("textures/grand_rev_system/rp_4.png"),
	}
	goldenSystemSprites := []rl.Texture2D{
		rl.LoadTexture("textures/golden_system/ggp_1.png"),
		rl.LoadTexture("textures/golden_system/ggp_2.png"),
		rl.LoadTexture("textures/golden_system/ggp_3.png"),
		rl.LoadTexture("textures/golden_system/ggp_4.png"),
		rl.LoadTexture("textures/golden_system/ggp_5.png"),
		rl.LoadTexture("textures/golden_system/ggp_6.png"),
		rl.LoadTexture("textures/golden_system/ggp_7.png"),
		rl.LoadTexture("textures/golden_system/golden_capital.png"),
		rl.LoadTexture("textures/golden_system/orange_star.png"),
	}
	collectiveSystemSprites := []rl.Texture2D{
		rl.LoadTexture("textures/collective_system/cbp_1.png"),
		rl.LoadTexture("textures/collective_system/cbp_2.png"),
		rl.LoadTexture("textures/collective_system/cbp_3.png"),
		rl.LoadTexture("textures/collective_system/cgp_1.png"),
		rl.LoadTexture("textures/collective_system/cgp_2.png"),
		rl.LoadTexture("textures/collective_system/cgp_3.png"),
		rl.LoadTexture("textures/collective_system/cgp_4.png"),
		rl.LoadTexture("textures/collective_system/collective_capital.png"),
		rl.LoadTexture("textures/collective_system/white_star.png"),
	}
	elraAnimations := []rl.Texture2D{
		rl.LoadTexture("textures/boss_mech_animations/elra/elra_idle.png"),
		rl.LoadTexture("textures/boss_mech_animations/elra/elra_guard.png"),
		rl.LoadTexture("textures/boss_mech_animations/elra/elra_laser.png"),
		rl.LoadTexture("textures/boss_mech_animations/elra/elra_melee.png"),
		rl.LoadTexture("textures/boss_mech_animations/elra/elra_missile.png"),
		rl.LoadTexture("textures/boss_mech_animations/elra/elra_parry.png"),
	}
	emperorAnimations := []rl.Texture2D{
		rl.LoadTexture("textures/boss_mech_animations/emperor/emperor_idle.png"),
		rl.LoadTexture("textures/boss_mech_animations/emperor/emperor_guard.png"),
		rl.LoadTexture("textures/boss_mech_animations/emperor/emperor_laser.png"),
		rl.LoadTexture("textures/boss_mech_animations/emperor/emperor_melee.png"),
		rl.LoadTexture("textures/boss_mech_animations/emperor/emperor_missile.png"),
		rl.LoadTexture("textures/boss_mech_animations/emperor/emperor_parry.png"),
	}
	josephAnimations := []rl.Texture2D{
		rl.LoadTexture("textures/boss_mech_animations/joseph/joseph_idle.png"),
		rl.LoadTexture("textures/boss_mech_animations/joseph/joseph_guard.png"),
		rl.LoadTexture("textures/boss_mech_animations/joseph/joseph_laser.png"),
		rl.LoadTexture("textures/boss_mech_animations/joseph/joseph_melee.png"),
		rl.LoadTexture("textures/boss_mech_animations/joseph/joseph_missile.png"),
		rl.LoadTexture("textures/boss_mech_animations/joseph/joseph_parry.png"),
	}
	pirateKingAnimations := []rl.Texture2D{
		rl.LoadTexture("textures/boss_mech_animations/pirate_king/pirate_idle.png"),
		rl.LoadTexture("textures/boss_mech_animations/pirate_king/pirate_guard.png"),
		rl.LoadTexture("textures/boss_mech_animations/pirate_king/pirate_laser.png"),
		rl.LoadTexture("textures/boss_mech_animations/pirate_king/pirate_melee.png"),
		rl.LoadTexture("textures/boss_mech_animations/pirate_king/pirate_missile.png"),
		rl.LoadTexture("textures/boss_mech_animations/pirate_king/pirate_parry.png"),
	}
	reytAnimations := []rl.Texture2D{
		rl.LoadTexture("textures/boss_mech_animations/reyt/reyt_idle.png"),
		rl.LoadTexture("textures/boss_mech_animations/reyt/reyt_guard.png"),
		rl.LoadTexture("textures/boss_mech_animations/reyt/reyt_laser.png"),
		rl.LoadTexture("textures/boss_mech_animations/reyt/reyt_melee.png"),
		rl.LoadTexture("textures/boss_mech_animations/reyt/reyt_missile.png"),
		rl.LoadTexture("textures/boss_mech_animations/reyt/reyt_parry.png"),
	}

	// load in audio
	battleSounds := []rl.Sound{
		rl.LoadSound("audio/battle_sounds/battle_transition_sound.wav"),
		rl.LoadSound("audio/battle_sounds/boss_defeat_sound.wav"),
		rl.LoadSound("audio/battle_sounds/core_switch_sound.wav"),
		rl.LoadSound("audio/battle_sounds/enemy_defeat_sound.wav"),
		rl.LoadSound("audio/battle_sounds/guard_sound.wav"),
		rl.LoadSound("audio/battle_sounds/laser_sound.wav"),
		rl.LoadSound("audio/battle_sounds/melee_sound.wav"),
		rl.LoadSound("audio/battle_sounds/missile_explosion_sound.wav"),
		rl.LoadSound("audio/battle_sounds/missile_launch_sound.wav"),
		rl.LoadSound("audio/battle_sounds/parry_sound.wav"),
	}
	menuSounds := []rl.Sound{
		rl.LoadSound("audio/menu_sounds/title_and_pause_menu_sound.wav"),
		rl.LoadSound("audio/menu_sounds/trade_menu_sound.wav"),
	}
	musicSounds := []rl.Music{
		rl.LoadMusicStream("audio/music/boss_battle_music.mp3"),
		rl.LoadMusicStream("audio/music/exploration_background_music.mp3"),
		rl.LoadMusicStream("audio/music/golden_empire_background_music.mp3"),
		rl.LoadMusicStream("audio/music/regular_battle_music.wav"),
		rl.LoadMusicStream("audio/music/title_screen_music.wav"),
	}

	// create game data struct to store data
	gameData := GameData{
		RegularMechAnimations:    regularMechAnimations,
		PlayerShipSprites:        playerShipSprites,
		MomentarySprites:         momentarySprites,
		SolarSystemSprites:       solarSystemSprites,
		BackgroundImages:         backgroundImages,
		ImperiumSystemSprites:    imperiumSystemSprites,
		GrandRevSystemSprites:    grandRevSystemSprites,
		GoldenSystemSprites:      goldenSystemSprites,
		CollectiveSystemSprites:  collectiveSystemSprites,
		ElraAnimations:           elraAnimations,
		EmperorAnimations:        emperorAnimations,
		JosephAnimations:         josephAnimations,
		PirateKingAnimations:     pirateKingAnimations,
		ReytAnimations:           reytAnimations,
		BattleSounds:             battleSounds,
		MenuSounds:               menuSounds,
		MusicSounds:              musicSounds,
		PauseMenuInitialized:     false,
		TradeMenuInitialized:     false,
		TitleScreenInitialized:   false,
		UnitedNationsInitialized: false,
		GrandRevInitialized:      false,
		ImperiumInitialized:      false,
		CollectiveInitialized:    false,
		GoldenEmpireInitialized:  false,
		EndScreenInitialized:     false,
	}

	// return gameData
	return gameData
}
