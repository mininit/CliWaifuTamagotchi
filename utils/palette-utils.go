package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ==============================
// PALETTE STRUCT
// ==============================

// Palette stores UI color values
type Palette struct {
	Background string `json:"background"`
	Foreground string `json:"foreground"`
	Border     string `json:"border"`
	Accent     string `json:"accent"`
	Title      string `json:"title"`
}

// cachedPalette stores the palette loaded for this session
var cachedPalette *Palette

// ==============================
// DEFAULT PALETTE
// ==============================

// DefaultPalette returns a default color palette
func DefaultPalette() *Palette {
	return &Palette{
		Background: "#1e1e2e",
		Foreground: "#cdd6f4",
		Border:     "#cba6f7",
		Accent:     "#eba0ac",
		Title:      "#b4befe",
	}
}

// ==============================
// FILE HANDLING
// ==============================

// CreatePaletteFile creates palette.json in ~/.config/cliwaifutamagotchi if missing
func CreatePaletteFile() error {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "cliwaifutamagotchi")
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	palettePath := filepath.Join(configDir, "palette.json")
	if _, err := os.Stat(palettePath); err == nil {
		return nil
	}

	file, err := os.Create(palettePath)
	if err != nil {
		return fmt.Errorf("failed to create palette file: %w", err)
	}
	defer file.Close()

	palette := DefaultPalette()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(palette); err != nil {
		return fmt.Errorf("failed to write palette file: %w", err)
	}

	return nil
}

// LoadPalette loads the palette from file (or default if missing), cached for session
func LoadPalette() (*Palette, error) {
	if cachedPalette != nil {
		return cachedPalette, nil
	}

	configDir := filepath.Join(os.Getenv("HOME"), ".config", "cliwaifutamagotchi")
	palettePath := filepath.Join(configDir, "palette.json")

	if _, err := os.Stat(palettePath); os.IsNotExist(err) {
		if err := CreatePaletteFile(); err != nil {
			return nil, err
		}
	}

	file, err := os.Open(palettePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open palette file: %w", err)
	}
	defer file.Close()

	var p Palette
	if err := json.NewDecoder(file).Decode(&p); err != nil {
		return nil, fmt.Errorf("failed to decode palette file: %w", err)
	}

	cachedPalette = &p
	return cachedPalette, nil
}

// ==============================
// APPLY PALETTE TO WIDGETS
// ==============================

// ApplyTextViewPalette sets colors on one or more TextView widgets
func ApplyTextViewPalette(p *Palette, views ...*tview.TextView) {
	bgColor := tcell.GetColor(p.Background)
	fgColor := tcell.GetColor(p.Foreground)
	borderColor := tcell.GetColor(p.Border)
	titleColor := tcell.GetColor(p.Title)

	for _, v := range views {
		v.SetBackgroundColor(bgColor)
		v.SetTextColor(fgColor)
		v.SetBorder(true)
		v.SetBorderColor(borderColor)
		v.SetTitleColor(titleColor)
	}
}

// ApplyListPalette sets colors on one or more List widgets
func ApplyListPalette(p *Palette, lists ...*tview.List) {
	bgColor := tcell.GetColor(p.Background)
	fgColor := tcell.GetColor(p.Foreground)
	borderColor := tcell.GetColor(p.Border)
	titleColor := tcell.GetColor(p.Title)
	highlightColor := tcell.GetColor(p.Accent)

	for _, l := range lists {
		l.SetBackgroundColor(bgColor)
		l.SetMainTextColor(fgColor)
		l.SetSecondaryTextColor(fgColor)
		l.SetSelectedTextColor(fgColor)
		l.SetSelectedBackgroundColor(highlightColor)
		l.SetBorder(true)
		l.SetBorderColor(borderColor)
		l.SetTitleColor(titleColor)
	}
}
