package uniseg

import "unicode/utf8"

// WordBreakState is the type of the word break parser's states.
type WordBreakState int

// The states of the word break parser.
const (
	_ WordBreakState = iota // The zero value is reserved for the initial state.
	wbAny
	wbCR
	wbLF
	wbNewline
	wbWSegSpace
	wbHebrewLetter
	wbALetter
	wbWB7
	wbWB7c
	wbNumeric
	wbWB11
	wbKatakana
	wbExtendNumLet
	wbOddRI
	wbEvenRI
	wbMax = iota

	// This bit is set for any states followed by at least one zero-width joiner (see WB4 and WB3c).
	wbZWJBit WordBreakState = 16
)

type wbTransitionResult = transitionResult[WordBreakState, bool]

// The word break parser's state transitions. It's analogous to grTransitions,
// see comments there for details. Unicode version 16.0.0.
var wbTransitions = [wbMax * wbprMax]wbTransitionResult{
	// WB3b.
	int(wbAny)*wbprMax + int(wbprNewline): {wbNewline, true, 32},
	int(wbAny)*wbprMax + int(wbprCR):      {wbCR, true, 32},
	int(wbAny)*wbprMax + int(wbprLF):      {wbLF, true, 32},

	// WB3a.
	int(wbNewline)*wbprMax + int(wbprAny): {wbAny, true, 31},
	int(wbCR)*wbprMax + int(wbprAny):      {wbAny, true, 31},
	int(wbLF)*wbprMax + int(wbprAny):      {wbAny, true, 31},

	// WB3.
	int(wbCR)*wbprMax + int(wbprLF): {wbLF, false, 30},

	// WB3d.
	int(wbAny)*wbprMax + int(wbprWSegSpace):       {wbWSegSpace, true, 9990},
	int(wbWSegSpace)*wbprMax + int(wbprWSegSpace): {wbWSegSpace, false, 34},

	// WB5.
	int(wbAny)*wbprMax + int(wbprALetter):               {wbALetter, true, 9990},
	int(wbAny)*wbprMax + int(wbprHebrewLetter):          {wbHebrewLetter, true, 9990},
	int(wbALetter)*wbprMax + int(wbprALetter):           {wbALetter, false, 50},
	int(wbALetter)*wbprMax + int(wbprHebrewLetter):      {wbHebrewLetter, false, 50},
	int(wbHebrewLetter)*wbprMax + int(wbprALetter):      {wbALetter, false, 50},
	int(wbHebrewLetter)*wbprMax + int(wbprHebrewLetter): {wbHebrewLetter, false, 50},

	// WB7. Transitions to wbWB7 handled by transitionWordBreakState().
	int(wbWB7)*wbprMax + int(wbprALetter):      {wbALetter, false, 70},
	int(wbWB7)*wbprMax + int(wbprHebrewLetter): {wbHebrewLetter, false, 70},

	// WB7a.
	int(wbHebrewLetter)*wbprMax + int(wbprSingleQuote): {wbAny, false, 71},

	// WB7c. Transitions to wbWB7c handled by transitionWordBreakState().
	int(wbWB7c)*wbprMax + int(wbprHebrewLetter): {wbHebrewLetter, false, 73},

	// WB8.
	int(wbAny)*wbprMax + int(wbprNumeric):     {wbNumeric, true, 9990},
	int(wbNumeric)*wbprMax + int(wbprNumeric): {wbNumeric, false, 80},

	// WB9.
	int(wbALetter)*wbprMax + int(wbprNumeric):      {wbNumeric, false, 90},
	int(wbHebrewLetter)*wbprMax + int(wbprNumeric): {wbNumeric, false, 90},

	// WB10.
	int(wbNumeric)*wbprMax + int(wbprALetter):      {wbALetter, false, 100},
	int(wbNumeric)*wbprMax + int(wbprHebrewLetter): {wbHebrewLetter, false, 100},

	// WB11. Transitions to wbWB11 handled by transitionWordBreakState().
	int(wbWB11)*wbprMax + int(wbprNumeric): {wbNumeric, false, 110},

	// WB13.
	int(wbAny)*wbprMax + int(wbprKatakana):      {wbKatakana, true, 9990},
	int(wbKatakana)*wbprMax + int(wbprKatakana): {wbKatakana, false, 130},

	// WB13a.
	int(wbAny)*wbprMax + int(wbprExtendNumLet):          {wbExtendNumLet, true, 9990},
	int(wbALetter)*wbprMax + int(wbprExtendNumLet):      {wbExtendNumLet, false, 131},
	int(wbHebrewLetter)*wbprMax + int(wbprExtendNumLet): {wbExtendNumLet, false, 131},
	int(wbNumeric)*wbprMax + int(wbprExtendNumLet):      {wbExtendNumLet, false, 131},
	int(wbKatakana)*wbprMax + int(wbprExtendNumLet):     {wbExtendNumLet, false, 131},
	int(wbExtendNumLet)*wbprMax + int(wbprExtendNumLet): {wbExtendNumLet, false, 131},

	// WB13b.
	int(wbExtendNumLet)*wbprMax + int(wbprALetter):      {wbALetter, false, 132},
	int(wbExtendNumLet)*wbprMax + int(wbprHebrewLetter): {wbHebrewLetter, false, 132},
	int(wbExtendNumLet)*wbprMax + int(wbprNumeric):      {wbNumeric, false, 132},
	int(wbExtendNumLet)*wbprMax + int(wbprKatakana):     {wbKatakana, false, 132},
}

// transitionWordBreakState determines the new state of the word break parser
// given the current state and the next code point. It also returns whether a
// word boundary was detected. If more than one code point is needed to
// determine the new state, the byte slice or the string starting after rune "r"
// can be used (whichever is not nil or empty) for further lookups.
func transitionWordBreakState[T bytes](state WordBreakState, r rune, str T, decoder runeDecoder[T]) (newState WordBreakState, wordBreak bool) {
	// Determine the property of the next character.
	nextProperty := wordBreakLookup.search(r)

	// "Replacing Ignore Rules".
	switch nextProperty {
	case wbprZWJ:
		// WB4 (for zero-width joiners).
		if state == wbNewline || state == wbCR || state == wbLF {
			return wbAny | wbZWJBit, true // Make sure we don't apply WB4 to WB3a.
		}
		if state <= 0 {
			return wbAny | wbZWJBit, false
		}
		return state | wbZWJBit, false
	case wbprExtend, wbprFormat:
		// WB4 (for Extend and Format).
		if state == wbNewline || state == wbCR || state == wbLF {
			return wbAny, true // Make sure we don't apply WB4 to WB3a.
		}
		if state == wbWSegSpace || state == wbAny|wbZWJBit {
			return wbAny, false // We don't break but this is also not WB3d or WB3c.
		}
		if state <= 0 {
			return wbAny, false
		}
		return state, false
	}
	if state >= 0 && state&wbZWJBit != 0 && graphemeLookup.search(r) == prExtendedPictographic {
		// WB3c.
		return wbAny, false
	}
	if state > 0 {
		state = state &^ wbZWJBit
	}

	// Find the applicable transition in the table.
	transition := wbTransitions[int(state)*wbprMax+int(nextProperty)]
	newState, wordBreak, rule := transition.state, transition.boundary, transition.ruleNumber

	// For those rules that need to look up runes further in the string, we
	// determine the property after nextProperty, skipping over Format, Extend,
	// and ZWJ (according to WB4). It's -1 if not needed, if such a rune cannot
	// be determined (because the text ends or the rune is faulty).
	farProperty := wbProperty(-1)
	if rule > 60 &&
		(state == wbALetter || state == wbHebrewLetter || state == wbNumeric) &&
		(nextProperty == wbprMidLetter || nextProperty == wbprMidNumLet || nextProperty == wbprSingleQuote || // WB6.
			nextProperty == wbprDoubleQuote || // WB7b.
			nextProperty == wbprMidNum) { // WB12.
		for {
			r, length := decoder(str)
			str = str[length:]
			if r == utf8.RuneError {
				break
			}
			prop := wordBreakLookup.search(r)
			if prop == wbprExtend || prop == wbprFormat || prop == wbprZWJ {
				continue
			}
			farProperty = prop
			break
		}
	}

	// WB6.
	if rule > 60 &&
		(state == wbALetter || state == wbHebrewLetter) &&
		(nextProperty == wbprMidLetter || nextProperty == wbprMidNumLet || nextProperty == wbprSingleQuote) &&
		(farProperty == wbprALetter || farProperty == wbprHebrewLetter) {
		return wbWB7, false
	}

	// WB7b.
	if rule > 72 &&
		state == wbHebrewLetter &&
		nextProperty == wbprDoubleQuote &&
		farProperty == wbprHebrewLetter {
		return wbWB7c, false
	}

	// WB12.
	if rule > 120 &&
		state == wbNumeric &&
		(nextProperty == wbprMidNum || nextProperty == wbprMidNumLet || nextProperty == wbprSingleQuote) &&
		farProperty == wbprNumeric {
		return wbWB11, false
	}

	// WB15 and WB16.
	if newState == wbAny && nextProperty == wbprRegionalIndicator {
		if state != wbOddRI && state != wbEvenRI { // Includes state == 0.
			// Transition into the first RI.
			return wbOddRI, true
		}
		if state == wbOddRI {
			// Don't break pairs of Regional Indicators.
			return wbEvenRI, false
		}
		return wbOddRI, true // We can break after a pair.
	}

	return
}

func init() {
	// WB999: Any ÷ Any.
	resolveTransitions(wbTransitions[:], wbprMax, int(wbAny), int(wbprAny), wbTransitionResult{wbAny, true, 9990})
}
