package uniseg

import "unicode/utf8"

// FirstSentence returns the first sentence found in the given byte slice
// according to the rules of [Unicode Standard Annex #29, Sentence Boundaries].
// This function can be called continuously to extract all sentences from a byte
// slice, as illustrated in the example below.
//
// If you don't know the current state, for example when calling the function
// for the first time, you must pass 0. For consecutive calls, pass the state
// and rest slice returned by the previous call.
//
// The "rest" slice is the sub-slice of the original byte slice "b" starting
// after the last byte of the identified sentence. If the length of the "rest"
// slice is 0, the entire byte slice "b" has been processed. The "sentence" byte
// slice is the sub-slice of the input slice containing the identified sentence.
//
// Given an empty byte slice "b", the function returns nil values.
//
// [Unicode Standard Annex #29, Sentence Boundaries]: https://www.unicode.org/reports/tr29/tr29-41.html#Sentence_Boundaries
func FirstSentence(b []byte, state SentenceBreakState) (sentence, rest []byte, newState SentenceBreakState) {
	return firstSentence(b, state, utf8.DecodeRune)
}

// FirstSentence returns the first sentence found in the given byte slice
// according to the rules of [Unicode Standard Annex #29, Sentence Boundaries].
// This function can be called continuously to extract all sentences from a byte
// slice, as illustrated in the example below.
//
// If you don't know the current state, for example when calling the function
// for the first time, you must pass 0. For consecutive calls, pass the state
// and rest slice returned by the previous call.
//
// The "rest" slice is the sub-slice of the original byte slice "b" starting
// after the last byte of the identified sentence. If the length of the "rest"
// slice is 0, the entire byte slice "b" has been processed. The "sentence" byte
// slice is the sub-slice of the input slice containing the identified sentence.
//
// Given an empty byte slice "b", the function returns nil values.
//
// [Unicode Standard Annex #29, Sentence Boundaries]: https://www.unicode.org/reports/tr29/tr29-41.html#Sentence_Boundaries
func (*Parser) FirstSentence(b []byte, state SentenceBreakState) (sentence, rest []byte, newState SentenceBreakState) {
	return firstSentence(b, state, utf8.DecodeRune)
}

// FirstSentenceInString is like [FirstSentence] but its input and outputs are
// strings.
func FirstSentenceInString(str string, state SentenceBreakState) (sentence, rest string, newState SentenceBreakState) {
	return firstSentence(str, state, utf8.DecodeRuneInString)
}

// FirstSentenceInString is like [Parser.FirstSentence] but its input and outputs are
// strings.
func (*Parser) FirstSentenceInString(str string, state SentenceBreakState) (sentence, rest string, newState SentenceBreakState) {
	return firstSentence(str, state, utf8.DecodeRuneInString)
}

func firstSentence[T bytes](str T, state SentenceBreakState, decoder runeDecoder[T]) (sentence, rest T, newState SentenceBreakState) {
	var zero T

	// An empty byte slice returns nothing.
	if len(str) == 0 {
		return
	}

	// Extract the first rune.
	r, length := decoder(str)
	if len(str) <= length { // If we're already past the end, there is nothing else to parse.
		return str, zero, sbAny
	}

	// If we don't know the state, determine it now.
	if state <= 0 {
		state, _ = transitionSentenceBreakState(state, r, str[length:], decoder)
	}

	// Transition until we find a boundary.
	var boundary bool
	for {
		r, l := decoder(str[length:])
		state, boundary = transitionSentenceBreakState(state, r, str[length+l:], decoder)

		if boundary {
			return str[:length], str[length:], state
		}

		length += l
		if len(str) <= length {
			return str, zero, sbAny
		}
	}
}
