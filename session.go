package main

import (
	"fmt"
	"os"
)

type SessionData struct {
	NewStream     int
	CurrentStream int
	CurrentSystem int
	SData         ShipData
	player        PlayerShip
}

// resumes the game
func (s *SessionData) ResumeGame() {
	pAreYouSure = false
	pGamePaused = false
	pHasSaved = false
}

// returns to title screen, asks if you are sure if game hasn't been saved
func (s *SessionData) ReturnToTitle() {
	if !pHasSaved {
		pAreYouSure = true
	} else {
		s.CurrentSystem = 0
		s.ResumeGame()
	}
}

// wrapper (?) function to save player data, handles potential error
func (s *SessionData) SaveData() {
	if err := s.SData.SavePlayer(); err != nil {
		fmt.Println("Could not save data. Error occured.")
	} else if pGamePaused {
		pHasSaved = true
	}
}

// wrapper (?) function to load player data, handles potential error
func (s *SessionData) LoadData() {
	if err := s.SData.LoadPlayer(); err != nil {
		fmt.Println("Could not load data. Data likely does not exist.")
	} else {
		PlayerLoaded = true
	}
}

// function for creating a new save game, checks if save already exists
func (s *SessionData) CreateSave() {
	if _, err := os.Stat("data/steelheartedsave.json"); err != nil { // if does not exist
		s.SaveData()
		s.CurrentSystem = 1
		s.player.CurrentSystem = 1
	} else { // if exists
		tAreYouSure = true
	}
}

// function for quitting the game, does not save
func (s *SessionData) QuitGame() {
	s.CurrentSystem = -1
}

// quit from in game
func (s *SessionData) QuitInGame() {
	if !pHasSaved {
		pAreYouSure = true
		pQuitPressed = true
	} else {
		s.QuitGame()
	}
}

// upgrade the mech's frame, trade menu
func (s *SessionData) UpgradeTheFrame() {
	if InteractionSystem == 1 {
		if s.player.Mech.FrameUpgrade < 2 && s.player.Money >= 100 {
			s.player.Mech.FrameUpgrade = 2
			s.player.Money -= 100
			s.QuitInteraction()
		}
	} else if InteractionSystem == 2 {
		if s.player.Mech.FrameUpgrade < 3 && s.player.Money >= 200 {
			s.player.Mech.FrameUpgrade = 3
			s.player.Money -= 200
			s.QuitInteraction()
		}
	} else if InteractionSystem == 3 {
		if s.player.Mech.FrameUpgrade < 4 && s.player.Money >= 300 {
			s.player.Mech.FrameUpgrade = 4
			s.player.Money -= 300
			s.QuitInteraction()
		}
	} else if InteractionSystem == 4 {
		if s.player.Mech.FrameUpgrade < 5 && s.player.Money >= 400 {
			s.player.Mech.FrameUpgrade = 5
			s.player.Money -= 400
			s.QuitInteraction()
		}
	} else if InteractionSystem == 5 {
		if s.player.Mech.FrameUpgrade < 6 && s.player.Money >= 500 {
			s.player.Mech.FrameUpgrade = 6
			s.player.Money -= 500
			s.QuitInteraction()
		}
	}
}

// purchase the support core, trade menu
func (s *SessionData) PurchaseCore() {
	if InteractionSystem == 1 {
		if !s.player.Mech.SupportCores[3] && s.player.Money >= 50 {
			s.player.Mech.SupportCores[3] = true
			s.player.Money -= 50
			s.QuitInteraction()
		}
	} else if InteractionSystem == 2 {
		if !s.player.Mech.SupportCores[4] && s.player.Money >= 100 {
			s.player.Mech.SupportCores[4] = true
			s.player.Money -= 100
			s.QuitInteraction()
		}
	} else if InteractionSystem == 3 {
		if !s.player.Mech.SupportCores[5] && s.player.Money >= 150 {
			s.player.Mech.SupportCores[5] = true
			s.player.Money -= 150
			s.QuitInteraction()
		}
	} else if InteractionSystem == 4 {
		if !s.player.Mech.SupportCores[6] && s.player.Money >= 200 {
			s.player.Mech.SupportCores[6] = true
			s.player.Money -= 200
			s.QuitInteraction()
		}
	} else if InteractionSystem == 5 {
		if !s.player.Mech.SupportCores[7] && s.player.Money >= 250 {
			s.player.Mech.SupportCores[7] = true
			s.player.Money -= 250
			s.QuitInteraction()
		}
	}
}

// quit the trading interaction, trade menu
func (s *SessionData) QuitInteraction() {
	if isInteracting {
		isInteracting = false
	}
}
