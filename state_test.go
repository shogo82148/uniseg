package uniseg

import (
	"math"
	"slices"
	"testing"
)

// negativeStates are initial states which must be treated like 0 ("unknown").
var negativeStates = []int{-1, math.MinInt}

// stateTestInputs returns the inputs of all conformance test cases.
func stateTestInputs() []string {
	var inputs []string
	for _, cases := range [][]testCase{graphemeBreakTestCases, wordBreakTestCases, sentenceBreakTestCases, lineBreakTestCases} {
		for _, c := range cases {
			inputs = append(inputs, c.original)
		}
	}
	return inputs
}

// segmentAll splits str by calling first repeatedly, starting with the given
// initial state.
func segmentAll[S ~int](str string, initial S, first func(string, S) (string, string, S)) []string {
	var segments []string
	state := initial
	for len(str) > 0 {
		var segment string
		segment, str, state = first(str, state)
		segments = append(segments, segment)
	}
	return segments
}

// segmentAllBytes is like segmentAll but for byte slices.
func segmentAllBytes[S ~int](str string, initial S, first func([]byte, S) ([]byte, []byte, S)) []string {
	var segments []string
	state := initial
	b := []byte(str)
	for len(b) > 0 {
		var segment []byte
		segment, b, state = first(b, state)
		segments = append(segments, string(segment))
	}
	return segments
}

// checkNegativeState verifies that segmenting with a negative initial state
// gives the same result as with 0, without panicking.
func checkNegativeState[S ~int](t *testing.T, name string, segment func(string, S) []string) {
	t.Helper()
	for _, str := range stateTestInputs() {
		expected := segment(str, 0)
		for _, neg := range negativeStates {
			var got []string
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("%s(%q, %d) panicked: %v", name, str, neg, r)
					}
				}()
				got = segment(str, S(neg))
			}()
			if got != nil && !slices.Equal(got, expected) {
				t.Errorf("%s(%q, %d) = %q, with state 0 %q", name, str, neg, got, expected)
			}
		}
	}
}

func TestNegativeStateGrapheme(t *testing.T) {
	checkNegativeState(t, "FirstGraphemeClusterInString", func(str string, s GraphemeBreakState) []string {
		return segmentAll(str, s, func(str string, s GraphemeBreakState) (string, string, GraphemeBreakState) {
			c, rest, _, s := FirstGraphemeClusterInString(str, s)
			return c, rest, s
		})
	})
	checkNegativeState(t, "FirstGraphemeCluster", func(str string, s GraphemeBreakState) []string {
		return segmentAllBytes(str, s, func(b []byte, s GraphemeBreakState) ([]byte, []byte, GraphemeBreakState) {
			c, rest, _, s := FirstGraphemeCluster(b, s)
			return c, rest, s
		})
	})
}

func TestNegativeStateWord(t *testing.T) {
	checkNegativeState(t, "FirstWordInString", func(str string, s WordBreakState) []string {
		return segmentAll(str, s, FirstWordInString)
	})
	checkNegativeState(t, "FirstWord", func(str string, s WordBreakState) []string {
		return segmentAllBytes(str, s, FirstWord)
	})
}

func TestNegativeStateSentence(t *testing.T) {
	checkNegativeState(t, "FirstSentenceInString", func(str string, s SentenceBreakState) []string {
		return segmentAll(str, s, FirstSentenceInString)
	})
	checkNegativeState(t, "FirstSentence", func(str string, s SentenceBreakState) []string {
		return segmentAllBytes(str, s, FirstSentence)
	})
}

func TestNegativeStateLine(t *testing.T) {
	checkNegativeState(t, "FirstLineSegmentInString", func(str string, s LineBreakState) []string {
		return segmentAll(str, s, func(str string, s LineBreakState) (string, string, LineBreakState) {
			segment, rest, _, s := FirstLineSegmentInString(str, s)
			return segment, rest, s
		})
	})
	checkNegativeState(t, "FirstLineSegment", func(str string, s LineBreakState) []string {
		return segmentAllBytes(str, s, func(b []byte, s LineBreakState) ([]byte, []byte, LineBreakState) {
			segment, rest, _, s := FirstLineSegment(b, s)
			return segment, rest, s
		})
	})
}

func TestNegativeStateStep(t *testing.T) {
	checkNegativeState(t, "StepString", func(str string, s State) []string {
		return segmentAll(str, s, func(str string, s State) (string, string, State) {
			c, rest, _, s := StepString(str, s)
			return c, rest, s
		})
	})
	checkNegativeState(t, "Step", func(str string, s State) []string {
		return segmentAllBytes(str, s, func(b []byte, s State) ([]byte, []byte, State) {
			c, rest, _, s := Step(b, s)
			return c, rest, s
		})
	})
}
