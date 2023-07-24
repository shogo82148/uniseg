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

type wbTransitionResult struct {
	WordBreakState
	boundary   bool
	ruleNumber int
}

// The word break parser's state transitions. It's analogous to grTransitions,
// see comments there for details. Unicode version 15.0.0.
var wbTransitions = [wbMax * wbPropertyMax]wbTransitionResult{
	// WB3b.
	int(wbAny)*wbPropertyMax + int(wbprNewline): {wbNewline, true, 32},
	int(wbAny)*wbPropertyMax + int(wbprCR):      {wbCR, true, 32},
	int(wbAny)*wbPropertyMax + int(wbprLF):      {wbLF, true, 32},

	// WB3a.
	int(wbNewline)*wbPropertyMax + int(wbprAny): {wbAny, true, 31},
	int(wbCR)*wbPropertyMax + int(wbprAny):      {wbAny, true, 31},
	int(wbLF)*wbPropertyMax + int(wbprAny):      {wbAny, true, 31},

	// WB3.
	int(wbCR)*wbPropertyMax + int(wbprLF): {wbLF, false, 30},

	// WB3d.
	int(wbAny)*wbPropertyMax + int(wbprWSegSpace):       {wbWSegSpace, true, 9990},
	int(wbWSegSpace)*wbPropertyMax + int(wbprWSegSpace): {wbWSegSpace, false, 34},

	// WB5.
	int(wbAny)*wbPropertyMax + int(wbprALetter):               {wbALetter, true, 9990},
	int(wbAny)*wbPropertyMax + int(wbprHebrewLetter):          {wbHebrewLetter, true, 9990},
	int(wbALetter)*wbPropertyMax + int(wbprALetter):           {wbALetter, false, 50},
	int(wbALetter)*wbPropertyMax + int(wbprHebrewLetter):      {wbHebrewLetter, false, 50},
	int(wbHebrewLetter)*wbPropertyMax + int(wbprALetter):      {wbALetter, false, 50},
	int(wbHebrewLetter)*wbPropertyMax + int(wbprHebrewLetter): {wbHebrewLetter, false, 50},

	// WB7. Transitions to wbWB7 handled by transitionWordBreakState().
	int(wbWB7)*wbPropertyMax + int(wbprALetter):      {wbALetter, false, 70},
	int(wbWB7)*wbPropertyMax + int(wbprHebrewLetter): {wbHebrewLetter, false, 70},

	// WB7a.
	int(wbHebrewLetter)*wbPropertyMax + int(wbprSingleQuote): {wbAny, false, 71},

	// WB7c. Transitions to wbWB7c handled by transitionWordBreakState().
	int(wbWB7c)*wbPropertyMax + int(wbprHebrewLetter): {wbHebrewLetter, false, 73},

	// WB8.
	int(wbAny)*wbPropertyMax + int(wbprNumeric):     {wbNumeric, true, 9990},
	int(wbNumeric)*wbPropertyMax + int(wbprNumeric): {wbNumeric, false, 80},

	// WB9.
	int(wbALetter)*wbPropertyMax + int(wbprNumeric):      {wbNumeric, false, 90},
	int(wbHebrewLetter)*wbPropertyMax + int(wbprNumeric): {wbNumeric, false, 90},

	// WB10.
	int(wbNumeric)*wbPropertyMax + int(wbprALetter):      {wbALetter, false, 100},
	int(wbNumeric)*wbPropertyMax + int(wbprHebrewLetter): {wbHebrewLetter, false, 100},

	// WB11. Transitions to wbWB11 handled by transitionWordBreakState().
	int(wbWB11)*wbPropertyMax + int(wbprNumeric): {wbNumeric, false, 110},

	// WB13.
	int(wbAny)*wbPropertyMax + int(wbprKatakana):      {wbKatakana, true, 9990},
	int(wbKatakana)*wbPropertyMax + int(wbprKatakana): {wbKatakana, false, 130},

	// WB13a.
	int(wbAny)*wbPropertyMax + int(wbprExtendNumLet):          {wbExtendNumLet, true, 9990},
	int(wbALetter)*wbPropertyMax + int(wbprExtendNumLet):      {wbExtendNumLet, false, 131},
	int(wbHebrewLetter)*wbPropertyMax + int(wbprExtendNumLet): {wbExtendNumLet, false, 131},
	int(wbNumeric)*wbPropertyMax + int(wbprExtendNumLet):      {wbExtendNumLet, false, 131},
	int(wbKatakana)*wbPropertyMax + int(wbprExtendNumLet):     {wbExtendNumLet, false, 131},
	int(wbExtendNumLet)*wbPropertyMax + int(wbprExtendNumLet): {wbExtendNumLet, false, 131},

	// WB13b.
	int(wbExtendNumLet)*wbPropertyMax + int(wbprALetter):      {wbALetter, false, 132},
	int(wbExtendNumLet)*wbPropertyMax + int(wbprHebrewLetter): {wbHebrewLetter, false, 132},
	int(wbExtendNumLet)*wbPropertyMax + int(wbprNumeric):      {wbNumeric, false, 132},
	int(wbExtendNumLet)*wbPropertyMax + int(wbprKatakana):     {wbKatakana, false, 132},
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
	transition := wbTransitions[int(state)*wbPropertyMax+int(nextProperty)]
	if transition.ruleNumber > 0 {
		// We have a specific transition. We'll use it.
		newState, wordBreak, rule = transition.WordBreakState, transition.boundary, transition.ruleNumber
	} else {
		// No specific transition found. Try the less specific ones.
		transAnyProp := wbTransitions[int(state)*wbPropertyMax+int(wbprAny)]
		transAnyState := wbTransitions[int(wbAny)*wbPropertyMax+int(nextProperty)]
		if transAnyProp.ruleNumber > 0 && transAnyState.ruleNumber > 0 {
			// Both apply. We'll use a mix (see comments for grTransitions).
			newState, wordBreak, rule = transAnyState.WordBreakState, transAnyState.boundary, transAnyState.ruleNumber
			if transAnyProp.ruleNumber < transAnyState.ruleNumber {
				wordBreak, rule = transAnyProp.boundary, transAnyProp.ruleNumber
			}
		} else if transAnyProp.ruleNumber > 0 {
			// We only have a specific state.
			newState, wordBreak, rule = transAnyProp.WordBreakState, transAnyProp.boundary, transAnyProp.ruleNumber
			// This branch will probably never be reached because okAnyState will
			// always be true given the current transition map. But we keep it here
			// for future modifications to the transition map where this may not be
			// true anymore.
		} else if transAnyState.ruleNumber > 0 {
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
