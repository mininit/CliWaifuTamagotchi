package main

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"

	"cliwt/utils"
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
	encouragements, err := utils.LoadEncouragements("assets/words-of-encouragement.txt")
	if err != nil {
		return nil, fmt.Errorf("could not load encouragements: %v", err)
	}

	return &Assets{
		head:           utils.LoadASCII(utils.BasePath + "/expressions/neutral"),
		headBlink:      utils.LoadASCII(utils.BasePath + "/expressions/neutral-blink"),
		happyHead:      utils.LoadASCII(utils.BasePath + "/expressions/-happy"),
		body:           utils.LoadASCII(utils.BasePath + "/clothes/hoodie"),
		encouragements: encouragements,
	}, nil
}

// ==============================
// UI CREATION
// ==============================
type UI struct {
	app          *tview.Application
	actionSpace  *tview.List
	happinessBar *tview.TextView
	waifuArt     *tview.TextView
	chatBox      *tview.TextView
	grid         *tview.Grid
	stopBlink   chan bool
}

func createUI(assets *Assets) *UI {
	app := tview.NewApplication()

	actionSpace := tview.NewList()
	actionSpace.SetBorder(true).SetTitle("| Action Space |")

	happinessBar := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	happinessBar.SetBorder(true).SetTitle("| Happiness Bar |")

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
		AddItem(actionSpace,  0, 0, 1, 1, 0, 0,  true).
		AddItem(happinessBar, 1, 0, 1, 1, 0, 0,  false).
		AddItem(waifuArt,     0, 1, 1, 1, 0, 75, false).
		AddItem(chatBox,      1, 1, 1, 1, 0, 0,  false)

	return &UI{app, actionSpace, happinessBar, waifuArt, chatBox, grid, make(chan bool)}
}

// ==============================
// ACTION SPACE SETUP
// ==============================
func setupActionSpace(ui *UI, assets *Assets, encourageLocked *bool, currentBody *string, keys utils.KeyBindings, waifuName string) {
	ui.actionSpace.AddItem("Encourage", "  Get a nice message.", rune(keys.Encourage[0]), func() {
		if !*encourageLocked {
			*encourageLocked = true
			utils.Encourage(ui.app, ui.waifuArt, ui.chatBox,
				assets.head, assets.happyHead, *currentBody, waifuName,
				assets.encouragements, 1*time.Second,
				func() { *encourageLocked = false })
		}
	})

	ui.actionSpace.AddItem("Gift", "  Give a gift.", rune(keys.Gift[0]), func() {
		if !utils.LockGridChanges {
			utils.GiftMenu(ui.app, ui.grid, ui.actionSpace, ui.waifuArt, ui.chatBox,
				assets.head, assets.happyHead, waifuName, currentBody)
		}
	})

	ui.actionSpace.AddItem("Dress Up", "  Change the outfit.", rune(keys.DressUp[0]), func() {
		if !utils.LockGridChanges {
			utils.DressUp(ui.app, ui.grid, ui.actionSpace,ui.waifuArt, ui.chatBox,
				assets.head, waifuName, currentBody)
		}
	})

	ui.actionSpace.AddItem("Background Mode", "  Remove all odd TUI.", rune(keys.BackgroundMode[0]), func() {
		utils.BackgroundMode(ui.app, ui.grid, ui.waifuArt, ui.chatBox, ui.happinessBar, ui.actionSpace, currentBody)
	})

	ui.actionSpace.AddItem("Quit", "  Exit the application.", rune(keys.Quit[0]), func() {
		ui.app.Stop()
	})
}

// ==============================
// GLOBAL KEYS
// ==============================
func setGlobalKeys(ui *UI, assets *Assets, encourageLocked *bool, currentBody *string, keys utils.KeyBindings, isVimNavigation bool, waifuName string) {
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		
		// Main keys for actions
		switch event.Rune() {
		case rune(keys.Encourage[0]):
			if !*encourageLocked {
				*encourageLocked = true
				utils.Encourage(ui.app, ui.waifuArt, ui.chatBox,
					assets.head, assets.happyHead, *currentBody, waifuName,
					assets.encouragements, 1*time.Second,
					func() { *encourageLocked = false })
			}
			return nil
		case rune(keys.Gift[0]):
			if !utils.LockGridChanges {
				utils.GiftMenu(ui.app, ui.grid, ui.actionSpace, ui.waifuArt, ui.chatBox,
					assets.head, assets.happyHead, waifuName, currentBody)
			}
			return nil
		case rune(keys.DressUp[0]):
			if !utils.LockGridChanges {
				utils.DressUp(ui.app, ui.grid, ui.actionSpace,ui.waifuArt,
					ui.chatBox, assets.head, waifuName, currentBody)
			}
			return nil
		case rune(keys.BackgroundMode[0]):
			utils.BackgroundMode(ui.app, ui.grid, ui.waifuArt, ui.chatBox, ui.happinessBar, ui.actionSpace, currentBody)
			return nil
		case rune(keys.Quit[0]):
			ui.app.Stop()
			return nil
		}

		// Vim-like navigation keys that work on any focused list
		if isVimNavigation {
			if focusedPrimitive := ui.app.GetFocus(); focusedPrimitive != nil {
				// Check if the focused primitive is a list to enable vim navigation
				if list, ok := focusedPrimitive.(*tview.List); ok {
					switch event.Rune() {
					case 'j':
						// Move down in the list
						currentItem := list.GetCurrentItem()
						if currentItem < list.GetItemCount()-1 {
							list.SetCurrentItem(currentItem + 1)
						}
						return nil
					case 'k':
						// Move up in the list
						currentItem := list.GetCurrentItem()
						if currentItem > 0 {
							list.SetCurrentItem(currentItem - 1)
						}
						return nil
					case 'l':
						// Select the current item
						currentIndex := list.GetCurrentItem()
						if currentIndex >= 0 && currentIndex < list.GetItemCount() {
							// Simulate Enter key to trigger the selected function
							return tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
						}
						return nil
					case 'h':
						// Go back/escape
						return tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone)
					}
				}
			}
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
	utils.HeadASCII      = &assets.head
	utils.BlinkHeadASCII = &assets.headBlink

	// ===== Set UI up
	// =====
	ui := createUI(assets)
	// Channel to handle UI changes without errors
	uiEvents := make(chan func(), 20)
	go func() {
		for fn := range uiEvents {
			ui.app.QueueUpdateDraw(fn)
		}
	}()
	utils.UIEventsChan = uiEvents
	// Create happiness bar's variable
	utils.HappinessBarRef = ui.happinessBar
	utils.GetHappinessBar()
	ui.happinessBar.SetText(utils.CurrentBar)

	// ===== Set palette up
	// =====
	// Get the palette
	palette, err := utils.LoadPalette()
	if err != nil {
		panic(fmt.Sprintf("Failed to load palette: %v", err))
	}
	// Apply palette to TextViews
	utils.ApplyTextViewPalette(palette, ui.happinessBar, ui.waifuArt, ui.chatBox)
	// Apply palette to Lists
	utils.ApplyListPalette(palette, ui.actionSpace)

	// ===== Set settings up
	// =====
	// Get the settings
	settings, err := utils.LoadSettings()
	if err != nil {
		panic(fmt.Sprintf("Failed to load settings: %v", err))
	}
	// Apply Waifu's name
	ui.waifuArt.SetTitle("| " + settings.Name + " |")
	// Apply deafult message in ChatBox
	ui.chatBox.SetText(settings.DefaultMessage)

	// ===== Variable work
	// =====
	var encourageLocked bool
	currentBody := assets.body

	// ===== Set functions
	// =====
	setupActionSpace(ui, assets, &encourageLocked, &currentBody, settings.Keys, settings.Name)
	setGlobalKeys(ui, assets, &encourageLocked, &currentBody, settings.Keys, settings.VimNavigation, settings.Name)

	// ===== Auto processes
	// =====
	ui.stopBlink = utils.StartBlinking(ui.app, ui.waifuArt,
		&assets.head, &assets.headBlink, &currentBody, 5*time.Second)

	// ===== No returns - Error handling
	// =====
	if err := utils.LoadClothes(utils.BasePath + "/clothes"); err != nil {
		panic(err)
	}
	if err := ui.app.SetRoot(ui.grid, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
	if err := utils.CreatePaletteFile(); err != nil {
		panic(err)
	}
}