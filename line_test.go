package uniseg

import (
	"runtime"
	"testing"
)

// Test all official Unicode test cases for line breaks using the byte slice
// function.
func TestLineCasesBytes(t *testing.T) {
	for testNum, testCase := range lineBreakTestCases {
		/*t.Logf(`Test case %d %q: Expecting %x, getting %x, code points %x"`,
		testNum,
		strings.TrimSpace(testCase.original),
		testCase.expected,
		decomposed(testCase.original),
		[]rune(testCase.original))*/
		var (
			segment []byte
			index   int
			state   LineBreakState
		)
		b := []byte(testCase.original)
	WordLoop:
		for index = 0; len(b) > 0; index++ {
			if index >= len(testCase.expected) {
				t.Errorf(`Test case %d %q: %q failed: More segments %d returned than expected %d`,
					testNum,
					testCase.name,
					testCase.original,
					index,
					len(testCase.expected))
				break
			}
			segment, b, _, state = FirstLineSegment(b, state)
			cluster := []rune(string(segment))
			if len(cluster) != len(testCase.expected[index]) {
				t.Errorf(`Test case %d %q: %q failed: Segment at index %d has %d codepoints %x, %d expected %x`,
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
					t.Errorf(`Test case %d %q: %q failed: Segment at index %d is %x, expected %x`,
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
			t.Errorf(`Test case %d %q failed: Fewer segments returned (%d) than expected (%d)`,
				testNum,
				testCase.original,
				index,
				len(testCase.expected))
		}
	}
}

// Test all official Unicode test cases for line breaks using the string
// function.
func TestLineCasesString(t *testing.T) {
	for testNum, testCase := range lineBreakTestCases {
		/*t.Logf(`Test case %d %q: Expecting %x, getting %x, code points %x"`,
		testNum,
		strings.TrimSpace(testCase.original),
		testCase.expected,
		decomposed(testCase.original),
		[]rune(testCase.original))*/
		var (
			segment string
			index   int
			state   LineBreakState
		)
		str := testCase.original
	WordLoop:
		for index = 0; len(str) > 0; index++ {
			if index >= len(testCase.expected) {
				t.Errorf(`Test case %d %q: %q failed: More segments %d returned than expected %d`,
					testNum,
					testCase.name,
					testCase.original,
					index,
					len(testCase.expected))
				break
			}
			segment, str, _, state = FirstLineSegmentInString(str, state)
			cluster := []rune(string(segment))
			if len(cluster) != len(testCase.expected[index]) {
				t.Errorf(`Test case %d %q: %q failed: Segment at index %d has %d codepoints %x, %d expected %x`,
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
					t.Errorf(`Test case %d %q: %q failed: Segment at index %d is %x, expected %x`,
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
			t.Errorf(`Test case %d %q: %q failed: Fewer segments returned (%d) than expected (%d)`,
				testNum,
				testCase.name,
				testCase.original,
				index,
				len(testCase.expected))
		}
	}
}

var hasTrailingLineBreakTestCases = []struct {
	input string
	want  bool
}{
	{"\v", true},     // prBK
	{"\r", true},     // prCR
	{"\n", true},     // prLF
	{"\u0085", true}, // prNL
	{" ", false},
	{"A", false},
	{"", false},
}

func TestHasTrailingLineBreak(t *testing.T) {
	for _, tt := range hasTrailingLineBreakTestCases {
		got := HasTrailingLineBreak([]byte(tt.input))
		if got != tt.want {
			t.Errorf("HasTrailingLineBreak(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestHasTrailingLineBreakInString(t *testing.T) {
	for _, tt := range hasTrailingLineBreakTestCases {
		got := HasTrailingLineBreakInString(tt.input)
		if got != tt.want {
			t.Errorf("HasTrailingLineBreak(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func collectLineSegmentsInString(input string) []string {
	var (
		segments []string
		state    LineBreakState
	)
	for len(input) > 0 {
		var segment string
		segment, input, _, state = FirstLineSegmentInString(input, state)
		segments = append(segments, segment)
	}
	return segments
}

func collectLineSegmentsInBytes(input []byte) []string {
	var (
		segments []string
		state    LineBreakState
	)
	for len(input) > 0 {
		var segment []byte
		segment, input, _, state = FirstLineSegment(input, state)
		segments = append(segments, string(segment))
	}
	return segments
}

func TestLB20aWordInitialHyphenContext(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "AL HY AL allows break",
			input:    "foo-bar",
			expected: []string{"foo-", "bar"},
		},
		{
			name:     "AL BAHyphen AL allows break",
			input:    "foo\u2010bar",
			expected: []string{"foo\u2010", "bar"},
		},
		{
			name:     "sot HY AL stays unbroken",
			input:    "-bar",
			expected: []string{"-bar"},
		},
		{
			name:     "SP HY AL stays unbroken",
			input:    " -bar",
			expected: []string{" ", "-bar"},
		},
		{
			name:     "SP BAHyphen AL stays unbroken",
			input:    " \u2010bar",
			expected: []string{" ", "\u2010bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/string", func(t *testing.T) {
			got := collectLineSegmentsInString(tt.input)
			if len(got) != len(tt.expected) {
				t.Fatalf("collectLineSegmentsInString(%q) returned %d segments, want %d (%q)", tt.input, len(got), len(tt.expected), got)
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Fatalf("collectLineSegmentsInString(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.expected[i])
				}
			}
		})

		t.Run(tt.name+"/bytes", func(t *testing.T) {
			got := collectLineSegmentsInBytes([]byte(tt.input))
			if len(got) != len(tt.expected) {
				t.Fatalf("collectLineSegmentsInBytes(%q) returned %d segments, want %d (%q)", tt.input, len(got), len(tt.expected), got)
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Fatalf("collectLineSegmentsInBytes(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.expected[i])
				}
			}
		})
	}
}

// Benchmark the use of the line break function for byte slices.
func BenchmarkLineFunctionBytes(b *testing.B) {
	input := []byte(benchmarkStr)
	for i := 0; i < b.N; i++ {
		var c []byte
		var boundaries bool
		var state LineBreakState
		str := input
		for len(str) > 0 {
			c, str, _, state = FirstLineSegment(str, state)

			// to avoid the compiler optimizing out the benchmark
			runtime.KeepAlive(c)
			runtime.KeepAlive(str)
			runtime.KeepAlive(boundaries)
			runtime.KeepAlive(state)
		}
	}
}

// Benchmark the use of the line break function for strings.
func BenchmarkLineFunctionString(b *testing.B) {
	input := benchmarkStr
	for i := 0; i < b.N; i++ {
		var c string
		var boundaries bool
		var state LineBreakState
		str := input
		for len(str) > 0 {
			c, str, boundaries, state = FirstLineSegmentInString(str, state)

			// to avoid the compiler optimizing out the benchmark
			runtime.KeepAlive(c)
			runtime.KeepAlive(str)
			runtime.KeepAlive(boundaries)
			runtime.KeepAlive(state)
		}
	}
}

func FuzzFirstLineInString(f *testing.F) {
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
		var state LineBreakState
		var b []byte
		str := input
		for len(str) > 0 {
			var line string
			line, str, _, state = FirstLineSegmentInString(str, state)
			b = append(b, line...)
		}

		// Check if the constructed string is the same as the original.
		if string(b) != input {
			t.Errorf("Fuzzing failed: %q != %q", string(b), input)
		}
	})
}
