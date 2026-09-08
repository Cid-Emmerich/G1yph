package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"strings"

	"github.com/Cid-Emmerich/G1yph/internal/config"
	"github.com/Cid-Emmerich/G1yph/internal/glyph"
	"github.com/Cid-Emmerich/G1yph/internal/paint"
	"github.com/Cid-Emmerich/G1yph/internal/theme"
)

// renderCanvas draws a glyph at time t into a fresh canvas.
func renderCanvas(cfg *config.Config, idx int, w, h int, t float64) *paint.Canvas {
	def := glyph.Registry[idx]
	anim := def.New()
	style := cfg.Style
	if indexOf(paint.StyleNames, style) < 0 {
		style = "blocks"
	}
	bw, bh := paint.BufferSize(style, w, h)
	bmp := paint.NewBitmap(bw, bh)
	cv := paint.NewCanvas(w, h)
	th := theme.Get(cfg.Theme)
	const step = 1.0 / 30
	clock := 0.0
	for clock < t {
		bmp.Clear()
		p := paint.NewPainter(bmp, cfg.Zoom, cfg.Mirror)
		anim.Draw(p, clock, step)
		clock += step
	}
	bmp.Clear()
	p := paint.NewPainter(bmp, cfg.Zoom, cfg.Mirror)
	anim.Draw(p, t, step)
	paint.Render(bmp, cv, paint.RenderOpts{Style: style, Charset: cfg.Charset, Colour: cfg.Colour,
		Palette: paint.Palette{Accent: th.Accent, Secondary: th.Secondary, Tertiary: th.Tertiary}, Time: t})
	return cv
}

// cellW and cellH are the pixel size of one terminal cell in PNG output.
const cellW, cellH = 8, 16

// paintCells draws a canvas into img at pixel offset ox,oy, approximating
// how a terminal would show it.
func paintCells(img *image.RGBA, cv *paint.Canvas, ox, oy int) {
	fill := func(x0, y0, x1, y1 int, c paint.RGB) {
		for y := y0; y < y1; y++ {
			for x := x0; x < x1; x++ {
				img.Set(ox+x, oy+y, color.RGBA{c.R, c.G, c.B, 255})
			}
		}
	}
	for y := 0; y < cv.H; y++ {
		for x := 0; x < cv.W; x++ {
			c := cv.Cells[y*cv.W+x]
			px, py := x*cellW, y*cellH
			if c.Ch == 0 {
				continue
			}
			if c.HasBg {
				fill(px, py, px+cellW, py+cellH, c.Bg)
			}
			switch c.Ch {
			case '▀':
				fill(px, py, px+cellW, py+cellH/2, c.Fg)
			case '▄':
				fill(px, py+cellH/2, px+cellW, py+cellH, c.Fg)
			case '█':
				fill(px, py, px+cellW, py+cellH, c.Fg)
			default:
				// text-ish glyphs: a smaller mark in the middle of the cell
				fill(px+2, py+4, px+cellW-2, py+cellH-4, c.Fg)
			}
		}
	}
}

// FramePNG writes one frame as a PNG.
func FramePNG(out io.Writer, cfg *config.Config, idx int, w, h int, t float64) error {
	cv := renderCanvas(cfg, idx, w, h, t)
	img := image.NewRGBA(image.Rect(0, 0, w*cellW, h*cellH))
	fillBg(img, 0, 0, w*cellW, h*cellH)
	paintCells(img, cv, 0, 0)
	return png.Encode(out, img)
}

func fillBg(img *image.RGBA, x0, y0, x1, y1 int) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			img.Set(x, y, color.RGBA{18, 18, 22, 255})
		}
	}
}

// SheetPNG writes a contact sheet of every glyph, cols per row, each
// tile w x h cells, at time t.
func SheetPNG(out io.Writer, cfg *config.Config, cols, w, h int, t float64) error {
	n := len(glyph.Registry)
	rows := (n + cols - 1) / cols
	img := image.NewRGBA(image.Rect(0, 0, cols*w*cellW, rows*h*cellH))
	fillBg(img, 0, 0, cols*w*cellW, rows*h*cellH)
	for i := range glyph.Registry {
		cv := renderCanvas(cfg, i, w, h, t)
		ox, oy := (i%cols)*w*cellW, (i/cols)*h*cellH
		paintCells(img, cv, ox, oy)
		// a faint tile border
		for x := 0; x < w*cellW; x++ {
			img.Set(ox+x, oy, color.RGBA{50, 50, 60, 255})
		}
		for y := 0; y < h*cellH; y++ {
			img.Set(ox, oy+y, color.RGBA{50, 50, 60, 255})
		}
	}
	return png.Encode(out, img)
}

// Frame renders one frame of a glyph at time t into w x h cells and
// returns it as ANSI true-colour text (colour=false gives plain text).
// It runs the animation up to t in small steps so particle systems have
// had time to fill in.
func Frame(cfg *config.Config, idx int, w, h int, t float64, colour bool) string {
	cv := renderCanvas(cfg, idx, w, h, t)
	style := cfg.Style
	var sb strings.Builder
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := cv.Cells[y*w+x]
			if c.Ch == 0 {
				sb.WriteByte(' ')
				continue
			}
			if colour {
				fmt.Fprintf(&sb, "\x1b[38;2;%d;%d;%dm", c.Fg.R, c.Fg.G, c.Fg.B)
				if c.HasBg {
					fmt.Fprintf(&sb, "\x1b[48;2;%d;%d;%dm", c.Bg.R, c.Bg.G, c.Bg.B)
				}
				sb.WriteRune(c.Ch)
				sb.WriteString("\x1b[0m")
			} else {
				if style == "blocks" && c.HasBg {
					sb.WriteRune('█')
				} else {
					sb.WriteRune(c.Ch)
				}
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
