package ui

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"

	"github.com/Cid-Emmerich/G1yph/internal/glyph"
	"github.com/Cid-Emmerich/G1yph/internal/paint"
	"github.com/Cid-Emmerich/G1yph/internal/theme"
)

func (a *App) handle(ev tcell.Event) {
	switch e := ev.(type) {
	case *tcell.EventResize:
		a.scr.Sync()
	case *tcell.EventMouse:
		a.mouse(e)
	case *tcell.EventKey:
		if a.picker {
			a.pickerKey(e)
			return
		}
		if a.help {
			a.helpKey(e)
			return
		}
		a.key(e)
	}
}

func (a *App) helpKey(e *tcell.EventKey) {
	switch e.Key() {
	case tcell.KeyEscape, tcell.KeyCtrlK, tcell.KeyEnter:
		a.help = false
	case tcell.KeyUp:
		a.helpScroll = max(0, a.helpScroll-1)
	case tcell.KeyDown:
		a.helpScroll++
	case tcell.KeyPgUp:
		a.helpScroll = max(0, a.helpScroll-10)
	case tcell.KeyPgDn:
		a.helpScroll += 10
	case tcell.KeyCtrlC:
		a.quit = true
	case tcell.KeyRune:
		switch e.Rune() {
		case 'q', '?':
			a.help = false
		case 'j':
			a.helpScroll++
		case 'k':
			a.helpScroll = max(0, a.helpScroll-1)
		}
	}
}

func (a *App) key(e *tcell.EventKey) {
	switch e.Key() {
	case tcell.KeyCtrlC:
		a.quit = true
	case tcell.KeyEscape:
		a.quit = true
	case tcell.KeyCtrlK:
		a.help = true
		a.helpScroll = 0
	case tcell.KeyRight, tcell.KeyEnter:
		a.step(1)
	case tcell.KeyLeft:
		a.step(-1)
	case tcell.KeyUp:
		a.stepGroup(-1)
	case tcell.KeyDown:
		a.stepGroup(1)
	case tcell.KeyRune:
		a.rune(e.Rune())
	}
}

func (a *App) rune(r rune) {
	switch r {
	case 'q':
		a.quit = true
	case '?':
		a.help = true
		a.helpScroll = 0
	case '/':
		a.openPicker()
	case ' ':
		a.paused = !a.paused
		if a.paused {
			a.say("paused")
		} else {
			a.say("playing")
		}
	case 'n':
		a.step(1)
	case 'b':
		a.step(-1)
	case ']':
		a.stepGroup(1)
	case '[':
		a.stepGroup(-1)
	case 'r':
		a.random()
	case 'a':
		a.auto = !a.auto
		a.autoAt = 0
		a.say(fmt.Sprintf("auto-cycle %s (%.0fs)", onOff(a.auto), a.autoSec))
	case '{':
		a.autoSec = clampf(a.autoSec-1, 1, 120)
		a.say(fmt.Sprintf("auto-cycle every %.0fs", a.autoSec))
	case '}':
		a.autoSec = clampf(a.autoSec+1, 1, 120)
		a.say(fmt.Sprintf("auto-cycle every %.0fs", a.autoSec))
	case 't', 'T':
		a.themeIdx = cyc(a.themeIdx, len(theme.Builtin), r == 't')
		a.th = theme.Builtin[a.themeIdx]
		a.say("theme " + a.th.Name)
	case 's', 'S':
		a.style = cyc(a.style, len(paint.StyleNames), r == 's')
		a.say("style " + paint.StyleNames[a.style])
	case 'c', 'C':
		a.charset = cyc(a.charset, len(paint.CharsetNames), r == 'c')
		if paint.StyleNames[a.style] != "ascii" {
			a.style = indexOf(paint.StyleNames, "ascii")
		}
		a.say("charset " + paint.CharsetNames[a.charset])
	case 'g', 'G':
		a.colour = cyc(a.colour, len(paint.ColourNames), r == 'g')
		a.say("colour " + paint.ColourNames[a.colour])
	case 'z':
		a.zoom = clampf(a.zoom+0.1, 0.3, 1.6)
		a.say(fmt.Sprintf("zoom %d%%", int(a.zoom*100+0.5)))
	case 'Z':
		a.zoom = clampf(a.zoom-0.1, 0.3, 1.6)
		a.say(fmt.Sprintf("zoom %d%%", int(a.zoom*100+0.5)))
	case '+', '=':
		a.speed = clampf(a.speed*1.25, 0.1, 6)
		a.say(fmt.Sprintf("speed %.2fx", a.speed))
	case '-', '_':
		a.speed = clampf(a.speed/1.25, 0.1, 6)
		a.say(fmt.Sprintf("speed %.2fx", a.speed))
	case 'l':
		a.label = LabelModes[cyc(indexOf(LabelModes, a.label), len(LabelModes), true)]
		a.say("label " + a.label)
	case 'd':
		a.bg = BackgroundModes[cyc(indexOf(BackgroundModes, a.bg), len(BackgroundModes), true)]
		a.say("backdrop " + a.bg)
	case 'o':
		a.border = !a.border
		a.say("border " + onOff(a.border))
	case 'x':
		a.mirror = !a.mirror
		a.say("mirror " + onOff(a.mirror))
	case 'R':
		a.style, a.charset, a.colour = 0, 0, 0
		a.zoom, a.speed = 1, 1
		a.label, a.bg = "name", "none"
		a.border, a.mirror = false, false
		a.say("look reset")
	}
}

func cyc(i, n int, fwd bool) int {
	if n == 0 {
		return 0
	}
	if fwd {
		return (i + 1) % n
	}
	return (i - 1 + n) % n
}

func onOff(b bool) string {
	if b {
		return "on"
	}
	return "off"
}

// ---------------------------------------------------------------------------
// Picker

func (a *App) openPicker() {
	a.picker = true
	a.query = ""
	a.pickCursor = 0
	a.refilter()
}

func (a *App) refilter() {
	a.pickRows = glyph.Search(a.query)
	if a.pickCursor >= len(a.pickRows) {
		a.pickCursor = max(0, len(a.pickRows)-1)
	}
}

func (a *App) pickerKey(e *tcell.EventKey) {
	switch e.Key() {
	case tcell.KeyEscape:
		a.picker = false
	case tcell.KeyCtrlC:
		a.quit = true
	case tcell.KeyEnter:
		if len(a.pickRows) > 0 {
			a.setGlyph(a.pickRows[a.pickCursor])
		}
		a.picker = false
	case tcell.KeyUp:
		a.pickCursor = max(0, a.pickCursor-1)
	case tcell.KeyDown:
		a.pickCursor = min(len(a.pickRows)-1, a.pickCursor+1)
	case tcell.KeyPgUp:
		a.pickCursor = max(0, a.pickCursor-10)
	case tcell.KeyPgDn:
		a.pickCursor = min(len(a.pickRows)-1, a.pickCursor+10)
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if a.query != "" {
			rs := []rune(a.query)
			a.query = string(rs[:len(rs)-1])
			a.refilter()
		}
	case tcell.KeyCtrlU:
		a.query = ""
		a.refilter()
	case tcell.KeyTab:
		// tab previews the highlighted glyph without closing
		if len(a.pickRows) > 0 {
			a.setGlyph(a.pickRows[a.pickCursor])
		}
	case tcell.KeyRune:
		r := e.Rune()
		if r == '/' && a.query == "" {
			return
		}
		if unicode.IsPrint(r) {
			a.query += string(r)
			a.refilter()
			// an exact emoji paste jumps straight there
			if idx := glyph.Find(a.query); idx >= 0 && strings.TrimSpace(a.query) != "" && !isASCII(a.query) {
				a.setGlyph(idx)
			}
		}
	}
	if a.pickCursor < 0 {
		a.pickCursor = 0
	}
}

func isASCII(s string) bool {
	for _, r := range s {
		if r > 127 {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// Mouse

func (a *App) mouse(e *tcell.EventMouse) {
	x, y := e.Position()
	btn := e.Buttons()
	switch {
	case btn&tcell.WheelDown != 0:
		if a.picker {
			a.pickCursor = min(len(a.pickRows)-1, a.pickCursor+1)
		} else if a.help {
			a.helpScroll++
		} else {
			a.step(1)
		}
	case btn&tcell.WheelUp != 0:
		if a.picker {
			a.pickCursor = max(0, a.pickCursor-1)
		} else if a.help {
			a.helpScroll = max(0, a.helpScroll-1)
		} else {
			a.step(-1)
		}
	case btn&tcell.Button1 != 0:
		if a.picker {
			bx, by, bw, bh := a.pickBox[0], a.pickBox[1], a.pickBox[2], a.pickBox[3]
			if x < bx || x >= bx+bw || y < by || y >= by+bh {
				a.picker = false
				return
			}
			row := y - (by + 3)
			if row >= 0 {
				first := a.pickFirst()
				if i := first + row; i < len(a.pickRows) {
					a.setGlyph(a.pickRows[i])
					a.picker = false
				}
			}
			return
		}
		if a.help {
			a.help = false
			return
		}
		if y == 0 {
			w, _ := a.scr.Size()
			if x >= w-12 {
				a.help = true
				return
			}
			a.openPicker()
			return
		}
		// click on the left third goes back, right two thirds forward
		w, _ := a.scr.Size()
		if x < w/3 {
			a.step(-1)
		} else {
			a.step(1)
		}
	}
}
