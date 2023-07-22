package uniseg

import "unicode/utf8"

// The bit masks used to extract boundary information returned by [Step].
const (
	MaskLine     = 3
	MaskWord     = 4
	MaskSentence = 8
)

// The number of bits to shift the boundary information returned by [Step] to
// obtain the monospace width of the grapheme cluster.
const ShiftWidth = 4

// The bit positions by which boundary flags are shifted by the [Step] function.
// These must correspond to the Mask constants.
const (
	shiftWord     = 2
	shiftSentence = 3
	// shiftWidth is ShiftWidth above. No mask as these are always the remaining bits.
)

// The bit positions by which states are shifted by the [Step] function. These
// values must ensure state values defined for each of the boundary algorithms
// don't overlap (and that they all still fit in a single int). These must
// correspond to the Mask constants.
const (
	shiftWordState     = 4
	shiftSentenceState = 9
	shiftLineState     = 13
	shiftPropState     = 21 // No mask as these are always the remaining bits.
)

// The bit mask used to extract the state returned by the [Step] function, after
// shifting. These values must correspond to the shift constants.
const (
	maskGraphemeState = 0xf
	maskWordState     = 0x1f
	maskSentenceState = 0xf
	maskLineState     = 0xff
)

// Step returns the first grapheme cluster (user-perceived character) found in
// the given byte slice. It also returns information about the boundary between
// that grapheme cluster and the one following it as well as the monospace width
// of the grapheme cluster. There are three types of boundary information: word
// boundaries, sentence boundaries, and line breaks. This function is therefore
// a combination of [FirstGraphemeCluster], [FirstWord], [FirstSentence], and
// [FirstLineSegment].
//
// The "boundaries" return value can be evaluated as follows:
//
//   - boundaries&MaskWord != 0: The boundary is a word boundary.
//   - boundaries&MaskWord == 0: The boundary is not a word boundary.
//   - boundaries&MaskSentence != 0: The boundary is a sentence boundary.
//   - boundaries&MaskSentence == 0: The boundary is not a sentence boundary.
//   - boundaries&MaskLine == LineDontBreak: You must not break the line at the
//     boundary.
//   - boundaries&MaskLine == LineMustBreak: You must break the line at the
//     boundary.
//   - boundaries&MaskLine == LineCanBreak: You may or may not break the line at
//     the boundary.
//   - boundaries >> ShiftWidth: The width of the grapheme cluster for most
//     monospace fonts where a value of 1 represents one character cell.
//
// This function can be called continuously to extract all grapheme clusters
// from a byte slice, as illustrated in the examples below.
//
// If you don't know which state to pass, for example when calling the function
// for the first time, you must pass 0. For consecutive calls, pass the state
// and rest slice returned by the previous call.
//
// The "rest" slice is the sub-slice of the original byte slice "b" starting
// after the last byte of the identified grapheme cluster. If the length of the
// "rest" slice is 0, the entire byte slice "b" has been processed. The
// "cluster" byte slice is the sub-slice of the input slice containing the
// first identified grapheme cluster.
//
// Given an empty byte slice "b", the function returns nil values.
//
// While slightly less convenient than using the Graphemes class, this function
// has much better performance and makes no allocations. It lends itself well to
// large byte slices.
//
// Note that in accordance with [UAX #14 LB3], the final segment will end with
// a mandatory line break (boundaries&MaskLine == LineMustBreak). You can choose
// to ignore this by checking if the length of the "rest" slice is 0 and calling
// [HasTrailingLineBreak] or [HasTrailingLineBreakInString] on the last rune.
//
// [UAX #14 LB3]: https://www.unicode.org/reports/tr14/tr14-49.html#Algorithm
func Step(b []byte, state int) (cluster, rest []byte, boundaries int, newState int) {
	return step(b, state, utf8.DecodeRune)
}

// StepString is like [Step] but its input and outputs are strings.
func StepString(str string, state int) (cluster, rest string, boundaries int, newState int) {
	return step(str, state, utf8.DecodeRuneInString)
}

func step[T bytes](str T, state int, decoder runeDecoder[T]) (cluster, rest T, boundaries int, newState int) {
	var zero T

	// An empty byte slice returns nothing.
	if len(str) == 0 {
		return
	}

	// Extract the first rune.
	r, length := decoder(str)
	if len(str) <= length { // If we're already past the end, there is nothing else to parse.
		prop := graphemeCodePoints.search(r)
		return str, zero, int(LineMustBreak) | (1 << shiftWord) | (1 << shiftSentence) | (runeWidth(r, prop) << ShiftWidth), int(grAny) | (int(wbAny) << shiftWordState) | (int(sbAny) << shiftSentenceState) | (int(lbAny) << shiftLineState)
	}

	// If we don't know the state, determine it now.
	var graphemeState grState
	var wordState WordBreakState
	var sentenceState SentenceBreakState
	var lineState LineBreakState
	var firstProp property
	remainder := str[length:]
	if state <= 0 {
		graphemeState, firstProp, _ = transitionGraphemeState(0, r)
		wordState, _ = transitionWordBreakState(0, r, remainder, decoder)
		sentenceState, _ = transitionSentenceBreakState(0, r, remainder, decoder)
		lineState, _ = transitionLineBreakState(0, r, remainder, decoder)
	} else {
		graphemeState = grState(state & maskGraphemeState)
		wordState = WordBreakState((state >> shiftWordState) & maskWordState)
		sentenceState = SentenceBreakState((state >> shiftSentenceState) & maskSentenceState)
		lineState = LineBreakState((state >> shiftLineState) & maskLineState)
		firstProp = property(state >> shiftPropState)
	}

	// Transition until we find a grapheme cluster boundary.
	width := runeWidth(r, firstProp)
	for {
		var (
			graphemeBoundary, wordBoundary, sentenceBoundary bool
			lineBreak                                        LineBreak
			prop                                             property
		)

		r, l := decoder(remainder)
		remainder = str[length+l:]

		graphemeState, prop, graphemeBoundary = transitionGraphemeState(graphemeState, r)
		wordState, wordBoundary = transitionWordBreakState(wordState, r, remainder, decoder)
		sentenceState, sentenceBoundary = transitionSentenceBreakState(sentenceState, r, remainder, decoder)
		lineState, lineBreak = transitionLineBreakState(lineState, r, remainder, decoder)

		if graphemeBoundary {
			boundary := int(lineBreak) | (width << ShiftWidth)
			if wordBoundary {
				boundary |= 1 << shiftWord
			}
			if sentenceBoundary {
				boundary |= 1 << shiftSentence
			}
			return str[:length], str[length:], boundary, int(graphemeState) | int(wordState<<shiftWordState) | int(sentenceState<<shiftSentenceState) | int(lineState<<shiftLineState) | int(prop<<shiftPropState)
		}

		if r == vs16 {
			width = 2
		} else if firstProp != prExtendedPictographic && firstProp != prRegionalIndicator && firstProp != prL {
			width += runeWidth(r, prop)
		} else if firstProp == prExtendedPictographic {
			if r == vs15 {
				width = 1
			} else {
				width = 2
			}
		}

		length += l
		if len(str) <= length {
			return str, zero, int(LineMustBreak) | (1 << shiftWord) | (1 << shiftSentence) | (width << ShiftWidth), int(grAny) | int(wbAny<<shiftWordState) | int(sbAny<<shiftSentenceState) | int(lbAny<<shiftLineState) | int(prop<<shiftPropState)
		}
	}
}
