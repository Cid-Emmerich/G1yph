# g1yph

An emoji animator for the terminal.

g1yph draws emoji as living, animated art in your terminal — a coffee cup with
curling steam, a waving flag, a blinking face. Every glyph is drawn from code in
unit coordinates rather than pasted as a character, so it stays crisp at any
terminal size, in any render style, in any colour.

```
              ██████████████████
             ███████████████████
             ███████  ███ ██
              ██████    █████
                ████████████
          ████████████████████████
          ████████████████████████
           ██████████████████████ ████
           ██████████████████████  █████
           ██████████████████████    ███
            ████████████████████     ███
            ████████████████████    ████
             ██████████████████   ████
             ██████████████████
             ██████████████████
      ████████████████████████████████
```

<sub>`g1yph frame coffee` — the same glyph animates in place when you run `g1yph coffee`.</sub>

## Features

- **100 hand-drawn animated glyphs** across nine groups: drinks, people, flags,
  math, buttons, nature, objects, animals and faces.
- **Four render styles** — `blocks` (half-block true colour), `braille`,
  `ascii` and `chunky` — with eight ASCII charsets to choose from.
- **Twelve colour modes**, from natural `emoji` colours to `rainbow`, `fire`,
  `ice`, `neon`, `matrix` and more.
- **Fourteen colour themes**: g1yph, wvfrm, nord, dracula, gruvbox, catppuccin,
  solarized, synthwave, sunset, forest, ocean, mono, amber, matrix.
- **Interactive picker** with fuzzy search by name, alias or pasted emoji, plus
  mouse support for the header tabs and scroll-to-change.
- **Auto-cycle mode** to tour every glyph hands-free.
- **Export to text or PNG** — a single frame, or a contact sheet of all 100.
- **Persistent config** at `~/.config/g1yph/g1yphrc`; settings are saved on quit.

## Installation

### With Go

Requires Go 1.25 or newer.

```sh
go install github.com/Cid-Emmerich/G1yph/cmd/g1yph@latest
```

The binary lands in `$(go env GOPATH)/bin`. Add it to your `PATH` if it isn't
already:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

### From source

```sh
git clone https://github.com/Cid-Emmerich/G1yph.git
cd G1yph
go build -o g1yph ./cmd/g1yph
./g1yph
```

To install the built binary system-wide:

```sh
sudo install -m 0755 g1yph /usr/local/bin/g1yph
```

### Requirements

- Go 1.25+ (to build).
- A terminal with 24-bit true colour for the best results — iTerm2, Kitty,
  WezTerm, Alacritty, Ghostty, Windows Terminal and modern GNOME Terminal all
  qualify. g1yph still runs on 256-colour terminals; try `-s ascii` or
  `-c mono` if colours look muddy.

## Usage

```sh
g1yph                      # open the picker
g1yph coffee               # animate one glyph by name
g1yph ☕                   # ...or by pasting the emoji
g1yph "thumbs up"          # ...or by alias
g1yph random               # start on a random glyph
g1yph all                  # cycle through every glyph
g1yph list                 # list every glyph with its aliases
g1yph help                 # full help
```

Options go **before** the command:

| Flag | Meaning |
| --- | --- |
| `-t <theme>` | colour theme (`g1yph themes`) |
| `-s <style>` | render style: `blocks`, `braille`, `ascii`, `chunky` |
| `-c <colour>` | colour mode: `emoji`, `theme`, `rainbow`, `fire`, … |
| `-k <charset>` | ASCII charset: `standard`, `detailed`, `blocks`, … |
| `-z <zoom>` | zoom, `0.3` – `1.6` |
| `-f <fps>` | frames per second (default `30`) |
| `-l <label>` | label: `name`, `big`, `off` |
| `-b <backdrop>` | backdrop: `none`, `stars`, `grid`, `dots` |
| `-x` | mirror horizontally |
| `-a` | auto-cycle glyphs |
| `-o` | draw a border |

For example:

```sh
g1yph -t nord -s ascii -c rainbow heart
```

Discover what's available with `g1yph themes`, `g1yph styles`,
`g1yph colours`, `g1yph charsets` and `g1yph groups`.

## Keyboard shortcuts

Press <kbd>ctrl+k</kbd> (or <kbd>?</kbd>) inside g1yph for the full list.

| Key | Action |
| --- | --- |
| <kbd>→</kbd> <kbd>n</kbd> / <kbd>←</kbd> <kbd>b</kbd> | next / previous glyph |
| <kbd>]</kbd> <kbd>[</kbd> | next / previous group |
| <kbd>r</kbd> | random glyph |
| <kbd>a</kbd> | auto-cycle on/off (<kbd>{</kbd> <kbd>}</kbd> change the interval) |
| <kbd>/</kbd> | search by name, alias or pasted emoji |
| <kbd>space</kbd> | pause / resume |
| <kbd>t</kbd> / <kbd>T</kbd> | next / previous theme |
| <kbd>s</kbd> / <kbd>S</kbd> | render style |
| <kbd>c</kbd> / <kbd>C</kbd> | ASCII charset |
| <kbd>g</kbd> / <kbd>G</kbd> | colour mode |
| <kbd>z</kbd> / <kbd>Z</kbd> | zoom in / out |
| <kbd>+</kbd> / <kbd>-</kbd> | faster / slower |
| <kbd>l</kbd> | label: name → big → off |
| <kbd>d</kbd> | backdrop: none → stars → grid → dots |
| <kbd>o</kbd> / <kbd>x</kbd> | border / mirror |
| <kbd>R</kbd> | reset look to defaults |
| <kbd>q</kbd> <kbd>ctrl+c</kbd> | quit (settings are saved) |

## Exporting

Print a single frame to stdout, or save it as a PNG:

```sh
g1yph frame coffee                          # 60×24 frame, in colour
g1yph frame coffee -w 80 -h 30 -plain       # sized, no ANSI colour
g1yph frame coffee -at 3.5 -png coffee.png  # the frame at t=3.5s, as a PNG
```

Save a contact sheet of every glyph:

```sh
g1yph sheet glyphs.png                      # defaults: 10 columns, 40×20 cells
g1yph sheet glyphs.png -cols 8 -w 48 -h 24 -at 2
```

Because `-plain` output is just text, it drops straight into a shell profile,
a MOTD or a `figlet`-style banner.

## Configuration

Settings live in `~/.config/g1yph/g1yphrc` as plain `key = value` lines. The
file is written when you quit, so the easiest way to configure g1yph is to set
things up interactively and press <kbd>q</kbd>. Editing by hand works too —
unknown or missing keys fall back to their defaults, so the file stays valid
across upgrades.

```ini
# g1yph configuration
auto = false
auto_seconds = 8.0
background = none
border = false
charset = standard
colour = emoji
fps = 30
label = name
mirror = false
speed = 1.00
style = blocks
theme = g1yph
zoom = 1.00
```

## Project layout

```
cmd/g1yph        CLI entry point and argument parsing
internal/glyph   the glyph library — one file per group
internal/paint   the rasteriser: painter, colour, gradients, fonts, particles
internal/ui      tcell application, drawing, picker, help, frame/PNG export
internal/theme   colour themes
internal/config  ~/.config/g1yph/g1yphrc
```

Each glyph is a `glyph.Def` with an `Anim` that draws into a `paint.Painter`
using unit coordinates, so adding one means writing a draw function and
registering it in the relevant group file.

## Development

```sh
go build ./...
go test ./...
go vet ./...
```

## Acknowledgements

Built on [tcell](https://github.com/gdamore/tcell) and
[go-runewidth](https://github.com/mattn/go-runewidth).

## License

Released under the [MIT License](LICENSE). Copyright (c) 2026 Cid Emmerich.
