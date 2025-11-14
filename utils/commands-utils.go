package utils

import (
	"fmt"
	"math/rand"
	"path"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ==============================
// ENCOURAGEMENT
// ==============================

// Encourage shows a random encouragement, swaps head for `duration`, then restores
func Encourage(app *tview.Application, waifuArt, chatBox *tview.TextView,
	head, happyHead, body string, encouragements []string,
	duration time.Duration, unlockFunc func()) {

	if len(encouragements) == 0 {
		unlockFunc()
		return
	}

	go func() {
		rand.Seed(time.Now().UnixNano())
		line := encouragements[rand.Intn(len(encouragements))]

		// Show happy face + message
		if UIEventsChan != nil {
			UIEventsChan <- func() {
				chatBox.SetText("Waifu: " + line)
				waifuArt.SetText(happyHead + "\n" + body)
				IncreaseHappiness(6)
			}
		}

		time.Sleep(duration)

		// Restore neutral
		if UIEventsChan != nil {
			UIEventsChan <- func() {
				waifuArt.SetText(head + "\n" + body)
				unlockFunc()
			}
		}
	}()
}

// ==============================
// DRESS-UP
// ==============================

var clothesCache []struct {
	Name string
	Data string
}

// DressUp allows the user to pick a clothes ASCII file from a scrollable list
func DressUp(app *tview.Application, waifuArt, chatBox *tview.TextView,
	head, blinkHead string, grid *tview.Grid, actionSpace *tview.List,
	currentBody *string) {

	if len(clothesCache) == 0 {
		if UIEventsChan != nil {
			UIEventsChan <- func() {
				chatBox.SetText("No clothes found!")
			}
		}
		return
	}

	list := tview.NewList()
	ApplyListPalette(cachedPalette, list) // Safe: the error could have occurred during the initialization
	for _, item := range clothesCache {
		display := "-" + item.Name
		list.AddItem(display, "", 0, func() {
			if UIEventsChan != nil {
				UIEventsChan <- func() {
					*currentBody = item.Data
					waifuArt.SetText(head + "\n" + *currentBody)
					chatBox.SetText("Waifu changed into: " + item.Name)
					IncreaseHappiness(3)
				}
			}

			closeDressUp(app, grid, list, actionSpace, waifuArt, head, blinkHead, currentBody)
		})
	}

	list.SetBorder(true).SetTitle("| Dress Up |").SetTitleAlign(tview.AlignCenter)
	list.SetDoneFunc(func() {
		closeDressUp(app, grid, list, actionSpace, waifuArt, head, blinkHead, currentBody)
	})

	// Swap in the dress-up list
	grid.RemoveItem(actionSpace)
	grid.AddItem(list, 0, 0, 1, 1, 0, 0, true)
	app.SetFocus(list)
}

// scanASCIIFiles recursively scans directory and returns paths and display names
func scanASCIIFiles(dir string) ([]string, []string, error) {
	var files []string
	var names []string

	var walk func(string, string) error
	walk = func(currentPath, relPath string) error {
		entries, err := ASCIIFS.ReadDir(currentPath)
		if err != nil {
			return err
		}
		for _, e := range entries {
			fullPath := path.Join(currentPath, e.Name())
			rel := path.Join(relPath, e.Name())
			if e.IsDir() {
				if err := walk(fullPath, rel); err != nil {
					return err
				}
			} else {
				files = append(files, fullPath)
				names = append(names, rel)
			}
		}
		return nil
	}

	if err := walk(dir, ""); err != nil {
		return nil, nil, err
	}
	return files, names, nil
}

// LoadClothes loads all clothes ASCII files from the specified directory into the cache
func LoadClothes(dir string) error {
	files, names, err := scanASCIIFiles(dir)
	if err != nil {
		return fmt.Errorf("failed to scan clothes: %v", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no clothes found")
	}

	clothesCache = make([]struct {
		Name string
		Data string
	}, len(files))

	for i, f := range files {
		data := LoadASCII(f)
		clothesCache[i] = struct {
			Name string
			Data string
		}{names[i], data}
	}

	return nil
}

// closeDressUp restores the actionSpace and restarts blinking
func closeDressUp(app *tview.Application, grid *tview.Grid, list *tview.List,
	actionSpace *tview.List, waifuArt *tview.TextView,
	head, blinkHead string, currentBody *string) {

	grid.RemoveItem(list)
	grid.AddItem(actionSpace, 0, 0, 1, 1, 0, 0, true)
	app.SetFocus(actionSpace)
}

// ==============================
// BACKGROUND MODE
// ==============================

var LockGridChanges bool = false

// BackgroundMode makes the UI focus only on the waifuArt view
func BackgroundMode(app *tview.Application, waifuArt, chatBox, happinessBar *tview.TextView,
	grid *tview.Grid, actionSpace *tview.List, currentBody *string) {

	// Block some keys
	LockGridChanges = true

	// Remove odds and show waifuArt only
	grid.RemoveItem(actionSpace)
	grid.RemoveItem(happinessBar)
	grid.RemoveItem(chatBox)
	grid.RemoveItem(waifuArt)

	grid.AddItem(waifuArt, 0, 0, 2, 2, 0, 34, true)
	app.SetFocus(waifuArt)

	// Output Handler (Esc)
	waifuArt.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			closeBackground(app, grid, waifuArt, actionSpace, happinessBar, chatBox)
		}
	})
}

// closeBackground restores the previous layout after BackgroundMode
func closeBackground(app *tview.Application, grid *tview.Grid,
	waifuArt *tview.TextView, actionSpace *tview.List,
	happinessBar, chatBox *tview.TextView) {

	grid.RemoveItem(waifuArt)

	grid.AddItem(actionSpace, 0, 0, 1, 1, 0, 0, true)
	grid.AddItem(happinessBar, 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(waifuArt, 0, 1, 1, 1, 0, 75, false)
	grid.AddItem(chatBox, 1, 1, 1, 1, 0, 0, false)

	app.SetFocus(actionSpace)
	waifuArt.SetDoneFunc(nil)

	LockGridChanges = false
}
