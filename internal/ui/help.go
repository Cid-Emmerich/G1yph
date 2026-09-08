package ui

// helpLines is the single source of truth for the ctrl+k help menu.
// Lines starting with "# " are section headers; others are "keys<TAB>description".
func helpLines() []string {
	return []string{
		"# Glyphs",
		"→  n  enter\tnext glyph",
		"←  b\tprevious glyph",
		"]  [\tnext / previous group (drinks, people, flags, math, buttons, nature, objects, animals, faces)",
		"r\ta random glyph",
		"a\tauto-cycle on / off (every 8 seconds; { } change the interval)",
		"/\tsearch: type a name, alias or paste an emoji, enter to pick",
		"space\tpause / resume the animation",
		"",
		"# Look",
		"t / T\tnext / previous colour theme",
		"s / S\trender style: blocks (half-block true colour) → braille → ascii → chunky",
		"c / C\tascii charset: standard, detailed, blocks, minimal, dots, lines, binary, hearts",
		"g / G\tcolour mode: emoji (natural), theme, horizontal, rainbow, spectrum, fire, ice, neon, heat, mono, pastel, matrix",
		"z / Z\tzoom in / out",
		"+ / -\tanimation faster / slower",
		"l\tlabel: name → big → off",
		"d\tbackdrop: none → stars → grid → dots",
		"o\tborder on / off",
		"x\tmirror (flip horizontally)",
		"R\treset look to defaults",
		"",
		"# General",
		"ctrl+k  ?\tthis help (↑↓ to scroll)",
		"esc\tclose help / search",
		"q  ctrl+c\tquit (your settings are saved)",
		"",
		"# Command line",
		"g1yph\topen with the picker",
		"g1yph <emoji or name>\tanimate one glyph, e.g. g1yph ☕  or  g1yph coffee",
		"g1yph random\tstart on a random glyph",
		"g1yph all\tcycle through every glyph",
		"g1yph list\tprint every glyph with its aliases",
		"g1yph -t nord -s ascii -g rainbow heart\tstart with a theme, style and colour mode",
		"",
		"mouse\tclick the header tabs, scroll to change glyph, click the picker",
	}
}
