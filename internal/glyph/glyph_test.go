package glyph

import (
	"testing"

	"github.com/Cid-Emmerich/G1yph/internal/paint"
)

// TestEveryGlyphDraws runs each animation for a few seconds at a few sizes
// and checks it paints something without panicking.
func TestEveryGlyphDraws(t *testing.T) {
	sizes := [][2]int{{40, 40}, {160, 80}, {12, 8}}
	for _, d := range Registry {
		for _, sz := range sizes {
			anim := d.New()
			b := paint.NewBitmap(sz[0], sz[1])
			painted := false
			for i := 0; i < 150; i++ {
				b.Clear()
				p := paint.NewPainter(b, 1, i%2 == 1)
				anim.Draw(p, float64(i)/30, 1.0/30)
				for _, px := range b.P {
					if px.A > 0 {
						painted = true
						break
					}
				}
			}
			if !painted && sz[0] >= 40 {
				t.Errorf("%s (%s) never painted anything at %dx%d", d.Name, d.Emoji, sz[0], sz[1])
			}
		}
	}
}

func TestFind(t *testing.T) {
	cases := map[string]string{"☕": "coffee", "coffee": "coffee", "❤️": "heart", "❤": "heart", "thumbs": "thumbs up", "Heart Eyes": "heart eyes", "🏴‍☠️": "pirate", "us": "usa", "➕": "plus", "✖": "plus"}
	for q, want := range cases {
		i := Find(q)
		if i < 0 || Registry[i].Name != want {
			got := "nothing"
			if i >= 0 {
				got = Registry[i].Name
			}
			t.Errorf("Find(%q) = %s, want %s", q, got, want)
		}
	}
	if Find("zzzzqqq") != -1 {
		t.Error("Find of nonsense should fail")
	}
}

func TestNoDuplicateNames(t *testing.T) {
	seen := map[string]bool{}
	for _, d := range Registry {
		if seen[d.Name] {
			t.Errorf("duplicate glyph name %q", d.Name)
		}
		seen[d.Name] = true
		if seen[d.Emoji] {
			t.Errorf("duplicate glyph emoji %q", d.Emoji)
		}
		seen[d.Emoji] = true
	}
}
