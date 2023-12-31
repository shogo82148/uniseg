package uniseg

import (
	"runtime"
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

// Benchmark the use of the word break function for byte slices.
func BenchmarkWordFunctionBytes(b *testing.B) {
	input := []byte(benchmarkStr)
	for i := 0; i < b.N; i++ {
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
	for i := 0; i < b.N; i++ {
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
