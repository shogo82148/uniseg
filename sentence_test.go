package uniseg

import (
	"runtime"
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

// Benchmark the use of the sentence break function for byte slices.
func BenchmarkSentenceFunctionBytes(b *testing.B) {
	input := []byte(benchmarkStr)
	for i := 0; i < b.N; i++ {
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
	for i := 0; i < b.N; i++ {
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
