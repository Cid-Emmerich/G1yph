// Package config reads and writes ~/.config/g1yph/g1yphrc, a simple
// key = value file. Missing keys keep their defaults so old files keep
// working after upgrades.
package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Config holds every user-tunable setting.
type Config struct {
	Theme       string  // colour theme name
	Style       string  // blocks, braille, ascii, chunky
	Charset     string  // ascii charset name
	Colour      string  // emoji, theme, rainbow, ...
	Zoom        float64 // 0.3 .. 1.5
	Speed       float64 // animation speed multiplier
	Label       string  // off, name, big
	Background  string  // none, stars, grid, dots
	Border      bool
	Mirror      bool
	Auto        bool    // cycle glyphs automatically
	AutoSeconds float64 // seconds per glyph when cycling
	FPS         int
	Last        string // last glyph shown

	Path string
}

// Default returns the baseline configuration.
func Default() Config {
	home, _ := os.UserHomeDir()
	return Config{
		Theme: "g1yph", Style: "blocks", Charset: "standard", Colour: "emoji",
		Zoom: 1, Speed: 1, Label: "name", Background: "none", Border: false,
		Mirror: false, Auto: false, AutoSeconds: 8, FPS: 30, Last: "",
		Path: filepath.Join(home, ".config", "g1yph", "g1yphrc"),
	}
}

// Load reads the config file over the defaults. A missing file is fine.
func Load() (*Config, error) {
	c := Default()
	f, err := os.Open(c.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return &c, nil
		}
		return &c, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		c.set(strings.TrimSpace(k), strings.TrimSpace(v))
	}
	return &c, sc.Err()
}

func (c *Config) set(k, v string) {
	b := func() bool { return v == "true" || v == "1" || v == "yes" || v == "on" }
	fl := func(def float64) float64 {
		if x, err := strconv.ParseFloat(v, 64); err == nil {
			return x
		}
		return def
	}
	switch k {
	case "theme":
		c.Theme = v
	case "style":
		c.Style = v
	case "charset":
		c.Charset = v
	case "colour", "color":
		c.Colour = v
	case "zoom":
		c.Zoom = fl(c.Zoom)
	case "speed":
		c.Speed = fl(c.Speed)
	case "label":
		c.Label = v
	case "background", "bg":
		c.Background = v
	case "border":
		c.Border = b()
	case "mirror":
		c.Mirror = b()
	case "auto":
		c.Auto = b()
	case "auto_seconds":
		c.AutoSeconds = fl(c.AutoSeconds)
	case "fps":
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			c.FPS = n
		}
	case "last":
		c.Last = v
	}
}

// Save writes the config file, creating the directory if needed.
func (c *Config) Save() error {
	if err := os.MkdirAll(filepath.Dir(c.Path), 0o755); err != nil {
		return err
	}
	kv := map[string]string{
		"theme": c.Theme, "style": c.Style, "charset": c.Charset, "colour": c.Colour,
		"zoom": fmt.Sprintf("%.2f", c.Zoom), "speed": fmt.Sprintf("%.2f", c.Speed),
		"label": c.Label, "background": c.Background,
		"border": strconv.FormatBool(c.Border), "mirror": strconv.FormatBool(c.Mirror),
		"auto": strconv.FormatBool(c.Auto), "auto_seconds": fmt.Sprintf("%.1f", c.AutoSeconds),
		"fps": strconv.Itoa(c.FPS), "last": c.Last,
	}
	keys := make([]string, 0, len(kv))
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	sb.WriteString("# g1yph configuration\n# Edit by hand, or change settings inside g1yph (saved on quit).\n\n")
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s = %s\n", k, kv[k])
	}
	return os.WriteFile(c.Path, []byte(sb.String()), 0o644)
}
