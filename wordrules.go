package uniseg

import "unicode/utf8"

// WordBreakState is the type of the word break parser's states.
type WordBreakState int

// The states of the word break parser.
const (
	wbAny WordBreakState = iota
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
	wbZWJBit = 16 // This bit is set for any states followed by at least one zero-width joiner (see WB4 and WB3c).
)

type wbStateProperty struct {
	WordBreakState
	property
}

type wbTransitionResult struct {
	WordBreakState
	boundary   bool
	ruleNumber int
}

// The word break parser's state transitions. It's analogous to grTransitions,
// see comments there for details. Unicode version 15.0.0.
var wbTransitions = map[wbStateProperty]wbTransitionResult{
	// WB3b.
	{wbAny, prNewline}: {wbNewline, true, 32},
	{wbAny, prCR}:      {wbCR, true, 32},
	{wbAny, prLF}:      {wbLF, true, 32},

	// WB3a.
	{wbNewline, prAny}: {wbAny, true, 31},
	{wbCR, prAny}:      {wbAny, true, 31},
	{wbLF, prAny}:      {wbAny, true, 31},

	// WB3.
	{wbCR, prLF}: {wbLF, false, 30},

	// WB3d.
	{wbAny, prWSegSpace}:       {wbWSegSpace, true, 9990},
	{wbWSegSpace, prWSegSpace}: {wbWSegSpace, false, 34},

	// WB5.
	{wbAny, prALetter}:               {wbALetter, true, 9990},
	{wbAny, prHebrewLetter}:          {wbHebrewLetter, true, 9990},
	{wbALetter, prALetter}:           {wbALetter, false, 50},
	{wbALetter, prHebrewLetter}:      {wbHebrewLetter, false, 50},
	{wbHebrewLetter, prALetter}:      {wbALetter, false, 50},
	{wbHebrewLetter, prHebrewLetter}: {wbHebrewLetter, false, 50},

	// WB7. Transitions to wbWB7 handled by transitionWordBreakState().
	{wbWB7, prALetter}:      {wbALetter, false, 70},
	{wbWB7, prHebrewLetter}: {wbHebrewLetter, false, 70},

	// WB7a.
	{wbHebrewLetter, prSingleQuote}: {wbAny, false, 71},

	// WB7c. Transitions to wbWB7c handled by transitionWordBreakState().
	{wbWB7c, prHebrewLetter}: {wbHebrewLetter, false, 73},

	// WB8.
	{wbAny, prNumeric}:     {wbNumeric, true, 9990},
	{wbNumeric, prNumeric}: {wbNumeric, false, 80},

	// WB9.
	{wbALetter, prNumeric}:      {wbNumeric, false, 90},
	{wbHebrewLetter, prNumeric}: {wbNumeric, false, 90},

	// WB10.
	{wbNumeric, prALetter}:      {wbALetter, false, 100},
	{wbNumeric, prHebrewLetter}: {wbHebrewLetter, false, 100},

	// WB11. Transitions to wbWB11 handled by transitionWordBreakState().
	{wbWB11, prNumeric}: {wbNumeric, false, 110},

	// WB13.
	{wbAny, prKatakana}:      {wbKatakana, true, 9990},
	{wbKatakana, prKatakana}: {wbKatakana, false, 130},

	// WB13a.
	{wbAny, prExtendNumLet}:          {wbExtendNumLet, true, 9990},
	{wbALetter, prExtendNumLet}:      {wbExtendNumLet, false, 131},
	{wbHebrewLetter, prExtendNumLet}: {wbExtendNumLet, false, 131},
	{wbNumeric, prExtendNumLet}:      {wbExtendNumLet, false, 131},
	{wbKatakana, prExtendNumLet}:     {wbExtendNumLet, false, 131},
	{wbExtendNumLet, prExtendNumLet}: {wbExtendNumLet, false, 131},

	// WB13b.
	{wbExtendNumLet, prALetter}:      {wbALetter, false, 132},
	{wbExtendNumLet, prHebrewLetter}: {wbHebrewLetter, false, 132},
	{wbExtendNumLet, prNumeric}:      {wbNumeric, false, 132},
	{wbExtendNumLet, prKatakana}:     {wbKatakana, false, 132},
}

// transitionWordBreakState determines the new state of the word break parser
// given the current state and the next code point. It also returns whether a
// word boundary was detected. If more than one code point is needed to
// determine the new state, the byte slice or the string starting after rune "r"
// can be used (whichever is not nil or empty) for further lookups.
func transitionWordBreakState[T bytes](state WordBreakState, r rune, str T, decoder runeDecoder[T]) (newState WordBreakState, wordBreak bool) {
	// Determine the property of the next character.
	nextProperty := workBreakCodePoints.search(r)

	// "Replacing Ignore Rules".
	if nextProperty == prZWJ {
		// WB4 (for zero-width joiners).
		if state == wbNewline || state == wbCR || state == wbLF {
			return wbAny | wbZWJBit, true // Make sure we don't apply WB4 to WB3a.
		}
		if state < 0 {
			return wbAny | wbZWJBit, false
		}
		return state | wbZWJBit, false
	} else if nextProperty == prExtend || nextProperty == prFormat {
		// WB4 (for Extend and Format).
		if state == wbNewline || state == wbCR || state == wbLF {
			return wbAny, true // Make sure we don't apply WB4 to WB3a.
		}
		if state == wbWSegSpace || state == wbAny|wbZWJBit {
			return wbAny, false // We don't break but this is also not WB3d or WB3c.
		}
		if state < 0 {
			return wbAny, false
		}
		return state, false
	} else if nextProperty == prExtendedPictographic && state >= 0 && state&wbZWJBit != 0 {
		// WB3c.
		return wbAny, false
	}
	if state >= 0 {
		state = state &^ wbZWJBit
	}

	// Find the applicable transition in the table.
	var rule int
	transition, ok := wbTransitions[wbStateProperty{state, nextProperty}]
	if ok {
		// We have a specific transition. We'll use it.
		newState, wordBreak, rule = transition.WordBreakState, transition.boundary, transition.ruleNumber
	} else {
		// No specific transition found. Try the less specific ones.
		transAnyProp, okAnyProp := wbTransitions[wbStateProperty{state, prAny}]
		transAnyState, okAnyState := wbTransitions[wbStateProperty{wbAny, nextProperty}]
		if okAnyProp && okAnyState {
			// Both apply. We'll use a mix (see comments for grTransitions).
			newState, wordBreak, rule = transAnyState.WordBreakState, transAnyState.boundary, transAnyState.ruleNumber
			if transAnyProp.ruleNumber < transAnyState.ruleNumber {
				wordBreak, rule = transAnyProp.boundary, transAnyProp.ruleNumber
			}
		} else if okAnyProp {
			// We only have a specific state.
			newState, wordBreak, rule = transAnyProp.WordBreakState, transAnyProp.boundary, transAnyProp.ruleNumber
			// This branch will probably never be reached because okAnyState will
			// always be true given the current transition map. But we keep it here
			// for future modifications to the transition map where this may not be
			// true anymore.
		} else if okAnyState {
			// We only have a specific property.
			newState, wordBreak, rule = transAnyState.WordBreakState, transAnyState.boundary, transAnyState.ruleNumber
		} else {
			// No known transition. WB999: Any ÷ Any.
			newState, wordBreak, rule = wbAny, true, 9990
		}
	}

	// For those rules that need to look up runes further in the string, we
	// determine the property after nextProperty, skipping over Format, Extend,
	// and ZWJ (according to WB4). It's -1 if not needed, if such a rune cannot
	// be determined (because the text ends or the rune is faulty).
	farProperty := property(-1)
	if rule > 60 &&
		(state == wbALetter || state == wbHebrewLetter || state == wbNumeric) &&
		(nextProperty == prMidLetter || nextProperty == prMidNumLet || nextProperty == prSingleQuote || // WB6.
			nextProperty == prDoubleQuote || // WB7b.
			nextProperty == prMidNum) { // WB12.
		for {
			r, length := decoder(str)
			str = str[length:]
			if r == utf8.RuneError {
				break
			}
			prop := workBreakCodePoints.search(r)
			if prop == prExtend || prop == prFormat || prop == prZWJ {
				continue
			}
			farProperty = prop
			break
		}
	}

	// WB6.
	if rule > 60 &&
		(state == wbALetter || state == wbHebrewLetter) &&
		(nextProperty == prMidLetter || nextProperty == prMidNumLet || nextProperty == prSingleQuote) &&
		(farProperty == prALetter || farProperty == prHebrewLetter) {
		return wbWB7, false
	}

	// WB7b.
	if rule > 72 &&
		state == wbHebrewLetter &&
		nextProperty == prDoubleQuote &&
		farProperty == prHebrewLetter {
		return wbWB7c, false
	}

	// WB12.
	if rule > 120 &&
		state == wbNumeric &&
		(nextProperty == prMidNum || nextProperty == prMidNumLet || nextProperty == prSingleQuote) &&
		farProperty == prNumeric {
		return wbWB11, false
	}

	// WB15 and WB16.
	if newState == wbAny && nextProperty == prRegionalIndicator {
		if state != wbOddRI && state != wbEvenRI { // Includes state == -1.
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
