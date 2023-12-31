package uniseg

import (
	"runtime"
	"testing"
	"unicode/utf8"
)

const benchmarkStr = "This is 🏳️\u200d🌈, a test string ツ for grapheme cluster testing. 🏋🏽\u200d♀️🙂🙂 It's only relevant for benchmark tests."

type testCase = struct {
	name     string
	original string
	expected [][]rune
}

// The test cases for the simple test function.
var testCases = []testCase{
	{original: "", expected: [][]rune{}},
	{original: "x", expected: [][]rune{{0x78}}},
	{original: "basic", expected: [][]rune{{0x62}, {0x61}, {0x73}, {0x69}, {0x63}}},
	{original: "möp", expected: [][]rune{{0x6d}, {0x6f, 0x308}, {0x70}}},
	{original: "\r\n", expected: [][]rune{{0xd, 0xa}}},
	{original: "\n\n", expected: [][]rune{{0xa}, {0xa}}},
	{original: "\t*", expected: [][]rune{{0x9}, {0x2a}}},
	{original: "뢴", expected: [][]rune{{0x1105, 0x116c, 0x11ab}}},
	{original: "ܐ\u070fܒܓܕ", expected: [][]rune{{0x710}, {0x70f, 0x712}, {0x713}, {0x715}}},
	{original: "ำ", expected: [][]rune{{0xe33}}},
	{original: "ำำ", expected: [][]rune{{0xe33, 0xe33}}},
	{original: "สระอำ", expected: [][]rune{{0xe2a}, {0xe23}, {0xe30}, {0xe2d, 0xe33}}},
	{original: "*뢴*", expected: [][]rune{{0x2a}, {0x1105, 0x116c, 0x11ab}, {0x2a}}},
	{original: "*👩\u200d❤️\u200d💋\u200d👩*", expected: [][]rune{{0x2a}, {0x1f469, 0x200d, 0x2764, 0xfe0f, 0x200d, 0x1f48b, 0x200d, 0x1f469}, {0x2a}}},
	{original: "👩\u200d❤️\u200d💋\u200d👩", expected: [][]rune{{0x1f469, 0x200d, 0x2764, 0xfe0f, 0x200d, 0x1f48b, 0x200d, 0x1f469}}},
	{original: "🏋🏽‍♀️", expected: [][]rune{{0x1f3cb, 0x1f3fd, 0x200d, 0x2640, 0xfe0f}}},
	{original: "🙂", expected: [][]rune{{0x1f642}}},
	{original: "🙂🙂", expected: [][]rune{{0x1f642}, {0x1f642}}},
	{original: "🇩🇪", expected: [][]rune{{0x1f1e9, 0x1f1ea}}},
	{original: "🏳️\u200d🌈", expected: [][]rune{{0x1f3f3, 0xfe0f, 0x200d, 0x1f308}}},
	{original: "\t🏳️\u200d🌈", expected: [][]rune{{0x9}, {0x1f3f3, 0xfe0f, 0x200d, 0x1f308}}},
	{original: "\t🏳️\u200d🌈\t", expected: [][]rune{{0x9}, {0x1f3f3, 0xfe0f, 0x200d, 0x1f308}, {0x9}}},
	{original: "\r\n\uFE0E", expected: [][]rune{{13, 10}, {0xfe0e}}},
}

// Run all lists of test cases using the Graphemes class.
func TestGraphemesClass(t *testing.T) {
	allCases := append(testCases, graphemeBreakTestCases...)
	for testNum, testCase := range allCases {
		/*t.Logf(`Test case %d %q: Expecting %x, getting %x, code points %x"`,
		testNum,
		strings.TrimSpace(testCase.original),
		testCase.expected,
		decomposed(testCase.original),
		[]rune(testCase.original))*/
		gr := NewGraphemes(testCase.original)
		var index int
	GraphemeLoop:
		for index = 0; gr.Next(); index++ {
			if index >= len(testCase.expected) {
				t.Errorf(`Test case %d %q %q failed: More grapheme clusters returned than expected %d`,
					testNum,
					testCase.name,
					testCase.original,
					len(testCase.expected))
				break
			}
			cluster := gr.Runes()
			if len(cluster) != len(testCase.expected[index]) {
				t.Errorf(`Test case %d %q %q failed: Grapheme cluster at index %d has %d codepoints %x, %d expected %x`,
					testNum,
					testCase.name,
					testCase.original,
					index,
					len(cluster),
					cluster,
					len(testCase.expected[index]),
					testCase.expected[index])
				break
			}
			for i, r := range cluster {
				if r != testCase.expected[index][i] {
					t.Errorf(`Test case %d %q %q failed: Grapheme cluster at index %d is %x, expected %x`,
						testNum,
						testCase.name,
						testCase.original,
						index,
						cluster,
						testCase.expected[index])
					break GraphemeLoop
				}
			}
		}
		if index < len(testCase.expected) {
			t.Errorf(`Test case %d %q %q failed: Fewer grapheme clusters returned (%d) than expected (%d)`,
				testNum,
				testCase.name,
				testCase.original,
				index,
				len(testCase.expected))
		}
	}
}

// Run the standard Unicode test cases for word boundaries using the Graphemes
// class.
func TestGraphemesClassWord(t *testing.T) {
	for testNum, testCase := range wordBreakTestCases {
		if testNum == 1703 {
			// This test case reveals an inconsistency in the Unicode rule set,
			// namely the handling of ZWJ within two RI graphemes. (Grapheme
			// rules will restart the RI count, word rules will ignore the ZWJ.)
			// An error has been reported.
			continue
		}
		/*t.Logf(`Test case %d %q: Expecting %x, getting %x, code points %x"`,
		testNum,
		strings.TrimSpace(testCase.original),
		testCase.expected,
		decomposed(testCase.original),
		[]rune(testCase.original))*/
		gr := NewGraphemes(testCase.original)
		var (
			index   int
			cluster []rune
		)
	GraphemeLoop:
		for gr.Next() {
			if index >= len(testCase.expected) {
				t.Errorf(`Test case %d %q %q failed: More words returned than expected %d`,
					testNum,
					testCase.name,
					testCase.original,
					len(testCase.expected))
				break
			}
			cluster = append(cluster, gr.Runes()...)
			if gr.IsWordBoundary() {
				if len(cluster) != len(testCase.expected[index]) {
					t.Errorf(`Test case %d %q %q failed: Word at index %d has %d codepoints %x, %d expected %x`,
						testNum,
						testCase.name,
						testCase.original,
						index,
						len(cluster),
						cluster,
						len(testCase.expected[index]),
						testCase.expected[index])
					break
				}
				for i, r := range cluster {
					if r != testCase.expected[index][i] {
						t.Errorf(`Test case %d %q %q failed: Word at index %d is %x, expected %x`,
							testNum,
							testCase.name,
							testCase.original,
							index,
							cluster,
							testCase.expected[index])
						break GraphemeLoop
					}
				}
				cluster = nil
				index++
			}
		}
		if index < len(testCase.expected) {
			t.Errorf(`Test case %d %q %q failed: Fewer words returned (%d) than expected (%d)`,
				testNum,
				testCase.name,
				testCase.original,
				index,
				len(testCase.expected))
		}
	}
}

// Run the standard Unicode test cases for sentence boundaries using the
// Graphemes class.
func TestGraphemesClassSentence(t *testing.T) {
	for testNum, testCase := range sentenceBreakTestCases {
		/*t.Logf(`Test case %d %q: Expecting %x, getting %x, code points %x"`,
		testNum,
		strings.TrimSpace(testCase.original),
		testCase.expected,
		decomposed(testCase.original),
		[]rune(testCase.original))*/
		gr := NewGraphemes(testCase.original)
		var (
			index   int
			cluster []rune
		)
	GraphemeLoop:
		for gr.Next() {
			if index >= len(testCase.expected) {
				t.Errorf(`Test case %d %q failed: More sentences returned than expected %d`,
					testNum,
					testCase.original,
					len(testCase.expected))
				break
			}
			cluster = append(cluster, gr.Runes()...)
			if gr.IsSentenceBoundary() {
				if len(cluster) != len(testCase.expected[index]) {
					t.Errorf(`Test case %d %q failed: Sentence at index %d has %d codepoints %x, %d expected %x`,
						testNum,
						testCase.original,
						index,
						len(cluster),
						cluster,
						len(testCase.expected[index]),
						testCase.expected[index])
					break
				}
				for i, r := range cluster {
					if r != testCase.expected[index][i] {
						t.Errorf(`Test case %d %q failed: Sentence at index %d is %x, expected %x`,
							testNum,
							testCase.original,
							index,
							cluster,
							testCase.expected[index])
						break GraphemeLoop
					}
				}
				cluster = nil
				index++
			}
		}
		if index < len(testCase.expected) {
			t.Errorf(`Test case %d %q failed: Fewer sentences returned (%d) than expected (%d)`,
				testNum,
				testCase.original,
				index,
				len(testCase.expected))
		}
	}
}

// Test the Str() function.
func TestGraphemesStr(t *testing.T) {
	gr := NewGraphemes("möp")
	gr.Next()
	gr.Next()
	gr.Next()
	if str := gr.Str(); str != "p" {
		t.Errorf(`Expected "p", got %q`, str)
	}
}

// Test the Bytes() function.
func TestGraphemesBytes(t *testing.T) {
	gr := NewGraphemes("A👩\u200d❤️\u200d💋\u200d👩B")
	gr.Next()
	gr.Next()
	gr.Next()
	b := gr.Bytes()
	if len(b) != 1 {
		t.Fatalf(`Expected len("B") == 1, got %d`, len(b))
	}
	if b[0] != 'B' {
		t.Errorf(`Expected "B", got %q`, string(b[0]))
	}
}

// Test the Positions() function.
func TestGraphemesPositions(t *testing.T) {
	gr := NewGraphemes("A👩\u200d❤️\u200d💋\u200d👩B")
	gr.Next()
	gr.Next()
	from, to := gr.Positions()
	if from != 1 || to != 28 {
		t.Errorf(`Expected from=%d to=%d, got from=%d to=%d`, 1, 28, from, to)
	}
}

// Test the Reset() function.
func TestGraphemesReset(t *testing.T) {
	gr := NewGraphemes("möp")
	gr.Next()
	gr.Next()
	gr.Next()
	gr.Reset()
	gr.Next()
	if str := gr.Str(); str != "m" {
		t.Errorf(`Expected "m", got %q`, str)
	}
}

// Test retrieving clusters before calling Next().
func TestGraphemesEarly(t *testing.T) {
	gr := NewGraphemes("test")
	r := gr.Runes()
	if r != nil {
		t.Errorf(`Expected nil rune slice, got %x`, r)
	}
	str := gr.Str()
	if str != "" {
		t.Errorf(`Expected empty string, got %q`, str)
	}
	b := gr.Bytes()
	if b != nil {
		t.Errorf(`Expected byte rune slice, got %x`, b)
	}
	from, to := gr.Positions()
	if from != 0 || to != 0 {
		t.Errorf(`Expected from=%d to=%d, got from=%d to=%d`, 0, 0, from, to)
	}
}

// Test retrieving more clusters after retrieving the last cluster.
func TestGraphemesLate(t *testing.T) {
	gr := NewGraphemes("x")
	gr.Next()
	gr.Next()
	r := gr.Runes()
	if r != nil {
		t.Errorf(`Expected nil rune slice, got %x`, r)
	}
	str := gr.Str()
	if str != "" {
		t.Errorf(`Expected empty string, got %q`, str)
	}
	b := gr.Bytes()
	if b != nil {
		t.Errorf(`Expected byte rune slice, got %x`, b)
	}
	from, to := gr.Positions()
	if from != 1 || to != 1 {
		t.Errorf(`Expected from=%d to=%d, got from=%d to=%d`, 1, 1, from, to)
	}
}

// Test the GraphemeClusterCount function.
func TestGraphemesCount(t *testing.T) {
	if n := GraphemeClusterCount("🇩🇪🏳️\u200d🌈"); n != 2 {
		t.Errorf(`Expected 2 grapheme clusters, got %d`, n)
	}
}

// Test the ReverseString function.
func TestReverseString(t *testing.T) {
	for _, testCase := range testCases {
		var r []rune
		for index := len(testCase.expected) - 1; index >= 0; index-- {
			r = append(r, testCase.expected[index]...)
		}
		if string(r) != ReverseString(testCase.original) {
			t.Errorf(`Expected reverse of %q to be %q, got %q`, testCase.original, string(r), ReverseString(testCase.original))
		}
	}

	// Three additional ones, for good measure.
	if ReverseString("🇩🇪🏳️\u200d🌈") != "🏳️\u200d🌈🇩🇪" {
		t.Error("Flags weren't reversed correctly")
	}
	if ReverseString("🏳️\u200d🌈") != "🏳️\u200d🌈" {
		t.Error("Flag wasn't reversed correctly")
	}
	if ReverseString("") != "" {
		t.Error("Empty string wasn't reversed correctly")
	}
}

// Run all lists of test cases using the Graphemes function for byte slices.
func TestGraphemesFunctionBytes(t *testing.T) {
	allCases := append(testCases, graphemeBreakTestCases...)
	for testNum, testCase := range allCases {
		/*t.Logf(`Test case %d %q: Expecting %x, getting %x, code points %x"`,
		testNum,
		strings.TrimSpace(testCase.original),
		testCase.expected,
		decomposed(testCase.original),
		[]rune(testCase.original))*/
		b := []byte(testCase.original)
		var (
			index int
			c     []byte
			state GraphemeBreakState
		)
	GraphemeLoop:
		for len(b) > 0 {
			c, b, _, state = FirstGraphemeCluster(b, state)

			if index >= len(testCase.expected) {
				t.Errorf(`Test case %d %q failed: More grapheme clusters returned than expected %d`,
					testNum,
					testCase.original,
					len(testCase.expected))
				break
			}

			cluster := []rune(string(c))
			if len(cluster) != len(testCase.expected[index]) {
				t.Errorf(`Test case %d %q failed: Grapheme cluster at index %d has %d codepoints %x, %d expected %x`,
					testNum,
					testCase.original,
					index,
					len(cluster),
					cluster,
					len(testCase.expected[index]),
					testCase.expected[index])
				break
			}
			for i, r := range cluster {
				if r != testCase.expected[index][i] {
					t.Errorf(`Test case %d %q failed: Grapheme cluster at index %d is %x, expected %x`,
						testNum,
						testCase.original,
						index,
						cluster,
						testCase.expected[index])
					break GraphemeLoop
				}
			}

			index++
		}
		if index < len(testCase.expected) {
			t.Errorf(`Test case %d %q failed: Fewer grapheme clusters returned (%d) than expected (%d)`,
				testNum,
				testCase.original,
				index,
				len(testCase.expected))
		}
	}
}

// Run all lists of test cases using the Graphemes function for strings.
func TestGraphemesFunctionString(t *testing.T) {
	allCases := append(testCases, graphemeBreakTestCases...)
	for testNum, testCase := range allCases {
		/*t.Logf(`Test case %d %q: Expecting %x, getting %x, code points %x"`,
		testNum,
		strings.TrimSpace(testCase.original),
		testCase.expected,
		decomposed(testCase.original),
		[]rune(testCase.original))*/
		str := testCase.original
		var (
			index int
			c     string
			state GraphemeBreakState
		)
	GraphemeLoop:
		for len(str) > 0 {
			c, str, _, state = FirstGraphemeClusterInString(str, state)

			if index >= len(testCase.expected) {
				t.Errorf(`Test case %d %q failed: More grapheme clusters returned than expected %d`,
					testNum,
					testCase.original,
					len(testCase.expected))
				break
			}

			cluster := []rune(c)
			if len(cluster) != len(testCase.expected[index]) {
				t.Errorf(`Test case %d %q failed: Grapheme cluster at index %d has %d codepoints %x, %d expected %x`,
					testNum,
					testCase.original,
					index,
					len(cluster),
					cluster,
					len(testCase.expected[index]),
					testCase.expected[index])
				break
			}
			for i, r := range cluster {
				if r != testCase.expected[index][i] {
					t.Errorf(`Test case %d %q failed: Grapheme cluster at index %d is %x, expected %x`,
						testNum,
						testCase.original,
						index,
						cluster,
						testCase.expected[index])
					break GraphemeLoop
				}
			}

			index++
		}
		if index < len(testCase.expected) {
			t.Errorf(`Test case %d %q failed: Fewer grapheme clusters returned (%d) than expected (%d)`,
				testNum,
				testCase.original,
				index,
				len(testCase.expected))
		}
	}
}

// Benchmark the use of the Graphemes class.
func BenchmarkGraphemesClass(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := NewGraphemes(benchmarkStr)
		for g.Next() {
			runtime.KeepAlive(g.Runes())
		}
	}
}

// Benchmark the use of the Graphemes function for byte slices.
func BenchmarkGraphemesFunctionBytes(b *testing.B) {
	input := []byte(benchmarkStr)
	for i := 0; i < b.N; i++ {
		var c []byte
		var width int
		var state GraphemeBreakState
		str := input
		for len(str) > 0 {
			c, str, width, state = FirstGraphemeCluster(str, state)

			// to avoid the compiler optimizing out the benchmark
			runtime.KeepAlive(c)
			runtime.KeepAlive(str)
			runtime.KeepAlive(width)
			runtime.KeepAlive(state)
		}
	}
}

// Benchmark the use of the Graphemes function for strings.
func BenchmarkGraphemesFunctionString(b *testing.B) {
	input := benchmarkStr
	for i := 0; i < b.N; i++ {
		var c string
		var width int
		var state GraphemeBreakState
		str := input
		for len(str) > 0 {
			c, str, _, state = FirstGraphemeClusterInString(str, state)

			// to avoid the compiler optimizing out the benchmark
			runtime.KeepAlive(c)
			runtime.KeepAlive(str)
			runtime.KeepAlive(width)
			runtime.KeepAlive(state)
		}
	}
}

func FuzzGraphemeClusterCount(f *testing.F) {
	for _, test := range wordBreakTestCases {
		f.Add(test.original)
	}
	for _, test := range sentenceBreakTestCases {
		f.Add(test.original)
	}
	for _, test := range lineBreakTestCases {
		f.Add(test.original)
	}
	for _, test := range graphemeBreakTestCases {
		f.Add(test.original)
	}
	for _, test := range testCases {
		f.Add(test.original)
	}
	f.Fuzz(func(t *testing.T, input string) {
		count := GraphemeClusterCount(input)
		if count < 0 {
			t.Errorf("negative count: %d", count)
		}
	})
}

func FuzzReverseString(f *testing.F) {
	for _, test := range wordBreakTestCases {
		f.Add(test.original)
	}
	for _, test := range sentenceBreakTestCases {
		f.Add(test.original)
	}
	for _, test := range lineBreakTestCases {
		f.Add(test.original)
	}
	for _, test := range graphemeBreakTestCases {
		f.Add(test.original)
	}
	for _, test := range testCases {
		f.Add(test.original)
	}
	f.Fuzz(func(t *testing.T, input string) {
		reversed := ReverseString(input)
		if !utf8.ValidString(input) {
			return
		}
		if !utf8.ValidString(reversed) {
			t.Errorf("reversed string is not valid: %q", reversed)
		}
	})
}
