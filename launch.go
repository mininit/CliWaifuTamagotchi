package main

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"

	"CliWaifuTamagotchi/utils"
)

// ==============================
// ASSETS
// ==============================
type Assets struct {
	head           string
	headBlink      string
	happyHead      string
	body           string
	encouragements []string
}

// Load all ASCII and text assets
func loadAssets() (*Assets, error) {
	load := func(path string) string {
		data, err := utils.LoadASCII(path)
		if err != nil {
			panic(fmt.Sprintf("Failed to load %s: %v", path, err))
		}
		return data
	}

	encouragements, err := utils.LoadEncouragements("assets/words-of-encouragement.txt")
	if err != nil {
		return nil, fmt.Errorf("could not load encouragements: %v", err)
	}

	return &Assets{
		head:           load("ascii-arts/expressions/neutral"),
		headBlink:      load("ascii-arts/expressions/neutral-blink"),
		happyHead:      load("ascii-arts/expressions/-happy"),
		body:           load("ascii-arts/clothes/seifuku"),
		encouragements: encouragements,
	}, nil
}

// ==============================
// UI CREATION
// ==============================
type UI struct {
	app         *tview.Application
	actionSpace *tview.List
	waifuArt    *tview.TextView
	chatBox     *tview.TextView
	grid        *tview.Grid
	stopBlink   chan bool
}

func createUI(assets *Assets) *UI {
	app := tview.NewApplication()

	actionSpace := tview.NewList()
	actionSpace.SetTitle("| Action Space |").SetBorder(true)
	
	waifuArt := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	waifuArt.SetText(assets.head + "\n" + assets.body)
	waifuArt.SetBorder(true).SetTitle("| Waifu |")

	chatBox := tview.NewTextView().
		SetTextAlign(tview.AlignCenter)
	chatBox.SetText("...")
	chatBox.SetBorder(true).SetTitle("| Chatbox |")

	grid := tview.NewGrid().
		SetRows(0, 3).
		SetColumns(40, 0).
		SetBorders(false).
		AddItem(actionSpace, 0, 0, 2, 1, 0, 0, true).
		AddItem(waifuArt, 0, 1, 1, 1, 0, 75, false).
		AddItem(chatBox, 1, 1, 1, 1, 0, 0, false)

	return &UI{app, actionSpace, waifuArt, chatBox, grid, make(chan bool)}
}

// ==============================
// MENU SETUP
// ==============================
func setupMenu(ui *UI, assets *Assets, encourageLocked *bool, currentBody *string) {
	ui.actionSpace.AddItem("Encourage", "Get a nice message.", '1', func() {
		if !*encourageLocked {
			*encourageLocked = true
			utils.Encourage(ui.app, ui.waifuArt, ui.chatBox,
				assets.head, assets.happyHead, *currentBody,
				assets.encouragements, 1*time.Second,
				func() { *encourageLocked = false })
		}
	})

	ui.actionSpace.AddItem("Dress Up", "Change her outfit.", '2', func() {

		utils.DressUp(ui.app, ui.waifuArt, ui.chatBox,
			assets.head, assets.headBlink,
			ui.grid, ui.actionSpace, currentBody)
	})

	ui.actionSpace.AddItem("Quit", "Exit the application.", 'q', func() {
		ui.app.Stop()
	})
}

// ==============================
// GLOBAL KEYS
// ==============================
func setGlobalKeys(ui *UI, assets *Assets, encourageLocked *bool, currentBody *string) {
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case '1':
			if !*encourageLocked {
				*encourageLocked = true
				utils.Encourage(ui.app, ui.waifuArt, ui.chatBox,
					assets.head, assets.happyHead, *currentBody,
					assets.encouragements, 1*time.Second,
					func() { *encourageLocked = false })
			}
			return nil
		case '2':
			utils.DressUp(ui.app, ui.waifuArt, ui.chatBox,
				assets.head, assets.headBlink,
				ui.grid, ui.actionSpace, currentBody)
			return nil
		case 'q':
			ui.app.Stop()
			return nil
		}
		return event
	})
}

// ==============================
// MAIN
// ==============================
func main() {
	// ===== Load assets
	// =====
	assets, err := loadAssets()
	if err != nil {
		panic(err)
	}

	// ===== Set UI up
	// =====
	ui := createUI(assets)
	// Get the palette
	palette, err := utils.LoadPalette()
	if err != nil {
		panic(fmt.Sprintf("Failed to load palette: %v", err))
	}
	// Apply to TextViews
	utils.ApplyTextViewPalette(palette, ui.waifuArt, ui.chatBox)
	// Apply to Lists
	utils.ApplyListPalette(palette, ui.actionSpace)

	// ===== Variable work
	// =====
	var encourageLocked bool
	currentBody := assets.body

	// ===== Set functions
	// =====
	setupMenu(ui, assets, &encourageLocked, &currentBody)
	setGlobalKeys(ui, assets, &encourageLocked, &currentBody)

	// ===== Auto processes
	// =====
	ui.stopBlink = utils.StartBlinking(ui.app, ui.waifuArt,
		&assets.head, &assets.headBlink, &currentBody, 5*time.Second)

	// ===== No returns - Error handling
	// =====
	if err := utils.LoadClothes("ascii-arts/clothes"); err != nil {
		panic(err)
	}
	if err := ui.app.SetRoot(ui.grid, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
	if err := utils.CreatePaletteFile(); err != nil {
		panic(err)
	}
}