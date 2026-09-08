package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"

	"github.com/Cid-Emmerich/G1yph/internal/glyph"
	"github.com/Cid-Emmerich/G1yph/internal/paint"
)

func tc(c paint.RGB) tcell.Color { return tcell.NewRGBColor(int32(c.R), int32(c.G), int32(c.B)) }

func (a *App) st(fg paint.RGB) tcell.Style { return tcell.StyleDefault.Foreground(tc(fg)) }

// puts writes a string clipped to maxW cells (wide runes take two cells)
// and returns the width used.
func (a *App) puts(x, y int, s string, style tcell.Style, maxW int) int {
	i := 0
	for _, r := range s {
		w := runewidth.RuneWidth(r)
		if w == 0 {
			continue
		}
		if i+w > maxW {
			break
		}
		a.scr.SetContent(x+i, y, r, nil, style)
		i += w
	}
	return i
}

func strWidth(s string) int { return runewidth.StringWidth(s) }

func fit(s string, w int) string {
	if strWidth(s) <= w {
		return s
	}
	if w <= 1 {
		return ""
	}
	out := ""
	for _, r := range s {
		if strWidth(out)+runewidth.RuneWidth(r) > w-1 {
			break
		}
		out += string(r)
	}
	return out + "…"
}

func (a *App) fillRow(y, x0, x1 int, style tcell.Style) {
	for x := x0; x < x1; x++ {
		a.scr.SetContent(x, y, ' ', nil, style)
	}
}

// draw renders the whole screen.
func (a *App) draw() {
	a.scr.Clear()
	w, h := a.scr.Size()
	if w < 16 || h < 5 {
		a.puts(0, 0, "window too small", a.st(a.th.Warn), w)
		a.scr.Show()
		return
	}
	a.drawHeader(w)
	a.drawGlyph(w, h)
	a.drawStatus(w, h)
	if a.help {
		a.drawHelp(w, h)
	}
	if a.picker {
		a.drawPicker(w, h)
	}
	a.scr.Show()
}

func (a *App) drawHeader(w int) {
	def := glyph.Registry[a.idx]
	x := 1
	x += a.puts(x, 0, "g1yph", a.st(a.th.Accent).Bold(true), w)
	x += 2
	x += a.puts(x, 0, def.Emoji, a.st(a.th.Text), w-x)
	x += 1
	x += a.puts(x, 0, def.Name, a.st(a.th.Text).Bold(true), w-x)
	if def.Blurb != "" {
		x += a.puts(x, 0, "  ·  ", a.st(a.th.Select), w-x)
		x += a.puts(x, 0, fit(def.Blurb, w-x-14), a.st(a.th.Muted), w-x-14)
	}
	right := "ctrl+k help"
	a.puts(w-len(right)-1, 0, right, a.st(a.th.Muted), w)
}

// canvasRect returns the cell rectangle the animation draws in.
func (a *App) canvasRect(w, h int) (x0, y0, cw, ch int) {
	x0, y0 = 0, 1
	cw, ch = w, h-2
	if a.border {
		x0, y0, cw, ch = 1, 2, w-2, h-4
	}
	return
}

func (a *App) drawGlyph(w, h int) {
	x0, y0, cw, ch := a.canvasRect(w, h)
	if cw < 4 || ch < 2 {
		return
	}
	style := paint.StyleNames[a.style]
	bw, bh := paint.BufferSize(style, cw, ch)
	if a.bmp == nil || a.bmp.W != bw || a.bmp.H != bh {
		a.bmp = paint.NewBitmap(bw, bh)
	} else {
		a.bmp.Clear()
	}
	if a.cv == nil || a.cv.W != cw || a.cv.H != ch {
		a.cv = paint.NewCanvas(cw, ch)
	}
	// split off a strip at the bottom for the big label
	glyphH := bh
	if a.label == "big" && ch >= 8 {
		glyphH = int(float64(bh) * 0.80)
	}
	full := paint.NewPainter(a.bmp, 1, false)
	a.drawBackdrop(full)
	p := paint.NewPainterRegion(a.bmp, 0, 0, bw, glyphH, a.zoom, a.mirror)
	a.anim.Draw(p, a.t, a.frameDt)
	if glyphH < bh {
		a.drawBigLabel(bw, glyphH, bh)
	}
	paint.Render(a.bmp, a.cv, paint.RenderOpts{
		Style: style, Charset: paint.CharsetNames[a.charset], Colour: paint.ColourNames[a.colour],
		Palette: paint.Palette{Accent: a.th.Accent, Secondary: a.th.Secondary, Tertiary: a.th.Tertiary},
		Time: a.t,
	})
	for y := 0; y < ch; y++ {
		for x := 0; x < cw; x++ {
			c := a.cv.Cells[y*cw+x]
			if c.Ch == 0 {
				continue
			}
			s := tcell.StyleDefault.Foreground(tc(c.Fg))
			if c.HasBg {
				s = s.Background(tc(c.Bg))
			}
			a.scr.SetContent(x0+x, y0+y, c.Ch, nil, s)
		}
	}
	if a.border {
		a.drawBox(0, 1, w, h-2, a.st(a.th.Muted))
	}
	if a.label == "name" {
		def := glyph.Registry[a.idx]
		s := def.Emoji + " " + def.Name
		sw := strWidth(s)
		a.puts(x0+(cw-sw)/2, y0+ch-1, s, a.st(a.th.Text).Bold(true), cw)
	}
	if a.toast != "" && time.Now().Before(a.toastTill) {
		msg := " " + a.toast + " "
		a.puts(x0+cw-len(msg)-1, y0, msg, a.st(a.th.Text).Background(tc(a.th.Select)), cw)
	}
}

func (a *App) drawBackdrop(p *paint.Painter) {
	dim := paint.Mix(a.th.Select, paint.Black, 0.3)
	switch a.bg {
	case "stars":
		w, h := p.R-p.L, p.Bt-p.T
		for i := 0; i < 60; i++ {
			x := p.L + w*hash(i*7)
			y := p.T + h*hash(i*7+3)
			b := 0.5 + 0.5*math.Sin(a.t*(1+hash(i)*2)+hash(i+1)*6.28)
			p.Dot(x, y, 0.008, paint.Mix(dim, a.th.Muted, b))
		}
	case "grid":
		for x := math.Floor(p.L*10) / 10; x < p.R; x += 0.1 {
			p.Line(x, p.T, x, p.Bt, 0.004, dim)
		}
		for y := math.Floor(p.T*10) / 10; y < p.Bt; y += 0.1 {
			p.Line(p.L, y, p.R, y, 0.004, dim)
		}
	case "dots":
		for x := math.Floor(p.L*10) / 10; x < p.R; x += 0.1 {
			for y := math.Floor(p.T*10) / 10; y < p.Bt; y += 0.1 {
				p.Dot(x, y, 0.01, dim)
			}
		}
	}
}

func hash(i int) float64 {
	x := uint32(i)*2654435761 + 0x9E3779B9
	x ^= x >> 15
	x *= 0x2C1B3C6D
	x ^= x >> 12
	return float64(x%10000) / 10000
}

// drawBigLabel writes the glyph name in the block font under the picture.
func (a *App) drawBigLabel(bw, top, bh int) {
	name := strings.ToUpper(glyph.Registry[a.idx].Name)
	n := len([]rune(name))
	if n == 0 {
		return
	}
	stripH := float64(bh - top)
	sizePx := math.Min(stripH*0.7, float64(bw)*0.9*5/float64(4*n-1))
	if sizePx < 5 {
		return
	}
	p := paint.NewPainterRegion(a.bmp, 0, top, bw, bh, 1, false)
	size := sizePx / p.S
	cx := (p.L + p.R) / 2
	y := p.T + (stripH-sizePx)/2/p.S
	p.TextCentred(cx, y, size, name, a.th.Accent)
}

func (a *App) drawStatus(w, h int) {
	y := h - 1
	state := "▶"
	if a.paused {
		state = "⏸"
	}
	parts := []string{
		fmt.Sprintf("%s %.2gx", state, a.speed),
		"theme " + a.th.Name,
		"style " + paint.StyleNames[a.style],
	}
	if paint.StyleNames[a.style] == "ascii" {
		parts = append(parts, "charset "+paint.CharsetNames[a.charset])
	}
	parts = append(parts, "colour "+paint.ColourNames[a.colour], fmt.Sprintf("zoom %d%%", int(a.zoom*100+0.5)))
	if a.bg != "none" {
		parts = append(parts, "backdrop "+a.bg)
	}
	if a.mirror {
		parts = append(parts, "mirror")
	}
	if a.auto {
		parts = append(parts, fmt.Sprintf("auto %.0fs", a.autoSec))
	}
	x := 1
	for i, p := range parts {
		style := a.st(a.th.Muted)
		if i == 0 {
			style = a.st(a.th.Accent)
		}
		if x+len(p) > w-24 {
			break
		}
		x += a.puts(x, y, p, style, w-x)
		if i < len(parts)-1 {
			x += a.puts(x, y, " · ", a.st(a.th.Select), w-x)
		}
	}
	def := glyph.Registry[a.idx]
	right := fmt.Sprintf("%s %d/%d · / search", def.Group, a.idx+1, len(glyph.Registry))
	if len(right)+x+2 < w {
		a.puts(w-len(right)-1, y, right, a.st(a.th.Muted), w)
	}
}

// ---------------------------------------------------------------------------
// Overlays

func (a *App) drawBox(x, y, w, h int, style tcell.Style) {
	for i := 1; i < w-1; i++ {
		a.scr.SetContent(x+i, y, '─', nil, style)
		a.scr.SetContent(x+i, y+h-1, '─', nil, style)
	}
	for j := 1; j < h-1; j++ {
		a.scr.SetContent(x, y+j, '│', nil, style)
		a.scr.SetContent(x+w-1, y+j, '│', nil, style)
	}
	a.scr.SetContent(x, y, '╭', nil, style)
	a.scr.SetContent(x+w-1, y, '╮', nil, style)
	a.scr.SetContent(x, y+h-1, '╰', nil, style)
	a.scr.SetContent(x+w-1, y+h-1, '╯', nil, style)
}

func (a *App) drawHelp(w, h int) {
	lines := helpLines()
	bw := min(w-4, 92)
	bh := min(h-2, len(lines)+2)
	x0 := (w - bw) / 2
	y0 := (h - bh) / 2
	bg := tc(a.th.Select)
	base := tcell.StyleDefault.Background(bg).Foreground(tc(a.th.Text))
	for y := y0; y < y0+bh; y++ {
		a.fillRow(y, x0, x0+bw, base)
	}
	a.drawBox(x0, y0, bw, bh, base.Foreground(tc(a.th.Accent)))
	title := " g1yph shortcuts  (ctrl+k / esc to close, ↑↓ to scroll) "
	a.puts(x0+(bw-len(title))/2, y0, title, base.Foreground(tc(a.th.Accent)).Bold(true), bw)
	inner := bh - 2
	maxScroll := max(0, len(lines)-inner)
	if a.helpScroll > maxScroll {
		a.helpScroll = maxScroll
	}
	for i := 0; i < inner; i++ {
		li := a.helpScroll + i
		if li >= len(lines) {
			break
		}
		l := lines[li]
		y := y0 + 1 + i
		if strings.HasPrefix(l, "# ") {
			a.puts(x0+2, y, l[2:], base.Foreground(tc(a.th.Secondary)).Bold(true), bw-4)
			continue
		}
		key, desc, ok := strings.Cut(l, "\t")
		if !ok {
			a.puts(x0+2, y, l, base, bw-4)
			continue
		}
		const keyW = 26
		a.puts(x0+3, y, fit(key, keyW-1), base.Foreground(tc(a.th.Accent)), keyW-1)
		a.puts(x0+3+keyW, y, fit(strings.TrimSpace(desc), bw-keyW-5), base, bw-keyW-5)
	}
	if maxScroll > 0 {
		pct := fmt.Sprintf(" %d/%d ", a.helpScroll+inner, len(lines))
		a.puts(x0+bw-len(pct)-2, y0+bh-1, pct, base.Foreground(tc(a.th.Muted)), bw)
	}
}

// pickFirst returns the first visible row index so the cursor stays on screen.
func (a *App) pickFirst() int {
	_, h := a.scr.Size()
	rows := a.pickerRows(h)
	if rows <= 0 {
		return 0
	}
	first := 0
	if a.pickCursor >= rows {
		first = a.pickCursor - rows + 1
	}
	return first
}

func (a *App) pickerRows(h int) int {
	bh := min(h-2, 24)
	return bh - 4
}

func (a *App) drawPicker(w, h int) {
	bw := min(w-4, 64)
	bh := min(h-2, 24)
	x0 := (w - bw) / 2
	y0 := (h - bh) / 2
	a.pickBox = [4]int{x0, y0, bw, bh}
	bg := tc(a.th.Select)
	base := tcell.StyleDefault.Background(bg).Foreground(tc(a.th.Text))
	for y := y0; y < y0+bh; y++ {
		a.fillRow(y, x0, x0+bw, base)
	}
	a.drawBox(x0, y0, bw, bh, base.Foreground(tc(a.th.Accent)))
	title := fmt.Sprintf(" pick a glyph  (%d) ", len(glyph.Registry))
	a.puts(x0+(bw-len(title))/2, y0, title, base.Foreground(tc(a.th.Accent)).Bold(true), bw)
	// search line
	a.puts(x0+2, y0+1, "> ", base.Foreground(tc(a.th.Accent)), bw-4)
	q := a.query
	if q == "" {
		a.puts(x0+4, y0+1, "type a name, or paste an emoji", base.Foreground(tc(a.th.Muted)), bw-6)
	} else {
		a.puts(x0+4, y0+1, q, base.Bold(true), bw-6)
	}
	a.scr.ShowCursor(x0+4+strWidth(q), y0+1)
	a.puts(x0+2, y0+2, strings.Repeat("─", bw-4), base.Foreground(tc(a.th.Muted)), bw-4)
	rows := a.pickerRows(h)
	first := a.pickFirst()
	for i := 0; i < rows; i++ {
		ri := first + i
		if ri >= len(a.pickRows) {
			break
		}
		def := glyph.Registry[a.pickRows[ri]]
		y := y0 + 3 + i
		style := base
		if ri == a.pickCursor {
			style = tcell.StyleDefault.Background(tc(a.th.Accent)).Foreground(tc(a.th.Select)).Bold(true)
			a.fillRow(y, x0+1, x0+bw-1, style)
		}
		x := x0 + 2
		x += a.puts(x, y, def.Emoji, style, 3)
		x = x0 + 5
		x += a.puts(x, y, def.Name, style, 18)
		x = x0 + 24
		desc := def.Blurb
		if ri == a.pickCursor {
			a.puts(x, y, fit(desc, bw-26), style, bw-26)
		} else {
			a.puts(x, y, fit(desc, bw-26-10), style.Foreground(tc(a.th.Muted)), bw-26-10)
			a.puts(x0+bw-2-len(def.Group), y, def.Group, style.Foreground(tc(a.th.Muted)), bw)
		}
	}
	if len(a.pickRows) == 0 {
		a.puts(x0+2, y0+3, "nothing matches", base.Foreground(tc(a.th.Warn)), bw-4)
	}
	foot := " enter pick · tab preview · esc close "
	a.puts(x0+bw-len(foot)-1, y0+bh-1, foot, base.Foreground(tc(a.th.Muted)), bw)
	if len(a.pickRows) > rows {
		pos := fmt.Sprintf(" %d/%d ", a.pickCursor+1, len(a.pickRows))
		a.puts(x0+1, y0+bh-1, pos, base.Foreground(tc(a.th.Muted)), bw)
	}
}
