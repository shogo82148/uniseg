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
	wbZWJBit = 16 // This bit is set for any states followed by at least one zero-width joiner (see WB4 and WB3c).
)

type wbTransitionResult struct {
	WordBreakState
	boundary   bool
	ruleNumber int
}

// The word break parser's state transitions. It's analogous to grTransitions,
// see comments there for details. Unicode version 15.0.0.
func wbTransitions(wb WordBreakState, p wbProperty) (wbTransitionResult, bool) {
	switch wb {
	case wbAny:
		switch p {
		case wbprNewline: // WB3b.
			return wbTransitionResult{wbNewline, true, 32}, true
		case wbprCR: // WB3b.
			return wbTransitionResult{wbCR, true, 32}, true
		case wbprLF: // WB3b.
			return wbTransitionResult{wbLF, true, 32}, true
		case wbprWSegSpace: // WB3d.
			return wbTransitionResult{wbWSegSpace, true, 9990}, true
		case wbprALetter: // WB5.
			return wbTransitionResult{wbALetter, true, 9990}, true
		case wbprHebrewLetter: // WB5.
			return wbTransitionResult{wbHebrewLetter, true, 9990}, true
		case wbprNumeric: // WB8.
			return wbTransitionResult{wbNumeric, true, 9990}, true
		case wbprKatakana: // WB13.
			return wbTransitionResult{wbKatakana, true, 9990}, true
		case wbprExtendNumLet: // WB13a.
			return wbTransitionResult{wbExtendNumLet, true, 9990}, true
		}

	case wbCR:
		switch p {
		case wbprAny: // WB3a.
			return wbTransitionResult{wbAny, true, 31}, true
		case wbprLF: // WB3.
			return wbTransitionResult{wbLF, false, 30}, true
		}
	case wbLF:
		switch p {
		case wbprAny: // WB3a.
			return wbTransitionResult{wbAny, true, 31}, true
		}
	case wbNewline:
		switch p {
		case wbprAny: // WB3a.
			return wbTransitionResult{wbAny, true, 31}, true
		}
	case wbWSegSpace:
		switch p {
		case wbprWSegSpace: // WB3d.
			return wbTransitionResult{wbWSegSpace, false, 34}, true
		}
	case wbHebrewLetter:
		switch p {
		case wbprALetter: // WB5.
			return wbTransitionResult{wbALetter, false, 50}, true
		case wbprHebrewLetter: // WB5.
			return wbTransitionResult{wbHebrewLetter, false, 50}, true
		case wbprSingleQuote: // WB7a.
			return wbTransitionResult{wbAny, false, 71}, true
		case wbprNumeric: // WB9.
			return wbTransitionResult{wbNumeric, false, 90}, true
		case wbprExtendNumLet: // WB13a.
			return wbTransitionResult{wbExtendNumLet, false, 131}, true
		}
	case wbALetter:
		switch p {
		case wbprALetter: // WB5.
			return wbTransitionResult{wbALetter, false, 50}, true
		case wbprHebrewLetter: // WB5.
			return wbTransitionResult{wbHebrewLetter, false, 50}, true
		case wbprNumeric: // WB9.
			return wbTransitionResult{wbNumeric, false, 90}, true
		case wbprExtendNumLet: // WB13a.
			return wbTransitionResult{wbExtendNumLet, false, 131}, true
		}
	case wbWB7:
		switch p {
		case wbprALetter: // WB7.
			return wbTransitionResult{wbALetter, false, 70}, true
		case wbprHebrewLetter: // WB7.
			return wbTransitionResult{wbHebrewLetter, false, 70}, true
		}
	case wbWB7c:
		switch p {
		case wbprHebrewLetter: // WB7c.
			return wbTransitionResult{wbHebrewLetter, false, 73}, true
		}
	case wbNumeric:
		switch p {
		case wbprNumeric: // WB8.
			return wbTransitionResult{wbNumeric, false, 80}, true
		case wbprALetter: // WB10.
			return wbTransitionResult{wbALetter, false, 100}, true
		case wbprHebrewLetter: // WB10.
			return wbTransitionResult{wbHebrewLetter, false, 100}, true
		case wbprExtendNumLet: // WB13a.
			return wbTransitionResult{wbExtendNumLet, false, 131}, true
		}
	case wbWB11:
		switch p {
		case wbprNumeric: // WB11.
			return wbTransitionResult{wbNumeric, false, 110}, true
		}
	case wbKatakana:
		switch p {
		case wbprKatakana: // WB13.
			return wbTransitionResult{wbKatakana, false, 130}, true
		case wbprExtendNumLet: // WB13a.
			return wbTransitionResult{wbExtendNumLet, false, 131}, true
		}
	case wbExtendNumLet:
		switch p {
		case wbprExtendNumLet: // WB13a.
			return wbTransitionResult{wbExtendNumLet, false, 131}, true
		case wbprALetter: // WB13b.
			return wbTransitionResult{wbALetter, false, 132}, true
		case wbprHebrewLetter: // WB13b.
			return wbTransitionResult{wbHebrewLetter, false, 132}, true
		case wbprNumeric: // WB13b.
			return wbTransitionResult{wbNumeric, false, 132}, true
		case wbprKatakana: // WB13b.
			return wbTransitionResult{wbKatakana, false, 132}, true
		}
	}

	return wbTransitionResult{}, false
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
	case wbprExtendedPictographic:
		if state >= 0 && state&wbZWJBit != 0 {
			// WB3c.
			return wbAny, false
		}
	}
	if state > 0 {
		state = state &^ wbZWJBit
	}

	// Find the applicable transition in the table.
	var rule int
	transition, ok := wbTransitions(state, nextProperty)
	if ok {
		// We have a specific transition. We'll use it.
		newState, wordBreak, rule = transition.WordBreakState, transition.boundary, transition.ruleNumber
	} else {
		// No specific transition found. Try the less specific ones.
		transAnyProp, okAnyProp := wbTransitions(state, wbprAny)
		transAnyState, okAnyState := wbTransitions(wbAny, nextProperty)
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
			prop := workBreakCodePoints.search(r)
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
