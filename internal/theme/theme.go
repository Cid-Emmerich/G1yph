// Package theme defines the colour themes used by the UI chrome and by the
// "theme" colour mode of the renderer. The terminal background is never
// painted, so G1yph blends into whatever the terminal already uses.
package theme

import (
	"strings"

	"github.com/Cid-Emmerich/G1yph/internal/paint"
)

// Theme is a small palette.
type Theme struct {
	Name      string
	Text      paint.RGB
	Muted     paint.RGB
	Accent    paint.RGB
	Secondary paint.RGB
	Tertiary  paint.RGB
	Select    paint.RGB
	Warn      paint.RGB
}

func rgb(h uint32) paint.RGB { return paint.Hex(h) }

// Builtin themes in cycling order. The first is the default.
var Builtin = []Theme{
	{Name: "g1yph", Text: rgb(0xF2EFEA), Muted: rgb(0x6F6C7A), Accent: rgb(0xFFB454), Secondary: rgb(0xFF6E9C), Tertiary: rgb(0x7CE0FF), Select: rgb(0x33303F), Warn: rgb(0xFF5C5C)},
	{Name: "wvfrm", Text: rgb(0xE6E6F0), Muted: rgb(0x6C6F85), Accent: rgb(0x7AA2F7), Secondary: rgb(0xBB9AF7), Tertiary: rgb(0x7DCFFF), Select: rgb(0x2A2E45), Warn: rgb(0xF7768E)},
	{Name: "nord", Text: rgb(0xECEFF4), Muted: rgb(0x4C566A), Accent: rgb(0x88C0D0), Secondary: rgb(0x81A1C1), Tertiary: rgb(0xB48EAD), Select: rgb(0x3B4252), Warn: rgb(0xBF616A)},
	{Name: "dracula", Text: rgb(0xF8F8F2), Muted: rgb(0x6272A4), Accent: rgb(0xBD93F9), Secondary: rgb(0xFF79C6), Tertiary: rgb(0x8BE9FD), Select: rgb(0x44475A), Warn: rgb(0xFF5555)},
	{Name: "gruvbox", Text: rgb(0xEBDBB2), Muted: rgb(0x928374), Accent: rgb(0xFABD2F), Secondary: rgb(0xFE8019), Tertiary: rgb(0xB8BB26), Select: rgb(0x3C3836), Warn: rgb(0xFB4934)},
	{Name: "catppuccin", Text: rgb(0xCDD6F4), Muted: rgb(0x6C7086), Accent: rgb(0xCBA6F7), Secondary: rgb(0xF5C2E7), Tertiary: rgb(0x89DCEB), Select: rgb(0x313244), Warn: rgb(0xF38BA8)},
	{Name: "solarized", Text: rgb(0xEEE8D5), Muted: rgb(0x586E75), Accent: rgb(0x2AA198), Secondary: rgb(0x268BD2), Tertiary: rgb(0xB58900), Select: rgb(0x073642), Warn: rgb(0xDC322F)},
	{Name: "synthwave", Text: rgb(0xF4EEFF), Muted: rgb(0x7A5C99), Accent: rgb(0xFF2ED2), Secondary: rgb(0x2DE2E6), Tertiary: rgb(0xF6F740), Select: rgb(0x3A1F5D), Warn: rgb(0xFF6E6E)},
	{Name: "sunset", Text: rgb(0xFFF1E6), Muted: rgb(0x8C6A5D), Accent: rgb(0xFF7B54), Secondary: rgb(0xFFB26B), Tertiary: rgb(0xFFD56F), Select: rgb(0x4A2B2B), Warn: rgb(0xFF4C4C)},
	{Name: "forest", Text: rgb(0xE3EBD9), Muted: rgb(0x6B7F62), Accent: rgb(0x8FD694), Secondary: rgb(0x4FA37A), Tertiary: rgb(0xC9E4A6), Select: rgb(0x2C3A2A), Warn: rgb(0xE07A5F)},
	{Name: "ocean", Text: rgb(0xE0F2FE), Muted: rgb(0x557A95), Accent: rgb(0x38BDF8), Secondary: rgb(0x818CF8), Tertiary: rgb(0x34D399), Select: rgb(0x1E3A5F), Warn: rgb(0xFB7185)},
	{Name: "mono", Text: rgb(0xF0F0F0), Muted: rgb(0x707070), Accent: rgb(0xFFFFFF), Secondary: rgb(0xBDBDBD), Tertiary: rgb(0x8A8A8A), Select: rgb(0x333333), Warn: rgb(0xFFFFFF)},
	{Name: "amber", Text: rgb(0xFFE7B3), Muted: rgb(0x8A6A2E), Accent: rgb(0xFFB000), Secondary: rgb(0xFF8C00), Tertiary: rgb(0xFFD866), Select: rgb(0x3A2A0A), Warn: rgb(0xFF5A36)},
	{Name: "matrix", Text: rgb(0xCCFFCC), Muted: rgb(0x2E6B2E), Accent: rgb(0x00FF41), Secondary: rgb(0x00B32C), Tertiary: rgb(0x9CFFB0), Select: rgb(0x0B2E12), Warn: rgb(0xFF3B3B)},
}

// Names lists every theme name.
func Names() []string {
	out := make([]string, len(Builtin))
	for i, t := range Builtin {
		out[i] = t.Name
	}
	return out
}

// Get returns a theme by name (the default when unknown).
func Get(name string) Theme {
	name = strings.ToLower(name)
	for _, t := range Builtin {
		if t.Name == name {
			return t
		}
	}
	return Builtin[0]
}

// Index returns the position of a theme name (0 when unknown).
func Index(name string) int {
	for i, t := range Builtin {
		if t.Name == strings.ToLower(name) {
			return i
		}
	}
	return 0
}
