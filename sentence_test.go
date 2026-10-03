package uniseg

import (
	"runtime"
	"slices"
	"testing"
)

// Test all official Unicode test cases for sentence boundaries using the byte
// slice function.
func TestSentenceCasesBytes(t *testing.T) {
	for testNum, testCase := range sentenceBreakTestCases {
		/*t.Logf(`Test case %d %q: Expecting %x, getting %x, code points %x"`,
		testNum,
		strings.TrimSpace(testCase.original),
		testCase.expected,
		decomposed(testCase.original),
		[]rune(testCase.original))*/
		var (
			sentence []byte
			index    int
			state    SentenceBreakState
		)
		b := []byte(testCase.original)
	WordLoop:
		for index = 0; len(b) > 0; index++ {
			if index >= len(testCase.expected) {
				t.Errorf(`Test case %d %q failed: More sentences %d returned than expected %d`,
					testNum,
					testCase.original,
					index,
					len(testCase.expected))
				break
			}
			sentence, b, state = FirstSentence(b, state)
			cluster := []rune(string(sentence))
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
					break WordLoop
				}
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

// Test all official Unicode test cases for sentence boundaries using the string
// function.
func TestSentenceCasesString(t *testing.T) {
	for testNum, testCase := range sentenceBreakTestCases {
		/*t.Logf(`Test case %d %q: Expecting %x, getting %x, code points %x"`,
		testNum,
		strings.TrimSpace(testCase.original),
		testCase.expected,
		decomposed(testCase.original),
		[]rune(testCase.original))*/
		var (
			sentence string
			index    int
			state    SentenceBreakState
		)
		str := testCase.original
	WordLoop:
		for index = 0; len(str) > 0; index++ {
			if index >= len(testCase.expected) {
				t.Errorf(`Test case %d %q %q failed: More sentences %d returned than expected %d`,
					testNum,
					testCase.name,
					testCase.original,
					index,
					len(testCase.expected))
				break
			}
			sentence, str, state = FirstSentenceInString(str, state)
			cluster := []rune(string(sentence))
			if len(cluster) != len(testCase.expected[index]) {
				t.Errorf(`Test case %d %q %q failed: Sentence at index %d has %d codepoints %x, %d expected %x`,
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
					t.Errorf(`Test case %d %q %q failed: Sentence at index %d is %x, expected %x`,
						testNum,
						testCase.name,
						testCase.original,
						index,
						cluster,
						testCase.expected[index])
					break WordLoop
				}
			}
		}
		if index < len(testCase.expected) {
			t.Errorf(`Test case %d %q %q failed: Fewer sentences returned (%d) than expected (%d)`,
				testNum,
				testCase.name,
				testCase.original,
				index,
				len(testCase.expected))
		}
	}
}

// Test all official Unicode test cases for sentence boundaries using the
// Sentences iterator.
func TestSentences(t *testing.T) {
	for testNum, testCase := range sentenceBreakTestCases {
		index := 0
		offset := 0
		for i, s := range Sentences([]byte(testCase.original)) {
			if i != offset {
				t.Errorf(`Test case %d %q failed: Sentence at index %d starts at %d, expected %d`,
					testNum,
					testCase.original,
					index,
					i,
					offset)
				break
			}
			if index >= len(testCase.expected) {
				t.Errorf(`Test case %d %q failed: More sentences returned than expected %d`,
					testNum,
					testCase.original,
					len(testCase.expected))
				break
			}
			if sentence := string(s); sentence != string(testCase.expected[index]) {
				t.Errorf(`Test case %d %q failed: Sentence at index %d is %x, expected %x`,
					testNum,
					testCase.original,
					index,
					[]rune(sentence),
					testCase.expected[index])
				break
			}
			offset += len(s)
			index++
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

// Test all official Unicode test cases for sentence boundaries using the
// SentencesInString iterator.
func TestSentencesInString(t *testing.T) {
	for testNum, testCase := range sentenceBreakTestCases {
		index := 0
		offset := 0
		for i, s := range SentencesInString(testCase.original) {
			if i != offset {
				t.Errorf(`Test case %d %q failed: Sentence at index %d starts at %d, expected %d`,
					testNum,
					testCase.original,
					index,
					i,
					offset)
				break
			}
			if index >= len(testCase.expected) {
				t.Errorf(`Test case %d %q failed: More sentences returned than expected %d`,
					testNum,
					testCase.original,
					len(testCase.expected))
				break
			}
			if s != string(testCase.expected[index]) {
				t.Errorf(`Test case %d %q failed: Sentence at index %d is %x, expected %x`,
					testNum,
					testCase.original,
					index,
					[]rune(s),
					testCase.expected[index])
				break
			}
			offset += len(s)
			index++
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

// Test that the Sentences iterators stop when the loop body breaks.
func TestSentencesEarlyBreak(t *testing.T) {
	const input = "This is a test. Is it? Yes! It is."
	expected := []string{"This is a test. ", "Is it? "}

	var gotBytes []string
	for _, s := range Sentences([]byte(input)) {
		gotBytes = append(gotBytes, string(s))
		if len(gotBytes) == 2 {
			break
		}
	}
	if !slices.Equal(gotBytes, expected) {
		t.Errorf("Sentences: got %q, expected %q", gotBytes, expected)
	}

	var gotString []string
	for _, s := range SentencesInString(input) {
		gotString = append(gotString, s)
		if len(gotString) == 2 {
			break
		}
	}
	if !slices.Equal(gotString, expected) {
		t.Errorf("SentencesInString: got %q, expected %q", gotString, expected)
	}
}

// Test that the Sentences iterators can be used more than once, even after an
// earlier iteration completed or stopped early.
func TestSentencesReuse(t *testing.T) {
	const input = "This is a test. Is it? Yes! It is."
	expected := []string{"This is a test. ", "Is it? ", "Yes! ", "It is."}
	expectedIdx := []int{0, 16, 23, 28}

	seqBytes := Sentences([]byte(input))
	seqString := SentencesInString(input)

	// Stop early once, then run the iterators to completion twice.
	for range seqBytes {
		break
	}
	for range seqString {
		break
	}
	for range 2 {
		var gotBytes []string
		var gotBytesIdx []int
		for i, s := range seqBytes {
			gotBytes = append(gotBytes, string(s))
			gotBytesIdx = append(gotBytesIdx, i)
		}
		if !slices.Equal(gotBytes, expected) || !slices.Equal(gotBytesIdx, expectedIdx) {
			t.Errorf("Sentences: got %q at %v, expected %q at %v", gotBytes, gotBytesIdx, expected, expectedIdx)
		}

		var gotString []string
		var gotStringIdx []int
		for i, s := range seqString {
			gotString = append(gotString, s)
			gotStringIdx = append(gotStringIdx, i)
		}
		if !slices.Equal(gotString, expected) || !slices.Equal(gotStringIdx, expectedIdx) {
			t.Errorf("SentencesInString: got %q at %v, expected %q at %v", gotString, gotStringIdx, expected, expectedIdx)
		}
	}
}

// Test that the Parser methods produce the same results as the package-level
// functions.
func TestParserSentences(t *testing.T) {
	parsers := []*Parser{
		{},
		{EastAsianWidth: true},
		{EastAsianWidth: true, WideEmoji: true},
	}
	for _, p := range parsers {
		for testNum, testCase := range sentenceBreakTestCases {
			var expected []string
			for _, s := range SentencesInString(testCase.original) {
				expected = append(expected, s)
			}

			var gotBytes []string
			for _, s := range p.Sentences([]byte(testCase.original)) {
				gotBytes = append(gotBytes, string(s))
			}
			if !slices.Equal(gotBytes, expected) {
				t.Errorf(`Test case %d %q failed with parser %+v: Sentences returned %q, expected %q`,
					testNum, testCase.original, *p, gotBytes, expected)
			}

			var gotString []string
			for _, s := range p.SentencesInString(testCase.original) {
				gotString = append(gotString, s)
			}
			if !slices.Equal(gotString, expected) {
				t.Errorf(`Test case %d %q failed with parser %+v: SentencesInString returned %q, expected %q`,
					testNum, testCase.original, *p, gotString, expected)
			}
		}
	}
}

// Benchmark the use of the sentence break function for byte slices.
func BenchmarkSentenceFunctionBytes(b *testing.B) {
	input := []byte(benchmarkStr)
	for b.Loop() {
		var c []byte
		var state SentenceBreakState
		str := input
		for len(str) > 0 {
			c, str, state = FirstSentence(str, state)

			// to avoid the compiler optimizing out the benchmark
			runtime.KeepAlive(c)
			runtime.KeepAlive(str)
			runtime.KeepAlive(state)
		}
	}
}

// Benchmark the use of the sentence break function for strings.
func BenchmarkSentenceFunctionString(b *testing.B) {
	input := benchmarkStr
	for b.Loop() {
		var c string
		var state SentenceBreakState
		str := input
		for len(str) > 0 {
			c, str, state = FirstSentenceInString(str, state)

			// to avoid the compiler optimizing out the benchmark
			runtime.KeepAlive(c)
			runtime.KeepAlive(str)
			runtime.KeepAlive(state)
		}
	}
}

func FuzzFirstSentenceInString(f *testing.F) {
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
		var state SentenceBreakState
		var b []byte
		str := input
		for len(str) > 0 {
			var sentence string
			sentence, str, state = FirstSentenceInString(str, state)
			b = append(b, sentence...)
		}

		// Check if the constructed string is the same as the original.
		if string(b) != input {
			t.Errorf("Fuzzing failed: %q != %q", string(b), input)
		}
	})
}
