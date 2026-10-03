package uniseg

import (
	"runtime"
	"slices"
	"testing"
)

// Test all official Unicode test cases for word boundaries using the byte slice
// function.
func TestWordCasesBytes(t *testing.T) {
	for testNum, testCase := range wordBreakTestCases {
		/*t.Logf(`Test case %d %q: Expecting %x, getting %x, code points %x"`,
		testNum,
		strings.TrimSpace(testCase.original),
		testCase.expected,
		decomposed(testCase.original),
		[]rune(testCase.original))*/
		var (
			word  []byte
			index int
		)
		var state WordBreakState
		b := []byte(testCase.original)
	WordLoop:
		for index = 0; len(b) > 0; index++ {
			if index >= len(testCase.expected) {
				t.Errorf(`Test case %d %q failed: More words %d returned than expected %d`,
					testNum,
					testCase.original,
					index,
					len(testCase.expected))
				break
			}
			word, b, state = FirstWord(b, state)
			cluster := []rune(string(word))
			if len(cluster) != len(testCase.expected[index]) {
				t.Errorf(`Test case %d %q failed: Word at index %d has %d codepoints %x, %d expected %x`,
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
					t.Errorf(`Test case %d %q failed: Word at index %d is %x, expected %x`,
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
			t.Errorf(`Test case %d %q failed: Fewer words returned (%d) than expected (%d)`,
				testNum,
				testCase.original,
				index,
				len(testCase.expected))
		}
	}
}

// Test all official Unicode test cases for word boundaries using the string
// function.
func TestWordCasesString(t *testing.T) {
	for testNum, testCase := range wordBreakTestCases {
		/*t.Logf(`Test case %d %q: Expecting %x, getting %x, code points %x"`,
		testNum,
		strings.TrimSpace(testCase.original),
		testCase.expected,
		decomposed(testCase.original),
		[]rune(testCase.original))*/
		var (
			word  string
			index int
		)
		var state WordBreakState
		str := testCase.original
	WordLoop:
		for index = 0; len(str) > 0; index++ {
			if index >= len(testCase.expected) {
				t.Errorf(`Test case %d %q failed: More words %d returned than expected %d`,
					testNum,
					testCase.original,
					index,
					len(testCase.expected))
				break
			}
			word, str, state = FirstWordInString(str, state)
			cluster := []rune(string(word))
			if len(cluster) != len(testCase.expected[index]) {
				t.Errorf(`Test case %d %q failed: Word at index %d has %d codepoints %x, %d expected %x`,
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
					t.Errorf(`Test case %d %q failed: Word at index %d is %x, expected %x`,
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
			t.Errorf(`Test case %d %q failed: Fewer words returned (%d) than expected (%d)`,
				testNum,
				testCase.original,
				index,
				len(testCase.expected))
		}
	}
}

// Test all official Unicode test cases for word boundaries using the Words
// iterator.
func TestWords(t *testing.T) {
	for testNum, testCase := range wordBreakTestCases {
		index := 0
		offset := 0
		for i, w := range Words([]byte(testCase.original)) {
			if i != offset {
				t.Errorf(`Test case %d %q failed: Word at index %d starts at %d, expected %d`,
					testNum,
					testCase.original,
					index,
					i,
					offset)
				break
			}
			if index >= len(testCase.expected) {
				t.Errorf(`Test case %d %q failed: More words returned than expected %d`,
					testNum,
					testCase.original,
					len(testCase.expected))
				break
			}
			if word := string(w); word != string(testCase.expected[index]) {
				t.Errorf(`Test case %d %q failed: Word at index %d is %x, expected %x`,
					testNum,
					testCase.original,
					index,
					[]rune(word),
					testCase.expected[index])
				break
			}
			offset += len(w)
			index++
		}
		if index < len(testCase.expected) {
			t.Errorf(`Test case %d %q failed: Fewer words returned (%d) than expected (%d)`,
				testNum,
				testCase.original,
				index,
				len(testCase.expected))
		}
	}
}

// Test all official Unicode test cases for word boundaries using the
// WordsInString iterator.
func TestWordsInString(t *testing.T) {
	for testNum, testCase := range wordBreakTestCases {
		index := 0
		offset := 0
		for i, w := range WordsInString(testCase.original) {
			if i != offset {
				t.Errorf(`Test case %d %q failed: Word at index %d starts at %d, expected %d`,
					testNum,
					testCase.original,
					index,
					i,
					offset)
				break
			}
			if index >= len(testCase.expected) {
				t.Errorf(`Test case %d %q failed: More words returned than expected %d`,
					testNum,
					testCase.original,
					len(testCase.expected))
				break
			}
			if w != string(testCase.expected[index]) {
				t.Errorf(`Test case %d %q failed: Word at index %d is %x, expected %x`,
					testNum,
					testCase.original,
					index,
					[]rune(w),
					testCase.expected[index])
				break
			}
			offset += len(w)
			index++
		}
		if index < len(testCase.expected) {
			t.Errorf(`Test case %d %q failed: Fewer words returned (%d) than expected (%d)`,
				testNum,
				testCase.original,
				index,
				len(testCase.expected))
		}
	}
}

// Test that the Words iterators stop when the loop body breaks.
func TestWordsEarlyBreak(t *testing.T) {
	const input = "Hello, world! 🇩🇪"
	expected := []string{"Hello", ","}

	var gotBytes []string
	for _, w := range Words([]byte(input)) {
		gotBytes = append(gotBytes, string(w))
		if len(gotBytes) == 2 {
			break
		}
	}
	if !slices.Equal(gotBytes, expected) {
		t.Errorf("Words: got %q, expected %q", gotBytes, expected)
	}

	var gotString []string
	for _, w := range WordsInString(input) {
		gotString = append(gotString, w)
		if len(gotString) == 2 {
			break
		}
	}
	if !slices.Equal(gotString, expected) {
		t.Errorf("WordsInString: got %q, expected %q", gotString, expected)
	}
}

// Test that the Words iterators can be used more than once, even after an
// earlier iteration completed or stopped early.
func TestWordsReuse(t *testing.T) {
	const input = "Hello, world!"
	expected := []string{"Hello", ",", " ", "world", "!"}
	expectedIdx := []int{0, 5, 6, 7, 12}

	seqBytes := Words([]byte(input))
	seqString := WordsInString(input)

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
		for i, w := range seqBytes {
			gotBytes = append(gotBytes, string(w))
			gotBytesIdx = append(gotBytesIdx, i)
		}
		if !slices.Equal(gotBytes, expected) || !slices.Equal(gotBytesIdx, expectedIdx) {
			t.Errorf("Words: got %q at %v, expected %q at %v", gotBytes, gotBytesIdx, expected, expectedIdx)
		}

		var gotString []string
		var gotStringIdx []int
		for i, w := range seqString {
			gotString = append(gotString, w)
			gotStringIdx = append(gotStringIdx, i)
		}
		if !slices.Equal(gotString, expected) || !slices.Equal(gotStringIdx, expectedIdx) {
			t.Errorf("WordsInString: got %q at %v, expected %q at %v", gotString, gotStringIdx, expected, expectedIdx)
		}
	}
}

// Test that the Parser methods produce the same results as the package-level
// functions.
func TestParserWords(t *testing.T) {
	parsers := []*Parser{
		{},
		{EastAsianWidth: true},
		{EastAsianWidth: true, WideEmoji: true},
	}
	for _, p := range parsers {
		for testNum, testCase := range wordBreakTestCases {
			var expected []string
			for _, w := range WordsInString(testCase.original) {
				expected = append(expected, w)
			}

			var gotBytes []string
			for _, w := range p.Words([]byte(testCase.original)) {
				gotBytes = append(gotBytes, string(w))
			}
			if !slices.Equal(gotBytes, expected) {
				t.Errorf(`Test case %d %q failed with parser %+v: Words returned %q, expected %q`,
					testNum, testCase.original, *p, gotBytes, expected)
			}

			var gotString []string
			for _, w := range p.WordsInString(testCase.original) {
				gotString = append(gotString, w)
			}
			if !slices.Equal(gotString, expected) {
				t.Errorf(`Test case %d %q failed with parser %+v: WordsInString returned %q, expected %q`,
					testNum, testCase.original, *p, gotString, expected)
			}
		}
	}
}

// Benchmark the use of the word break function for byte slices.
func BenchmarkWordFunctionBytes(b *testing.B) {
	input := []byte(benchmarkStr)
	for b.Loop() {
		var c []byte
		var state WordBreakState
		str := input
		for len(str) > 0 {
			c, str, state = FirstWord(str, state)

			// to avoid the compiler optimizing out the benchmark
			runtime.KeepAlive(c)
			runtime.KeepAlive(str)
			runtime.KeepAlive(state)
		}
	}
}

// Benchmark the use of the word break function for strings.
func BenchmarkWordFunctionString(b *testing.B) {
	input := benchmarkStr
	for b.Loop() {
		var c string
		var state WordBreakState
		str := input
		for len(str) > 0 {
			c, str, state = FirstWordInString(str, state)

			// to avoid the compiler optimizing out the benchmark
			runtime.KeepAlive(c)
			runtime.KeepAlive(str)
			runtime.KeepAlive(state)
		}
	}
}

// Benchmark the use of the word break function for non-ASCII strings.
func BenchmarkWordFunctionMultilingual(b *testing.B) {
	input := benchmarkMultilingualStr
	for b.Loop() {
		var c string
		var state WordBreakState
		str := input
		for len(str) > 0 {
			c, str, state = FirstWordInString(str, state)

			// to avoid the compiler optimizing out the benchmark
			runtime.KeepAlive(c)
			runtime.KeepAlive(str)
			runtime.KeepAlive(state)
		}
	}
}

func FuzzFirstWordInString(f *testing.F) {
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
		var state WordBreakState
		var b []byte
		str := input
		for len(str) > 0 {
			var word string
			word, str, state = FirstWordInString(str, state)
			b = append(b, word...)
		}

		// Check if the constructed string is the same as the original.
		if string(b) != input {
			t.Errorf("Fuzzing failed: %q != %q", string(b), input)
		}
	})
}
