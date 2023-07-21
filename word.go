package uniseg

import "unicode/utf8"

// FirstWord returns the first word found in the given byte slice according to
// the rules of [Unicode Standard Annex #29, Word Boundaries]. This function can
// be called continuously to extract all words from a byte slice, as illustrated
// in the example below.
//
// If you don't know the current state, for example when calling the function
// for the first time, you must pass 0. For consecutive calls, pass the state
// and rest slice returned by the previous call.
//
// The "rest" slice is the sub-slice of the original byte slice "b" starting
// after the last byte of the identified word. If the length of the "rest" slice
// is 0, the entire byte slice "b" has been processed. The "word" byte slice is
// the sub-slice of the input slice containing the identified word.
//
// Given an empty byte slice "b", the function returns nil values.
//
// [Unicode Standard Annex #29, Word Boundaries]: https://www.unicode.org/reports/tr29/tr29-41.html#Word_Boundaries
func FirstWord(b []byte, state WordBreakState) (word, rest []byte, newState WordBreakState) {
	return firstWord(b, state, utf8.DecodeRune)
}

// FirstWordInString is like [FirstWord] but its input and outputs are strings.
func FirstWordInString(str string, state WordBreakState) (word, rest string, newState WordBreakState) {
	return firstWord(str, state, utf8.DecodeRuneInString)
}

func firstWord[T bytes](str T, state WordBreakState, decoder runeDecoder[T]) (word, rest T, newState WordBreakState) {
	var zero T

	// An empty byte slice returns nothing.
	if len(str) == 0 {
		return
	}

	// Extract the first rune.
	r, length := decoder(str)
	if len(str) <= length { // If we're already past the end, there is nothing else to parse.
		return str, zero, wbAny
	}

	// If we don't know the state, determine it now.
	if state <= 0 {
		state, _ = transitionWordBreakState(state, r, str[length:], decoder)
	}

	// Transition until we find a boundary.
	var boundary bool
	for {
		r, l := decoder(str[length:])
		state, boundary = transitionWordBreakState(state, r, str[length+l:], decoder)

		if boundary {
			return str[:length], str[length:], state
		}

		length += l
		if len(str) <= length {
			return str, zero, wbAny
		}
	}
}
