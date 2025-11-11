package utils

import (
	"os"
	"fmt"
	"time"
	"math/rand"
	"path/filepath"

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
		app.QueueUpdateDraw(func() {
			chatBox.SetText("Waifu: " + line)
			waifuArt.SetText(happyHead + "\n" + body)
		})

		time.Sleep(duration)

		// Restore neutral
		app.QueueUpdateDraw(func() {
			waifuArt.SetText(head + "\n" + body)
			unlockFunc()
		})
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
		app.QueueUpdateDraw(func() {
			chatBox.SetText("No clothes found!")
		})
		return
	}

	list := tview.NewList()
	ApplyListPalette(cachedPalette, list)  // Safe: the error could have occurred during the initialization
	for _, item := range clothesCache {
		display := "-" + item.Name
		list.AddItem(display, "", 0, func() {
			*currentBody = item.Data
			waifuArt.SetText(head + "\n" + *currentBody)
			chatBox.SetText("Waifu changed into: " + item.Name)

			closeDressUp(app, grid, list, actionSpace, waifuArt, head, blinkHead, currentBody)
		})
	}

	list.SetBorder(true).SetTitle("| Dress Up |").SetTitleAlign(tview.AlignCenter)
	list.SetDoneFunc(func() {
		closeDressUp(app, grid, list, actionSpace, waifuArt, head, blinkHead, currentBody)
	})

	// Swap in the dress-up list
	grid.RemoveItem(actionSpace)
	grid.AddItem(list, 0, 0, 2, 1, 0, 0, true)
	app.SetFocus(list)
}

// scanASCIIFiles recursively scans directory and returns paths and display names
func scanASCIIFiles(dir string) ([]string, []string) {
	var files []string
	var names []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		files = append(files, path)
		rel, _ := filepath.Rel(dir, path)
		names = append(names, rel)
		return nil
	})
	return files, names
}

// LoadClothes loads all clothes ASCII files from the specified directory into the cache
func LoadClothes(dir string) error {
	files, names := scanASCIIFiles(dir)
	if len(files) == 0 {
		return fmt.Errorf("no clothes found")
	}

	clothesCache = make([]struct {
		Name string
		Data string
	}, len(files))

	for i, f := range files {
		data, err := LoadASCII(f)
		if err != nil {
			return fmt.Errorf("failed to load %s: %w", f, err)
		}
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
	grid.AddItem(actionSpace, 0, 0, 2, 1, 0, 0, true)
	app.SetFocus(actionSpace)
}
