package uniseg

import "unicode/utf8"

// SentenceBreakState is the state of the sentence break parser.
type SentenceBreakState int

// The states of the sentence break parser.
const (
	_ SentenceBreakState = iota // The zero value is reserved for the initial state.
	sbAny
	sbCR
	sbParaSep
	sbATerm
	sbUpper
	sbLower
	sbSB7
	sbSB8Close
	sbSB8Sp
	sbSTerm
	sbSB8aClose
	sbSB8aSp
)

type sbTransitionResult struct {
	SentenceBreakState
	boundary   bool
	ruleNumber int
}

// The sentence break parser's state transitions. It's analogous to
// grTransitions, see comments there for details. Unicode version 15.0.0.
func sbTransitions(sb SentenceBreakState, pr property) (sbTransitionResult, bool) {
	switch sb {
	case sbAny:
		switch pr {
		case prCR: // SB3.
			return sbTransitionResult{sbCR, false, 9990}, true
		case prSep, prLF: // SB4.
			return sbTransitionResult{sbParaSep, false, 9990}, true
		case prATerm: // SB6.
			return sbTransitionResult{sbATerm, false, 9990}, true
		case prUpper: // SB7.
			return sbTransitionResult{sbUpper, false, 9990}, true
		case prLower: // SB7.
			return sbTransitionResult{sbLower, false, 9990}, true
		case prSTerm: // SB8a.
			return sbTransitionResult{sbSTerm, false, 9990}, true
		}
	case sbCR:
		switch pr {
		case prLF: // SB3.
			return sbTransitionResult{sbParaSep, false, 30}, true
		case prAny: // SB4.
			return sbTransitionResult{sbAny, true, 40}, true
		}
	case sbParaSep:
		switch pr {
		case prAny: // SB4.
			return sbTransitionResult{sbAny, true, 40}, true
		}
	case sbATerm:
		switch pr {
		case prNumeric: // SB6.
			return sbTransitionResult{sbAny, false, 60}, true
		case prSContinue: // SB8a.
			return sbTransitionResult{sbAny, false, 81}, true
		case prATerm: // SB8a.
			return sbTransitionResult{sbATerm, false, 81}, true
		case prSTerm: // SB8a.
			return sbTransitionResult{sbSTerm, false, 81}, true
		case prClose: // SB9.
			return sbTransitionResult{sbSB8Close, false, 90}, true
		case prSp: // SB9.
			return sbTransitionResult{sbSB8Sp, false, 90}, true
		case prSep, prCR, prLF: // SB9.
			return sbTransitionResult{sbParaSep, false, 90}, true
		case prAny: // SB11.
			return sbTransitionResult{sbAny, true, 110}, true
		}
	case sbUpper:
		switch pr {
		case prATerm: // SB7.
			return sbTransitionResult{sbSB7, false, 70}, true
		}
	case sbLower:
		switch pr {
		case prATerm: // SB7.
			return sbTransitionResult{sbSB7, false, 70}, true
		}
	case sbSB7:
		switch pr {
		case prNumeric: // SB6.
			// Because ATerm also appears in SB7.
			return sbTransitionResult{sbAny, false, 60}, true
		case prUpper: // SB7.
			return sbTransitionResult{sbUpper, false, 70}, true
		case prSContinue: // SB8a.
			return sbTransitionResult{sbAny, false, 81}, true
		case prATerm: // SB8a.
			return sbTransitionResult{sbATerm, false, 81}, true
		case prSTerm: // SB8a.
			return sbTransitionResult{sbSTerm, false, 81}, true
		case prClose: // SB9.
			return sbTransitionResult{sbSB8Close, false, 90}, true
		case prSp: // SB9.
			return sbTransitionResult{sbSB8Sp, false, 90}, true
		case prSep, prCR, prLF: // SB9.
			return sbTransitionResult{sbParaSep, false, 90}, true
		case prAny: // SB11.
			return sbTransitionResult{sbAny, true, 110}, true
		}
	case sbSB8Close:
		switch pr {
		case prSContinue: // SB8a.
			return sbTransitionResult{sbAny, false, 81}, true
		case prATerm: // SB8a.
			return sbTransitionResult{sbATerm, false, 81}, true
		case prSTerm: // SB8a.
			return sbTransitionResult{sbSTerm, false, 81}, true
		case prClose: // SB9.
			return sbTransitionResult{sbSB8Close, false, 90}, true
		case prSp: // SB9.
			return sbTransitionResult{sbSB8Sp, false, 90}, true
		case prSep, prCR, prLF: // SB9.
			return sbTransitionResult{sbParaSep, false, 90}, true
		case prAny: // SB11.
			return sbTransitionResult{sbAny, true, 110}, true
		}
	case sbSB8Sp:
		switch pr {
		case prSContinue: // SB8a.
			return sbTransitionResult{sbAny, false, 81}, true
		case prATerm: // SB8a.
			return sbTransitionResult{sbATerm, false, 81}, true
		case prSTerm: // SB8a.
			return sbTransitionResult{sbSTerm, false, 81}, true
		case prSp: // SB10.
			return sbTransitionResult{sbSB8Sp, false, 100}, true
		case prSep, prCR, prLF: // SB10.
			return sbTransitionResult{sbParaSep, false, 100}, true
		case prAny: // SB11.
			return sbTransitionResult{sbAny, true, 110}, true
		}
	case sbSTerm:
		switch pr {
		case prSContinue: // SB8a.
			return sbTransitionResult{sbAny, false, 81}, true
		case prATerm: // SB8a.
			return sbTransitionResult{sbATerm, false, 81}, true
		case prSTerm: // SB8a.
			return sbTransitionResult{sbSTerm, false, 81}, true
		case prClose: // SB9.
			return sbTransitionResult{sbSB8aClose, false, 90}, true
		case prSp: // SB9.
			return sbTransitionResult{sbSB8aSp, false, 90}, true
		case prSep, prCR, prLF: // SB9.
			return sbTransitionResult{sbParaSep, false, 90}, true
		case prAny: // SB11.
			return sbTransitionResult{sbAny, true, 110}, true
		}
	case sbSB8aClose:
		switch pr {
		case prSContinue: // SB8a.
			return sbTransitionResult{sbAny, false, 81}, true
		case prATerm: // SB8a.
			return sbTransitionResult{sbATerm, false, 81}, true
		case prSTerm: // SB8a.
			return sbTransitionResult{sbSTerm, false, 81}, true
		case prClose: // SB9.
			return sbTransitionResult{sbSB8aClose, false, 90}, true
		case prSp: // SB9.
			return sbTransitionResult{sbSB8aSp, false, 90}, true
		case prSep, prCR, prLF: // SB9.
			return sbTransitionResult{sbParaSep, false, 90}, true
		case prAny: // SB11.
			return sbTransitionResult{sbAny, true, 110}, true
		}
	case sbSB8aSp:
		switch pr {
		case prSContinue: // SB8a.
			return sbTransitionResult{sbAny, false, 81}, true
		case prATerm: // SB8a.
			return sbTransitionResult{sbATerm, false, 81}, true
		case prSTerm: // SB8a.
			return sbTransitionResult{sbSTerm, false, 81}, true
		case prSp: // SB10.
			return sbTransitionResult{sbSB8aSp, false, 100}, true
		case prAny: // SB11.
			return sbTransitionResult{sbAny, true, 110}, true
		}
	}

	return sbTransitionResult{}, false
}

// transitionSentenceBreakState determines the new state of the sentence break
// parser given the current state and the next code point. It also returns
// whether a sentence boundary was detected. If more than one code point is
// needed to determine the new state, the byte slice or the string starting
// after rune "r" can be used (whichever is not nil or empty) for further
// lookups.
func transitionSentenceBreakState[T bytes](state SentenceBreakState, r rune, str T, decoder runeDecoder[T]) (newState SentenceBreakState, sentenceBreak bool) {
	// Determine the property of the next character.
	nextProperty := sentenceBreakCodePoints.search(r)

	// SB5 (Replacing Ignore Rules).
	if nextProperty == prExtend || nextProperty == prFormat {
		if state == sbParaSep || state == sbCR {
			return sbAny, true // Make sure we don't apply SB5 to SB3 or SB4.
		}
		if state < 0 {
			return sbAny, true // SB1.
		}
		return state, false
	}

	// Find the applicable transition in the table.
	var rule int
	transition, ok := sbTransitions(state, nextProperty)
	if ok {
		// We have a specific transition. We'll use it.
		newState, sentenceBreak, rule = transition.SentenceBreakState, transition.boundary, transition.ruleNumber
	} else {
		// No specific transition found. Try the less specific ones.
		transAnyProp, okAnyProp := sbTransitions(state, prAny)
		transAnyState, okAnyState := sbTransitions(sbAny, nextProperty)
		if okAnyProp && okAnyState {
			// Both apply. We'll use a mix (see comments for grTransitions).
			newState, sentenceBreak, rule = transAnyState.SentenceBreakState, transAnyState.boundary, transAnyState.ruleNumber
			if transAnyProp.ruleNumber < transAnyState.ruleNumber {
				sentenceBreak, rule = transAnyProp.boundary, transAnyProp.ruleNumber
			}
		} else if okAnyProp {
			// We only have a specific state.
			newState, sentenceBreak, rule = transAnyProp.SentenceBreakState, transAnyProp.boundary, transAnyProp.ruleNumber
			// This branch will probably never be reached because okAnyState will
			// always be true given the current transition map. But we keep it here
			// for future modifications to the transition map where this may not be
			// true anymore.
		} else if okAnyState {
			// We only have a specific property.
			newState, sentenceBreak, rule = transAnyState.SentenceBreakState, transAnyState.boundary, transAnyState.ruleNumber
		} else {
			// No known transition. SB999: Any × Any.
			newState, sentenceBreak, rule = sbAny, false, 9990
		}
	}

	// SB8.
	if rule > 80 && (state == sbATerm || state == sbSB8Close || state == sbSB8Sp || state == sbSB7) {
		// Check the right side of the rule.
		var length int
		for nextProperty != prOLetter &&
			nextProperty != prUpper &&
			nextProperty != prLower &&
			nextProperty != prSep &&
			nextProperty != prCR &&
			nextProperty != prLF &&
			nextProperty != prATerm &&
			nextProperty != prSTerm {
			// Move on to the next rune.
			r, length = decoder(str)
			str = str[length:]
			if r == utf8.RuneError {
				break
			}
			nextProperty = sentenceBreakCodePoints.search(r)
		}
		if nextProperty == prLower {
			return sbLower, false
		}
	}

	return
}
