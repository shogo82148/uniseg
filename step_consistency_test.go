package uniseg

import (
	"math"
	"slices"
	"testing"
)

// stepConsistencyExtraInputs are targeted inputs that exercise line break
// state flags which must survive being packed into a [State].
var stepConsistencyExtraInputs = []string{
	"a b",
	"a b c",
	"hello world, this is a test.",
	"a - b",
	"-a b",
	"a -b",
	" -1",
	"well-known self-evident",
	"א-ב ג", // Hebrew letters with hyphen (LB21a).
	"א‐ב",   // Hebrew letter, hyphen, Hebrew letter.
	"◌क ◌",  // Dotted circles (LB28a).
	"ᨠ◌᩠ x", // Dotted circle between aksara.
	"( a ) b",
	"“a” b",          // Quotation marks (LB15).
	"« a » b",        // Pi/Pf quotation marks with spaces.
	"x‍\U0001f600 y", // ZWJ and Extended_Pictographic.
	"\U0001f1ef\U0001f1f5 \U0001f1fa\U0001f1f8",
	"、（a） b", // East Asian wide punctuation.
	"Mr. Smith went. He said \"Hi.\" Then (left.)  Next",
	"etc.)’ the",
	"a.\r\nb",
	"Hellò Ẃorld!",
}

// stepBoundaries runs [StepString] over str and collects the byte offsets of
// grapheme cluster, word, sentence, and line boundaries, as well as the width
// of each grapheme cluster.
func stepBoundaries(str string) (clusters, words, sentences, lines, mustLines, widths []int) {
	var (
		state      State
		boundaries Boundaries
		pos        int
	)
	rest := str
	for len(rest) > 0 {
		var c string
		c, rest, boundaries, state = StepString(rest, state)
		pos += len(c)
		clusters = append(clusters, pos)
		widths = append(widths, boundaries.Width())
		if boundaries.Word() {
			words = append(words, pos)
		}
		if boundaries.Sentence() {
			sentences = append(sentences, pos)
		}
		switch boundaries.Line() {
		case LineCanBreak:
			lines = append(lines, pos)
		case LineMustBreak:
			lines = append(lines, pos)
			mustLines = append(mustLines, pos)
		}
	}
	return
}

// keepAt returns the elements of offsets which are also contained in at.
func keepAt(offsets, at []int) []int {
	var kept []int
	for _, o := range offsets {
		if _, found := slices.BinarySearch(at, o); found {
			kept = append(kept, o)
		}
	}
	return kept
}

// checkStepConsistency verifies that [StepString] reports the same boundaries
// and widths as the standalone functions.
func checkStepConsistency(t *testing.T, str string) {
	t.Helper()
	clusters, words, sentences, lines, mustLines, widths := stepBoundaries(str)

	// Grapheme clusters and widths.
	var (
		expClusters []int
		expWidths   []int
		grState     GraphemeBreakState
		pos         int
	)
	for rest := str; len(rest) > 0; {
		var (
			c string
			w int
		)
		c, rest, w, grState = FirstGraphemeClusterInString(rest, grState)
		pos += len(c)
		expClusters = append(expClusters, pos)
		expWidths = append(expWidths, w)
	}
	if !slices.Equal(clusters, expClusters) {
		t.Errorf("%q: Step clusters %v, FirstGraphemeCluster %v", str, clusters, expClusters)
		return // The remaining comparisons depend on matching clusters.
	}
	if !slices.Equal(widths, expWidths) {
		t.Errorf("%q: Step widths %v, FirstGraphemeCluster widths %v", str, widths, expWidths)
	}
	var total int
	for _, w := range widths {
		total += w
	}
	if expected := StringWidth(str); total != expected {
		t.Errorf("%q: Step total width %d, StringWidth %d", str, total, expected)
	}

	// Words. Step only reports boundaries at grapheme cluster boundaries.
	var expWords []int
	var wbState WordBreakState
	pos = 0
	for rest := str; len(rest) > 0; {
		var w string
		w, rest, wbState = FirstWordInString(rest, wbState)
		pos += len(w)
		expWords = append(expWords, pos)
	}
	if expWords = keepAt(expWords, clusters); !slices.Equal(words, expWords) {
		t.Errorf("%q: Step word boundaries %v, FirstWord %v", str, words, expWords)
	}

	// Sentences.
	var expSentences []int
	var sbState SentenceBreakState
	pos = 0
	for rest := str; len(rest) > 0; {
		var s string
		s, rest, sbState = FirstSentenceInString(rest, sbState)
		pos += len(s)
		expSentences = append(expSentences, pos)
	}
	if expSentences = keepAt(expSentences, clusters); !slices.Equal(sentences, expSentences) {
		t.Errorf("%q: Step sentence boundaries %v, FirstSentence %v", str, sentences, expSentences)
	}

	// Line breaks.
	var expLines, expMustLines []int
	var lbState LineBreakState
	pos = 0
	for rest := str; len(rest) > 0; {
		var (
			l    string
			must bool
		)
		l, rest, must, lbState = FirstLineSegmentInString(rest, lbState)
		pos += len(l)
		expLines = append(expLines, pos)
		if must {
			expMustLines = append(expMustLines, pos)
		}
	}
	if expLines = keepAt(expLines, clusters); !slices.Equal(lines, expLines) {
		t.Errorf("%q: Step line breaks %v, FirstLineSegment %v", str, lines, expLines)
	}
	if expMustLines = keepAt(expMustLines, clusters); !slices.Equal(mustLines, expMustLines) {
		t.Errorf("%q: Step mandatory line breaks %v, FirstLineSegment %v", str, mustLines, expMustLines)
	}
}

// TestStepConsistency verifies that the combined [Step] parser agrees with the
// standalone segmentation functions.
func TestStepConsistency(t *testing.T) {
	var inputs []string
	for _, cases := range [][]testCase{lineBreakTestCases, wordBreakTestCases, sentenceBreakTestCases} {
		for _, c := range cases {
			inputs = append(inputs, c.original)
		}
	}
	for _, c := range widthTestCases {
		inputs = append(inputs, c.original)
	}
	inputs = append(inputs, stepConsistencyExtraInputs...)
	for _, str := range inputs {
		checkStepConsistency(t, str)
	}
}

// TestStepStateRoundTrip verifies that every state component survives being
// packed into a [State].
func TestStepStateRoundTrip(t *testing.T) {
	const (
		gr = grRIEven | grGB9cLinker
		wb = wbEvenRI | wbZWJBit
		sb = sbSB8aSp | sbSB3
		lb = lbVI | lbZWJBit | lbCPeaFWHBit | lbDottedCircleBit | lb15Bit | lbQUPfBit |
			lbLB20aBit | lbPrevLB20aBit | lbHLHyphenBit | lbExtPicCnBit
	)
	s := newState(gr, wb, sb, lb)
	if s <= 0 || int64(s) > math.MaxInt32 {
		t.Fatalf("packed state %#x does not fit in a positive 32-bit int", s)
	}
	gotGr, gotWb, gotSb, gotLb := s.unpack()
	if gotGr != gr || gotWb != wb || gotSb != sb || gotLb != lb {
		t.Errorf("unpack(newState(%#x, %#x, %#x, %#x)) = %#x, %#x, %#x, %#x", gr, wb, sb, lb, gotGr, gotWb, gotSb, gotLb)
	}
}
