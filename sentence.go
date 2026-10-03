package uniseg

import (
	"iter"
	"unicode/utf8"
)

// Sentences returns a sequence of sentences found in the given byte slice
// according to the rules of [Unicode Standard Annex #29, Sentence Boundaries].
// The sequence yields the starting index and the sentence as a byte slice.
//
// [Unicode Standard Annex #29, Sentence Boundaries]: https://www.unicode.org/reports/tr29/tr29-45.html#Sentence_Boundaries
func Sentences(b []byte) iter.Seq2[int, []byte] {
	return DefaultParser.Sentences(b)
}

// Sentences returns a sequence of sentences found in the given byte slice
// according to the rules of [Unicode Standard Annex #29, Sentence Boundaries].
// The sequence yields the starting index and the sentence as a byte slice.
//
// [Unicode Standard Annex #29, Sentence Boundaries]: https://www.unicode.org/reports/tr29/tr29-45.html#Sentence_Boundaries
func (p *Parser) Sentences(b []byte) iter.Seq2[int, []byte] {
	return func(yield func(idx int, sentence []byte) bool) {
		var state SentenceBreakState
		var index int
		rest := b
		for len(rest) > 0 {
			var sentence []byte
			sentence, rest, state = p.FirstSentence(rest, state)
			if !yield(index, sentence) {
				break
			}
			index += len(sentence)
		}
	}
}

// SentencesInString is like [Sentences] but its input and outputs are strings.
func SentencesInString(str string) iter.Seq2[int, string] {
	return DefaultParser.SentencesInString(str)
}

// SentencesInString is like [*Parser.Sentences] but its input and outputs are strings.
func (p *Parser) SentencesInString(str string) iter.Seq2[int, string] {
	return func(yield func(idx int, sentence string) bool) {
		var state SentenceBreakState
		var index int
		rest := str
		for len(rest) > 0 {
			var sentence string
			sentence, rest, state = p.FirstSentenceInString(rest, state)
			if !yield(index, sentence) {
				break
			}
			index += len(sentence)
		}
	}
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
// [Unicode Standard Annex #29, Sentence Boundaries]: https://www.unicode.org/reports/tr29/tr29-45.html#Sentence_Boundaries
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
// [Unicode Standard Annex #29, Sentence Boundaries]: https://www.unicode.org/reports/tr29/tr29-45.html#Sentence_Boundaries
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
