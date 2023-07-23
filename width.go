package uniseg

// runeWidth returns the monospace width for the given rune. The provided
// grapheme property is a value mapped by the [graphemeCodePoints] table.
//
// Every rune has a width of 1, except for runes with the following properties
// (evaluated in this order):
//
//   - Control, CR, LF, Extend, ZWJ: Width of 0
//   - \u2e3a, TWO-EM DASH: Width of 3
//   - \u2e3b, THREE-EM DASH: Width of 4
//   - East-Asian width Fullwidth and Wide: Width of 2 (Ambiguous and Neutral
//     have a width of 1)
//   - Regional Indicator: Width of 2
//   - Extended Pictographic: Width of 2, unless Emoji Presentation is "No".
func runeWidth(r rune, graphemeProperty property) int {
	switch graphemeProperty {
	case prControl, prCR, prLF, prExtend, prZWJ:
		return 0
	case prRegionalIndicator:
		return 2
	case prExtendedPictographic:
		if emojiPresentation.search(r) == prEmojiPresentation {
			return 2
		}
		return 1
	}

	switch r {
	case '\u2e3a': // TWO-EM DASH: Width of 3
		return 3
	case '\u2e3b': // THREE-EM DASH: Width of 4
		return 4
	}

	switch eastAsianWidth.search(r) {
	case prW, prF:
		return 2
	}

	return 1
}

// StringWidth returns the monospace width for the given string, that is, the
// number of same-size cells to be occupied by the string.
func StringWidth(s string) (width int) {
	var state GraphemeBreakState
	for len(s) > 0 {
		var w int
		_, s, w, state = FirstGraphemeClusterInString(s, state)
		width += w
	}
	return
}
