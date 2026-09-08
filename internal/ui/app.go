// Package ui is the terminal interface: the animation view, the picker
// and the help overlay.
package ui

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/Cid-Emmerich/G1yph/internal/config"
	"github.com/Cid-Emmerich/G1yph/internal/glyph"
	"github.com/Cid-Emmerich/G1yph/internal/paint"
	"github.com/Cid-Emmerich/G1yph/internal/theme"
)

// LabelModes and BackgroundModes list the cycling orders.
var (
	LabelModes      = []string{"name", "big", "off"}
	BackgroundModes = []string{"none", "stars", "grid", "dots"}
)

// App is the whole interactive program.
type App struct {
	scr tcell.Screen
	cfg *config.Config

	th       theme.Theme
	themeIdx int

	idx     int // registry index of the current glyph
	anim    glyph.Anim
	t       float64 // animation clock
	paused  bool
	speed   float64
	zoom    float64
	style   int
	charset int
	colour  int
	label   string
	bg      string
	border  bool
	mirror  bool
	auto    bool
	autoAt  float64
	autoSec float64

	bmp *paint.Bitmap
	cv  *paint.Canvas

	help       bool
	helpScroll int

	picker     bool
	query      string
	pickCursor int
	pickRows   []int
	pickBox    [4]int // x, y, w, h of the picker box for mouse hits

	toast     string
	toastTill time.Time
	quit      bool
	frameDt   float64
	rng       *rand.Rand
}

// New builds the app for a starting glyph index (-1 opens the picker on a
// random glyph).
func New(cfg *config.Config, start int) *App {
	a := &App{cfg: cfg, rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
	a.applyConfig()
	if start < 0 {
		start = a.rng.Intn(len(glyph.Registry))
		a.picker = true
	}
	a.setGlyph(start)
	return a
}

func (a *App) applyConfig() {
	c := a.cfg
	a.themeIdx = theme.Index(c.Theme)
	a.th = theme.Builtin[a.themeIdx]
	a.style = indexOf(paint.StyleNames, c.Style)
	a.charset = indexOf(paint.CharsetNames, c.Charset)
	a.colour = indexOf(paint.ColourNames, c.Colour)
	a.zoom = clampf(c.Zoom, 0.3, 1.6)
	a.speed = clampf(c.Speed, 0.1, 6)
	a.label = c.Label
	if indexOf(LabelModes, a.label) < 0 {
		a.label = "name"
	}
	a.bg = c.Background
	if indexOf(BackgroundModes, a.bg) < 0 {
		a.bg = "none"
	}
	a.border = c.Border
	a.mirror = c.Mirror
	a.auto = c.Auto
	a.autoSec = clampf(c.AutoSeconds, 1, 120)
}

func (a *App) saveConfig() {
	c := a.cfg
	c.Theme = a.th.Name
	c.Style = paint.StyleNames[a.style]
	c.Charset = paint.CharsetNames[a.charset]
	c.Colour = paint.ColourNames[a.colour]
	c.Zoom, c.Speed = a.zoom, a.speed
	c.Label, c.Background = a.label, a.bg
	c.Border, c.Mirror, c.Auto = a.border, a.mirror, a.auto
	c.AutoSeconds = a.autoSec
	c.Last = glyph.Registry[a.idx].Name
	_ = c.Save()
}

func indexOf(list []string, s string) int {
	for i, x := range list {
		if x == s {
			return i
		}
	}
	return -1
}

func clampf(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

// setGlyph switches to a registry index and restarts its animation.
func (a *App) setGlyph(i int) {
	n := len(glyph.Registry)
	i = ((i % n) + n) % n
	a.idx = i
	a.anim = glyph.Registry[i].New()
	a.t = 0
	a.autoAt = 0
}

func (a *App) step(d int) { a.setGlyph(a.idx + d) }

func (a *App) random() {
	if len(glyph.Registry) < 2 {
		return
	}
	i := a.rng.Intn(len(glyph.Registry) - 1)
	if i >= a.idx {
		i++
	}
	a.setGlyph(i)
}

// stepGroup jumps to the first glyph of the next/previous group.
func (a *App) stepGroup(d int) {
	groups := glyph.Groups()
	cur := glyph.Registry[a.idx].Group
	gi := 0
	for i, g := range groups {
		if g == cur {
			gi = i
		}
	}
	gi = ((gi+d)%len(groups) + len(groups)) % len(groups)
	for i, def := range glyph.Registry {
		if def.Group == groups[gi] {
			a.setGlyph(i)
			return
		}
	}
}

func (a *App) say(msg string) {
	a.toast = msg
	a.toastTill = time.Now().Add(1800 * time.Millisecond)
}

// Run is the main loop.
func (a *App) Run() error {
	scr, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := scr.Init(); err != nil {
		return err
	}
	a.scr = scr
	scr.EnableMouse()
	scr.HideCursor()
	defer func() {
		scr.Fini()
		a.saveConfig()
	}()

	events := make(chan tcell.Event, 32)
	go func() {
		for {
			ev := scr.PollEvent()
			if ev == nil {
				close(events)
				return
			}
			events <- ev
		}
	}()

	fps := a.cfg.FPS
	if fps < 5 || fps > 120 {
		fps = 30
	}
	tick := time.NewTicker(time.Second / time.Duration(fps))
	defer tick.Stop()
	last := time.Now()
	a.draw()
	for !a.quit {
		select {
		case ev, ok := <-events:
			if !ok {
				return nil
			}
			a.handle(ev)
			// keep the animation smooth by not redrawing here; the tick does it
		case now := <-tick.C:
			dt := now.Sub(last).Seconds()
			last = now
			if dt > 0.1 {
				dt = 0.1
			}
			if !a.paused {
				dt *= a.speed
				a.t += dt
				if a.auto && !a.picker && !a.help {
					a.autoAt += dt / a.speed
					if a.autoAt >= a.autoSec {
						a.step(1)
					}
				}
			} else {
				dt = 0
			}
			a.frameDt = dt
			a.draw()
		}
	}
	return nil
}

// ClosePicker starts on the current glyph without the picker overlay.
func (a *App) ClosePicker() { a.picker = false }
