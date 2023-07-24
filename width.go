package uniseg

// runeWidth returns the monospace width for the given rune. The provided
// grapheme property is a value mapped by the [graphemeCodePoints] table.
// runeWidth calculates the width of a given rune based on its grapheme property and the current parser settings.
func runeWidth(p *Parser, r rune, graphemeProperty property) int {
	// Check the grapheme property of the rune.
	switch graphemeProperty {
	case prControl, prCR, prLF, prExtend, prZWJ:
		// If the property is a control character, carriage return, line feed, extend character, or zero-width joiner, return a width of 0.
		return 0
	case prRegionalIndicator:
		// If the property is a regional indicator, return a width of 2.
		return 2
	case prExtendedPictographic:
		if emojiPresentation.search(r) == prEmojiPresentation {
			// If only the WideEmoji setting is true and the rune has an emoji presentation property, return a width of 2.
			return 2
		}
		if p.EastAsianWidth {
			return p.runeWidthAE(r, graphemeProperty)
		}
		// Otherwise, return a width of 1.
		return 1
	}

	// Check for specific runes that have a fixed width.
	switch r {
	case '\u2e3a': // TWO-EM DASH: Width of 3
		return 3
	case '\u2e3b': // THREE-EM DASH: Width of 4
		return 4
	}

	return p.runeWidthAE(r, graphemeProperty)
}

func (p *Parser) runeWidthAE(r rune, graphemeProperty property) int {
	if p.EastAsianWidth && p.WideEmoji {
		if graphemeProperty == prExtendedPictographic {
			return 2
		}
		if emoji.search(r) == prEmoji {
			return 2
		}
	}

	// Check the East Asian Width property of the rune.
	switch eastAsianWidth.search(r) {
	case eawprW, eawprF:
		// If the property is Wide or Fullwidth, return a width of 2.
		return 2
	case eawprA:
		if p.EastAsianWidth {
			// If the EastAsianWidth setting is true, return a width of 2.
			return 2
		} else {
			// Otherwise, return a width of 1.
			return 1
		}
	case eawprNa, eawprH, eawprN:
		// If the property is Neutral, Halfwidth, or Narrow, return a width of 1.
		return 1
	}

	// If the property is not recognized, return a width of 1.
	return 1
}

// StringWidth returns the monospace width for the given string, that is, the
// number of same-size cells to be occupied by the string.
func StringWidth(s string) (width int) {
	var state GraphemeBreakState
	for len(s) > 0 {
		var w int
		_, s, w, state = DefaultParser.FirstGraphemeClusterInString(s, state)
		width += w
	}
	return
}

// StringWidth returns the monospace width for the given string, that is, the
// number of same-size cells to be occupied by the string.
func (p *Parser) StringWidth(s string) (width int) {
	var state GraphemeBreakState
	for len(s) > 0 {
		var w int
		_, s, w, state = p.FirstGraphemeClusterInString(s, state)
		width += w
	}
	return
}
