package uniseg

import "unicode/utf8"

// State is the type of the state of the [Step] parser.
type State int

func newState(gr grState, wb WordBreakState, sb SentenceBreakState, lb LineBreakState, prop property) State {
	return State(gr) |
		State(wb<<shiftWordState) |
		State(sb<<shiftSentenceState) |
		// State(lb<<shiftLineState) |
		State(prop<<shiftPropState)
}

func (s State) unpack() (gr grState, wb WordBreakState, sb SentenceBreakState, lb LineBreakState, prop property) {
	gr = grState(s & maskGraphemeState)
	wb = WordBreakState((s >> shiftWordState) & maskWordState)
	sb = SentenceBreakState((s >> shiftSentenceState) & maskSentenceState)
	// lb = LineBreakState((s >> shiftLineState) & maskLineState)
	prop = property(s >> shiftPropState)
	return
}

// Boundaries is the type of the boundary information returned by [Step].
type Boundaries int

func newBoundaries(lb LineBreak, wb bool, sb bool, width int) Boundaries {
	var b Boundaries
	b |= Boundaries(lb<<shiftLine) | Boundaries(width<<shiftWidth)
	if wb {
		b |= 1 << shiftWord
	}
	if sb {
		b |= 1 << shiftSentence
	}
	return b
}

// Line returns the line break information from b.
func (b Boundaries) Line() LineBreak {
	return LineBreak(b&maskLine) >> shiftLine
}

// Word returns the word break information from b.
func (b Boundaries) Word() bool {
	return b&maskWord != 0
}

// Sentence returns the sentence break information from b.
func (b Boundaries) Sentence() bool {
	return b&maskSentence != 0
}

// Width returns the width information from b.
func (b Boundaries) Width() int {
	return int(b) >> shiftWidth
}

// The bit masks used to extract boundary information returned by [Step].
const (
	maskLine     = 0b0_0_11
	maskWord     = 0b0_1_00
	maskSentence = 0b1_0_00
)

// The bit positions by which boundary flags are shifted by the [Step] function.
// These must correspond to the Mask constants.
const (
	shiftLine     = 0
	shiftWord     = 2
	shiftSentence = 3
	shiftWidth    = 4
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
// a mandatory line break (boundaries&maskLine == LineMustBreak). You can choose
// to ignore this by checking if the length of the "rest" slice is 0 and calling
// [HasTrailingLineBreak] or [HasTrailingLineBreakInString] on the last rune.
//
// [UAX #14 LB3]: https://www.unicode.org/reports/tr14/tr14-49.html#Algorithm
func Step(b []byte, state State) (cluster, rest []byte, boundaries Boundaries, newState State) {
	return step(DefaultParser, b, state, utf8.DecodeRune)
}

// Step returns the first grapheme cluster (user-perceived character) found in
// the given byte slice. It also returns information about the boundary between
// that grapheme cluster and the one following it as well as the monospace width
// of the grapheme cluster. There are three types of boundary information: word
// boundaries, sentence boundaries, and line breaks. This function is therefore
// a combination of [FirstGraphemeCluster], [FirstWord], [FirstSentence], and
// [FirstLineSegment].
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
// a mandatory line break (boundaries&maskLine == LineMustBreak). You can choose
// to ignore this by checking if the length of the "rest" slice is 0 and calling
// [HasTrailingLineBreak] or [HasTrailingLineBreakInString] on the last rune.
//
// [UAX #14 LB3]: https://www.unicode.org/reports/tr14/tr14-49.html#Algorithm
func (p *Parser) Step(b []byte, state State) (cluster, rest []byte, boundaries Boundaries, newState State) {
	return step(p, b, state, utf8.DecodeRune)
}

// StepString is like [Step] but its input and outputs are strings.
func StepString(str string, state State) (cluster, rest string, boundaries Boundaries, newState State) {
	return step(DefaultParser, str, state, utf8.DecodeRuneInString)
}

// StepString is like [Parser.Step] but its input and outputs are strings.
func (p *Parser) StepString(str string, state State) (cluster, rest string, boundaries Boundaries, newState State) {
	return step(p, str, state, utf8.DecodeRuneInString)
}

func step[T bytes](p *Parser, str T, state State, decoder runeDecoder[T]) (cluster, rest T, boundaries Boundaries, _newState State) {
	var zero T

	// An empty byte slice returns nothing.
	if len(str) == 0 {
		return
	}

	// Extract the first rune.
	r, length := decoder(str)
	if len(str) <= length { // If we're already past the end, there is nothing else to parse.
		prop := graphemeCodePoints.search(r)
		boundaries := newBoundaries(LineMustBreak, true, true, runeWidth(p, r, prop))
		_newState := newState(grAny, wbAny, sbAny, *new(LineBreakState), prop)
		return str, zero, boundaries, _newState
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
		lineState, _ = transitionLineBreakState(*new(LineBreakState), r, remainder, decoder)
	} else {
		graphemeState, wordState, sentenceState, lineState, firstProp = state.unpack()
	}
	width := runeWidth(p, r, firstProp)

	// Transition until we find a grapheme cluster boundary.
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
			boundary := newBoundaries(lineBreak, wordBoundary, sentenceBoundary, width)
			_newState := newState(graphemeState, wordState, sentenceState, lineState, prop)
			return str[:length], str[length:], boundary, _newState
		}

		if r == vs16 {
			width = 2
		} else if firstProp != prExtendedPictographic && firstProp != prRegionalIndicator && firstProp != prL {
			width += runeWidth(p, r, prop)
		} else if firstProp == prExtendedPictographic {
			if r == vs15 {
				width = 1
			} else {
				width = 2
			}
		}

		length += l
		if len(str) <= length {
			boundaries := newBoundaries(LineMustBreak, true, true, width)
			_newState := newState(grAny, wbAny, sbAny, *new(LineBreakState), prop)
			return str, zero, boundaries, _newState
		}
	}
}
