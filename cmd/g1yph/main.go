// g1yph – an emoji animator for the terminal.
//
//	g1yph                       open with the picker
//	g1yph <emoji or name>       animate one glyph
//	g1yph random                a random glyph
//	g1yph all                   cycle through every glyph
//	g1yph list                  list glyphs
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Cid-Emmerich/G1yph/internal/config"
	"github.com/Cid-Emmerich/G1yph/internal/glyph"
	"github.com/Cid-Emmerich/G1yph/internal/paint"
	"github.com/Cid-Emmerich/G1yph/internal/theme"
	"github.com/Cid-Emmerich/G1yph/internal/ui"
)

const version = "0.1.0"

func usage() {
	fmt.Print(`g1yph ` + version + ` – emoji animator for the terminal

usage:
  g1yph                        open the picker
  g1yph <emoji or name>        animate a glyph:  g1yph ☕   g1yph coffee   g1yph "thumbs up"
  g1yph random                 start on a random glyph
  g1yph all                    cycle through every glyph (8 seconds each)
  g1yph list                   list every glyph with its aliases
  g1yph themes | styles | colours | charsets
  g1yph frame <glyph> [-w cols] [-h rows] [-at seconds] [-plain] [-png file]
                               print one frame to stdout (or save it as a PNG)
  g1yph sheet <file.png>       save a contact sheet of every glyph
  g1yph help

options (before the command):
  -t <theme>      colour theme (g1yph themes)
  -s <style>      render style: blocks, braille, ascii, chunky
  -c <colour>     colour mode: emoji, theme, rainbow, fire, ...
  -k <charset>    ascii charset: standard, detailed, blocks, ...
  -z <zoom>       zoom, 0.3 .. 1.6
  -f <fps>        frames per second (default 30)
  -l <label>      label: name, big, off
  -b <backdrop>   backdrop: none, stars, grid, dots
  -x              mirror
  -a              auto-cycle glyphs
  -o              draw a border

inside g1yph press ctrl+k for every keyboard shortcut.
config:  ~/.config/g1yph/g1yphrc
`)
}

func fail(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, "g1yph: "+msg+"\n", args...)
	os.Exit(1)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "g1yph: config: %v (using defaults)\n", err)
	}
	args := os.Args[1:]
	auto := false
	// options
	for len(args) > 0 && strings.HasPrefix(args[0], "-") && len(args[0]) > 1 {
		opt := args[0]
		args = args[1:]
		need := func() string {
			if len(args) == 0 {
				fail("%s needs a value", opt)
			}
			v := args[0]
			args = args[1:]
			return v
		}
		switch opt {
		case "-t", "--theme":
			cfg.Theme = need()
		case "-s", "--style":
			cfg.Style = need()
		case "-c", "--colour", "--color":
			cfg.Colour = need()
		case "-k", "--charset":
			cfg.Charset = need()
		case "-z", "--zoom":
			cfg.Zoom, _ = strconv.ParseFloat(need(), 64)
		case "-f", "--fps":
			cfg.FPS, _ = strconv.Atoi(need())
		case "-l", "--label":
			cfg.Label = need()
		case "-b", "--backdrop", "--bg":
			cfg.Background = need()
		case "-x", "--mirror":
			cfg.Mirror = true
		case "-a", "--auto":
			auto = true
		case "-o", "--border":
			cfg.Border = true
		case "-h", "--help":
			usage()
			return
		case "-v", "--version":
			fmt.Println("g1yph " + version)
			return
		default:
			fail("unknown option %s (try g1yph help)", opt)
		}
	}
	cfg.Auto = auto
	start := -1
	if len(args) > 0 {
		switch strings.ToLower(args[0]) {
		case "help", "-h", "--help":
			usage()
			return
		case "version":
			fmt.Println("g1yph " + version)
			return
		case "list", "ls", "glyphs", "emojis":
			list(strings.Join(args[1:], " "))
			return
		case "themes":
			fmt.Println(strings.Join(theme.Names(), "\n"))
			return
		case "styles":
			fmt.Println(strings.Join(paint.StyleNames, "\n"))
			return
		case "colours", "colors":
			fmt.Println(strings.Join(paint.ColourNames, "\n"))
			return
		case "charsets":
			fmt.Println(strings.Join(paint.CharsetNames, "\n"))
			return
		case "groups":
			fmt.Println(strings.Join(glyph.Groups(), "\n"))
			return
		case "frame", "render", "print":
			frame(cfg, args[1:])
			return
		case "sheet", "montage":
			sheet(cfg, args[1:])
			return
		case "random", "rand", "shuffle":
			start = -2
		case "all", "cycle", "tour", "demo":
			start = 0
			cfg.Auto = true
		default:
			q := strings.Join(args, " ")
			start = glyph.Find(q)
			if start < 0 {
				fmt.Fprintf(os.Stderr, "g1yph: no glyph for %q\n", q)
				if m := glyph.Search(strings.Fields(q)[0]); len(m) > 0 {
					fmt.Fprintln(os.Stderr, "did you mean:")
					for _, i := range m[:min(5, len(m))] {
						d := glyph.Registry[i]
						fmt.Fprintf(os.Stderr, "  %s  %s\n", d.Emoji, d.Name)
					}
				}
				fmt.Fprintln(os.Stderr, "run `g1yph list` to see every glyph")
				os.Exit(1)
			}
		}
	}
	var app *ui.App
	switch start {
	case -1:
		app = ui.New(cfg, -1)
	case -2:
		app = ui.New(cfg, -1)
		app.ClosePicker()
	default:
		app = ui.New(cfg, start)
	}
	if err := app.Run(); err != nil {
		fail("%v", err)
	}
}

func list(filter string) {
	idx := glyph.Search(filter)
	if filter == "" {
		idx = idx[:0]
		for i := range glyph.Registry {
			idx = append(idx, i)
		}
	}
	group := ""
	for _, i := range idx {
		d := glyph.Registry[i]
		if d.Group != group && filter == "" {
			group = d.Group
			fmt.Printf("\n%s\n", strings.ToUpper(group))
		}
		al := ""
		if len(d.Aliases) > 0 {
			al = "  (" + strings.Join(d.Aliases, ", ") + ")"
		}
		fmt.Printf("  %s  %-14s %s%s\n", d.Emoji, d.Name, d.Blurb, al)
	}
	if filter == "" {
		fmt.Printf("\n%d glyphs. Run: g1yph <name>\n", len(glyph.Registry))
	}
}

// frame prints a single frame of a glyph to stdout.
func frame(cfg *config.Config, args []string) {
	w, h, at, colour := 60, 24, 2.0, true
	pngPath := ""
	var words []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-w":
			i++
			if i < len(args) {
				w, _ = strconv.Atoi(args[i])
			}
		case "-h":
			i++
			if i < len(args) {
				h, _ = strconv.Atoi(args[i])
			}
		case "-at", "-t":
			i++
			if i < len(args) {
				at, _ = strconv.ParseFloat(args[i], 64)
			}
		case "-plain":
			colour = false
		case "-png":
			i++
			if i < len(args) {
				pngPath = args[i]
			}
		default:
			words = append(words, args[i])
		}
	}
	if len(words) == 0 {
		fail("frame needs a glyph, e.g. g1yph frame coffee")
	}
	idx := glyph.Find(strings.Join(words, " "))
	if idx < 0 {
		fail("no glyph for %q", strings.Join(words, " "))
	}
	if w < 4 || h < 2 {
		fail("frame size too small")
	}
	if pngPath != "" {
		f, err := os.Create(pngPath)
		if err != nil {
			fail("%v", err)
		}
		defer f.Close()
		if err := ui.FramePNG(f, cfg, idx, w, h, at); err != nil {
			fail("%v", err)
		}
		return
	}
	fmt.Print(ui.Frame(cfg, idx, w, h, at, colour))
}

// sheet writes a PNG contact sheet of every glyph.
func sheet(cfg *config.Config, args []string) {
	if len(args) == 0 {
		fail("sheet needs an output file, e.g. g1yph sheet glyphs.png")
	}
	at := 2.0
	cols, w, h := 10, 40, 20
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-at", "-t":
			i++
			if i < len(args) {
				at, _ = strconv.ParseFloat(args[i], 64)
			}
		case "-cols":
			i++
			if i < len(args) {
				cols, _ = strconv.Atoi(args[i])
			}
		case "-w":
			i++
			if i < len(args) {
				w, _ = strconv.Atoi(args[i])
			}
		case "-h":
			i++
			if i < len(args) {
				h, _ = strconv.Atoi(args[i])
			}
		}
	}
	f, err := os.Create(args[0])
	if err != nil {
		fail("%v", err)
	}
	defer f.Close()
	if err := ui.SheetPNG(f, cfg, max(cols, 1), max(w, 8), max(h, 4), at); err != nil {
		fail("%v", err)
	}
}
