package uniseg

import (
	"unicode/utf8"
)

// Graphemes implements an iterator over Unicode grapheme clusters, or
// user-perceived characters. While iterating, it also provides information
// about word boundaries, sentence boundaries, line breaks, and monospace
// character widths.
//
// After constructing the class via [NewGraphemes] for a given string "str",
// [Graphemes.Next] is called for every grapheme cluster in a loop until it
// returns false. Inside the loop, information about the grapheme cluster as
// well as boundary information and character width is available via the various
// methods (see examples below).
//
// Using this class to iterate over a string is convenient but it is much slower
// than using this package's [Step] or [StepString] functions or any of the
// other specialized functions starting with "First".
type Graphemes struct {
	parser *Parser

	// The original string.
	original string

	// The remaining string to be parsed.
	remaining string

	// The current grapheme cluster.
	cluster string

	// The byte offset of the current grapheme cluster relative to the original
	// string.
	offset int

	// The current boundary information of the [Step] parser.
	boundaries Boundaries

	// The current state of the [Step] parser.
	state State
}

// NewGraphemes returns a new grapheme cluster iterator.
func NewGraphemes(str string) *Graphemes {
	return &Graphemes{
		parser:    DefaultParser,
		original:  str,
		remaining: str,
	}
}

func (p *Parser) NewGraphemes(str string) *Graphemes {
	return &Graphemes{
		parser:    p,
		original:  str,
		remaining: str,
	}
}

// Next advances the iterator by one grapheme cluster and returns false if no
// clusters are left. This function must be called before the first cluster is
// accessed.
func (g *Graphemes) Next() bool {
	if len(g.remaining) == 0 {
		// We're already past the end.
		g.state = -2
		g.cluster = ""
		return false
	}
	g.offset += len(g.cluster)
	g.cluster, g.remaining, g.boundaries, g.state = step(g.parser, g.remaining, g.state, utf8.DecodeRuneInString)
	return true
}

// Runes returns a slice of runes (code points) which corresponds to the current
// grapheme cluster. If the iterator is already past the end or [Graphemes.Next]
// has not yet been called, nil is returned.
func (g *Graphemes) Runes() []rune {
	if g.state <= 0 {
		return nil
	}
	return []rune(g.cluster)
}

// Str returns a substring of the original string which corresponds to the
// current grapheme cluster. If the iterator is already past the end or
// [Graphemes.Next] has not yet been called, an empty string is returned.
func (g *Graphemes) Str() string {
	return g.cluster
}

// Bytes returns a byte slice which corresponds to the current grapheme cluster.
// If the iterator is already past the end or [Graphemes.Next] has not yet been
// called, nil is returned.
func (g *Graphemes) Bytes() []byte {
	if g.state <= 0 {
		return nil
	}
	return []byte(g.cluster)
}

// Positions returns the interval of the current grapheme cluster as byte
// positions into the original string. The first returned value "from" indexes
// the first byte and the second returned value "to" indexes the first byte that
// is not included anymore, i.e. str[from:to] is the current grapheme cluster of
// the original string "str". If [Graphemes.Next] has not yet been called, both
// values are 0. If the iterator is already past the end, both values are 1.
func (g *Graphemes) Positions() (int, int) {
	if g.state == -1 {
		return 0, 0
	} else if g.state == -2 {
		return 1, 1
	}
	return g.offset, g.offset + len(g.cluster)
}

// IsWordBoundary returns true if a word ends after the current grapheme
// cluster.
func (g *Graphemes) IsWordBoundary() bool {
	if g.state <= 0 {
		return true
	}
	return g.boundaries&maskWord != 0
}

// IsSentenceBoundary returns true if a sentence ends after the current
// grapheme cluster.
func (g *Graphemes) IsSentenceBoundary() bool {
	if g.state <= 0 {
		return true
	}
	return g.boundaries&maskSentence != 0
}

// LineBreak returns whether the line can be broken after the current grapheme
// cluster. A value of [LineDontBreak] means the line may not be broken, a value
// of [LineMustBreak] means the line must be broken, and a value of
// [LineCanBreak] means the line may or may not be broken.
func (g *Graphemes) LineBreak() LineBreak {
	if g.state == -1 {
		return LineDontBreak
	}
	if g.state == -2 {
		return LineMustBreak
	}
	return LineBreak(g.boundaries & maskLine)
}

// Width returns the monospace width of the current grapheme cluster.
func (g *Graphemes) Width() int {
	if g.state <= 0 {
		return 0
	}
	return g.boundaries.Width()
}

// Reset puts the iterator into its initial state such that the next call to
// [Graphemes.Next] sets it to the first grapheme cluster again.
func (g *Graphemes) Reset() {
	g.state = -1
	g.offset = 0
	g.cluster = ""
	g.remaining = g.original
}

// GraphemeClusterCount returns the number of user-perceived characters
// (grapheme clusters) for the given string.
func GraphemeClusterCount(s string) (n int) {
	var state GraphemeBreakState
	for len(s) > 0 {
		_, s, _, state = FirstGraphemeClusterInString(s, state)
		n++
	}
	return
}

// GraphemeClusterCount returns the number of user-perceived characters
// (grapheme clusters) for the given string.
func (p *Parser) GraphemeClusterCount(s string) (n int) {
	var state GraphemeBreakState
	for len(s) > 0 {
		_, s, _, state = p.FirstGraphemeClusterInString(s, state)
		n++
	}
	return
}

// ReverseString reverses the given string while observing grapheme cluster
// boundaries.
func ReverseString(s string) string {
	str := []byte(s)
	reversed := make([]byte, len(str))
	var state GraphemeBreakState
	index := len(str)
	for len(str) > 0 {
		var cluster []byte
		cluster, str, _, state = FirstGraphemeCluster(str, state)
		index -= len(cluster)
		copy(reversed[index:], cluster)
		if index <= len(str)/2 {
			break
		}
	}
	return string(reversed)
}

// GraphemeBreakState the type of the grapheme cluster parser's states.
type GraphemeBreakState int

// The number of bits the grapheme property must be shifted to make place for
// grapheme states.
const shiftGraphemePropState = 4

func newGraphemeBreakState(s grState, p property) GraphemeBreakState {
	return GraphemeBreakState(s)<<shiftGraphemePropState | GraphemeBreakState(p)
}

func (s GraphemeBreakState) unpack() (grState, property) {
	return grState(s >> shiftGraphemePropState), property(s & ((1 << shiftGraphemePropState) - 1))
}

// FirstGraphemeCluster returns the first grapheme cluster found in the given
// byte slice according to the rules of [Unicode Standard Annex #29, Grapheme
// Cluster Boundaries]. This function can be called continuously to extract all
// grapheme clusters from a byte slice, as illustrated in the example below.
//
// If you don't know the current state, for example when calling the function
// for the first time, you must pass 0. For consecutive calls, pass the state
// and rest slice returned by the previous call.
//
// The "rest" slice is the sub-slice of the original byte slice "b" starting
// after the last byte of the identified grapheme cluster. If the length of the
// "rest" slice is 0, the entire byte slice "b" has been processed. The
// "cluster" byte slice is the sub-slice of the input slice containing the
// identified grapheme cluster.
//
// The returned width is the width of the grapheme cluster for most monospace
// fonts where a value of 1 represents one character cell.
//
// Given an empty byte slice "b", the function returns nil values.
//
// While slightly less convenient than using the Graphemes class, this function
// has much better performance and makes no allocations. It lends itself well to
// large byte slices.
//
// [Unicode Standard Annex #29, Grapheme Cluster Boundaries]: https://www.unicode.org/reports/tr29/tr29-41.html#Grapheme_Cluster_Boundaries
func FirstGraphemeCluster(b []byte, state GraphemeBreakState) (cluster, rest []byte, width int, newState GraphemeBreakState) {
	return firstGraphemeCluster(DefaultParser, b, state, utf8.DecodeRune)
}

// FirstGraphemeCluster returns the first grapheme cluster found in the given
// byte slice according to the rules of [Unicode Standard Annex #29, Grapheme
// Cluster Boundaries]. This function can be called continuously to extract all
// grapheme clusters from a byte slice, as illustrated in the example below.
//
// If you don't know the current state, for example when calling the function
// for the first time, you must pass 0. For consecutive calls, pass the state
// and rest slice returned by the previous call.
//
// The "rest" slice is the sub-slice of the original byte slice "b" starting
// after the last byte of the identified grapheme cluster. If the length of the
// "rest" slice is 0, the entire byte slice "b" has been processed. The
// "cluster" byte slice is the sub-slice of the input slice containing the
// identified grapheme cluster.
//
// The returned width is the width of the grapheme cluster for most monospace
// fonts where a value of 1 represents one character cell.
//
// Given an empty byte slice "b", the function returns nil values.
//
// While slightly less convenient than using the Graphemes class, this function
// has much better performance and makes no allocations. It lends itself well to
// large byte slices.
//
// [Unicode Standard Annex #29, Grapheme Cluster Boundaries]: https://www.unicode.org/reports/tr29/tr29-41.html#Grapheme_Cluster_Boundaries
func (p *Parser) FirstGraphemeCluster(b []byte, state GraphemeBreakState) (cluster, rest []byte, width int, newState GraphemeBreakState) {
	return firstGraphemeCluster(p, b, state, utf8.DecodeRune)
}

// FirstGraphemeClusterInString is like [FirstGraphemeCluster] but its input and
// outputs are strings.
func FirstGraphemeClusterInString(str string, state GraphemeBreakState) (cluster, rest string, width int, newState GraphemeBreakState) {
	return firstGraphemeCluster(DefaultParser, str, state, utf8.DecodeRuneInString)
}

// FirstGraphemeClusterInString is like [Parser.FirstGraphemeCluster] but its input and
// outputs are strings.
func (p *Parser) FirstGraphemeClusterInString(str string, state GraphemeBreakState) (cluster, rest string, width int, newState GraphemeBreakState) {
	return firstGraphemeCluster(p, str, state, utf8.DecodeRuneInString)
}

func firstGraphemeCluster[T bytes](p *Parser, str T, state GraphemeBreakState, decoder runeDecoder[T]) (cluster, rest T, width int, newState GraphemeBreakState) {
	var zero T

	// An empty string returns nothing.
	if len(str) == 0 {
		return
	}

	// Extract the first rune.
	r, length := decoder(str)
	if len(str) <= length { // If we're already past the end, there is nothing else to parse.
		var prop property
		if state <= 0 {
			prop = graphemeCodePoints.search(r)
		} else {
			_, prop = state.unpack()
		}
		return str, zero, runeWidth(p, r, prop), newGraphemeBreakState(grAny, prop)
	}

	// If we don't know the state, determine it now.
	var myState grState
	var firstProp property
	if state <= 0 {
		myState, firstProp, _ = transitionGraphemeState(myState, r)
	} else {
		myState, firstProp = state.unpack()
	}
	width += runeWidth(p, r, firstProp)

	// Transition until we find a boundary.
	for {
		var (
			prop     property
			boundary bool
		)

		r, l := decoder(str[length:])
		myState, prop, boundary = transitionGraphemeState(myState, r)

		if boundary {
			return str[:length], str[length:], width, newGraphemeBreakState(myState, prop)
		}

		if firstProp == prExtendedPictographic {
			if r == vs15 {
				width = 1
			} else if r == vs16 {
				width = 2
			}
		} else if firstProp != prRegionalIndicator && firstProp != prL {
			width += runeWidth(p, r, prop)
		}

		length += l
		if len(str) <= length {
			return str, zero, width, newGraphemeBreakState(grAny, prop)
		}
	}
}
