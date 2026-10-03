package uniseg_test

import (
	"fmt"

	"github.com/shogo82148/uniseg"
)

func ExampleGraphemeClusterCount() {
	n := uniseg.GraphemeClusterCount("🇩🇪🏳️\u200d🌈")
	fmt.Println(n)
	// Output: 2
}

func ExampleFirstGraphemeCluster() {
	b := []byte("🇩🇪🏳️\u200d🌈!")
	var state uniseg.GraphemeBreakState
	var c []byte
	for len(b) > 0 {
		var width int
		c, b, width, state = uniseg.FirstGraphemeCluster(b, state)
		fmt.Println(string(c), width)
	}
	// Output:
	// 🇩🇪 2
	// 🏳️‍🌈 2
	// ! 1
}

func ExampleFirstGraphemeClusterInString() {
	str := "🇩🇪🏳️\u200d🌈!"
	var state uniseg.GraphemeBreakState
	var c string
	for len(str) > 0 {
		var width int
		c, str, width, state = uniseg.FirstGraphemeClusterInString(str, state)
		fmt.Println(c, width)
	}
	// Output:
	// 🇩🇪 2
	// 🏳️‍🌈 2
	// ! 1
}

func ExampleGraphemeClusters() {
	b := []byte("🇩🇪🏳️\u200d🌈!")
	for i, c := range uniseg.GraphemeClusters(b) {
		fmt.Println(i, string(c))
	}
	// Output:
	// 0 🇩🇪
	// 8 🏳️‍🌈
	// 22 !
}

func ExampleGraphemeClustersInString() {
	str := "🇩🇪🏳️\u200d🌈!"
	for i, c := range uniseg.GraphemeClustersInString(str) {
		fmt.Println(i, c)
	}
	// Output:
	// 0 🇩🇪
	// 8 🏳️‍🌈
	// 22 !
}

func ExampleFirstWord() {
	b := []byte("Hello, world!")
	var state uniseg.WordBreakState
	var c []byte
	for len(b) > 0 {
		c, b, state = uniseg.FirstWord(b, state)
		fmt.Printf("(%s)\n", string(c))
	}
	// Output:
	// (Hello)
	// (,)
	// ( )
	// (world)
	// (!)
}

func ExampleFirstWordInString() {
	str := "Hello, world!"
	var state uniseg.WordBreakState
	var c string
	for len(str) > 0 {
		c, str, state = uniseg.FirstWordInString(str, state)
		fmt.Printf("(%s)\n", c)
	}
	// Output:
	// (Hello)
	// (,)
	// ( )
	// (world)
	// (!)
}

func ExampleWords() {
	b := []byte("Hello, world!")
	for i, w := range uniseg.Words(b) {
		fmt.Printf("%d (%s)\n", i, string(w))
	}
	// Output:
	// 0 (Hello)
	// 5 (,)
	// 6 ( )
	// 7 (world)
	// 12 (!)
}

func ExampleWordsInString() {
	str := "Hello, world!"
	for i, w := range uniseg.WordsInString(str) {
		fmt.Printf("%d (%s)\n", i, w)
	}
	// Output:
	// 0 (Hello)
	// 5 (,)
	// 6 ( )
	// 7 (world)
	// 12 (!)
}

func ExampleFirstSentence() {
	b := []byte("This is sentence 1.0. And this is sentence two.")
	var state uniseg.SentenceBreakState
	var c []byte
	for len(b) > 0 {
		c, b, state = uniseg.FirstSentence(b, state)
		fmt.Printf("(%s)\n", string(c))
	}
	// Output:
	// (This is sentence 1.0. )
	// (And this is sentence two.)
}

func ExampleFirstSentenceInString() {
	str := "This is sentence 1.0. And this is sentence two."
	var state uniseg.SentenceBreakState
	var c string
	for len(str) > 0 {
		c, str, state = uniseg.FirstSentenceInString(str, state)
		fmt.Printf("(%s)\n", c)
	}
	// Output:
	// (This is sentence 1.0. )
	// (And this is sentence two.)
}

func ExampleSentences() {
	b := []byte("This is sentence 1.0. And this is sentence two.")
	for i, s := range uniseg.Sentences(b) {
		fmt.Printf("%d (%s)\n", i, string(s))
	}
	// Output:
	// 0 (This is sentence 1.0. )
	// 22 (And this is sentence two.)
}

func ExampleSentencesInString() {
	str := "This is sentence 1.0. And this is sentence two."
	for i, s := range uniseg.SentencesInString(str) {
		fmt.Printf("%d (%s)\n", i, s)
	}
	// Output:
	// 0 (This is sentence 1.0. )
	// 22 (And this is sentence two.)
}

func ExampleFirstLineSegment() {
	b := []byte("First line.\nSecond line.")
	var (
		c         []byte
		mustBreak bool
		state     uniseg.LineBreakState
	)
	for len(b) > 0 {
		c, b, mustBreak, state = uniseg.FirstLineSegment(b, state)
		fmt.Printf("(%s)", string(c))
		if mustBreak {
			fmt.Print("!")
		}
	}
	// Output:
	// (First )(line.
	// )!(Second )(line.)!
}

func ExampleFirstLineSegmentInString() {
	str := "First line.\nSecond line."
	var (
		c         string
		mustBreak bool
		state     uniseg.LineBreakState
	)
	for len(str) > 0 {
		c, str, mustBreak, state = uniseg.FirstLineSegmentInString(str, state)
		fmt.Printf("(%s)", c)
		if mustBreak {
			fmt.Println(" < must break")
		} else {
			fmt.Println(" < may break")
		}
	}
	// Output:
	// (First ) < may break
	// (line.
	// ) < must break
	// (Second ) < may break
	// (line.) < must break
}

func ExampleStep_graphemes() {
	b := []byte("🇩🇪🏳️\u200d🌈!")
	var (
		c          []byte
		boundaries uniseg.Boundaries
		state      uniseg.State
	)
	for len(b) > 0 {
		c, b, boundaries, state = uniseg.Step(b, state)
		fmt.Println(string(c), boundaries.Width())
	}
	// Output: 🇩🇪 2
	// 🏳️‍🌈 2
	// ! 1
}

func ExampleStepString_graphemes() {
	str := "🇩🇪🏳️\u200d🌈!"
	var c string
	var state uniseg.State
	for len(str) > 0 {
		var boundaries uniseg.Boundaries
		c, str, boundaries, state = uniseg.StepString(str, state)
		fmt.Println(string(c), boundaries.Width())
	}
	// Output:
	// 🇩🇪 2
	// 🏳️‍🌈 2
	// ! 1
}

func ExampleStep_word() {
	b := []byte("Hello, world!")
	var (
		c          []byte
		boundaries uniseg.Boundaries
		state      uniseg.State
	)
	for len(b) > 0 {
		c, b, boundaries, state = uniseg.Step(b, state)
		fmt.Print(string(c))
		if boundaries.Word() {
			fmt.Print("|")
		}
	}
	// Output: Hello|,| |world|!|
}

func ExampleStepString_word() {
	str := "Hello, world!"
	var (
		c          string
		boundaries uniseg.Boundaries
		state      uniseg.State
	)
	for len(str) > 0 {
		c, str, boundaries, state = uniseg.StepString(str, state)
		fmt.Print(c)
		if boundaries.Word() {
			fmt.Print("|")
		}
	}
	// Output: Hello|,| |world|!|
}

func ExampleStep_sentence() {
	b := []byte("This is sentence 1.0. And this is sentence two.")
	var (
		c          []byte
		boundaries uniseg.Boundaries
		state      uniseg.State
	)
	for len(b) > 0 {
		c, b, boundaries, state = uniseg.Step(b, state)
		fmt.Print(string(c))
		if boundaries.Sentence() {
			fmt.Print("|")
		}
	}
	// Output: This is sentence 1.0. |And this is sentence two.|
}

func ExampleStepString_sentence() {
	str := "This is sentence 1.0. And this is sentence two."
	var (
		c          string
		boundaries uniseg.Boundaries
		state      uniseg.State
	)
	for len(str) > 0 {
		c, str, boundaries, state = uniseg.StepString(str, state)
		fmt.Print(c)
		if boundaries.Sentence() {
			fmt.Print("|")
		}
	}
	// Output: This is sentence 1.0. |And this is sentence two.|
}

func ExampleStep_lineBreaking() {
	b := []byte("First line.\nSecond line.")
	var (
		c          []byte
		boundaries uniseg.Boundaries
		state      uniseg.State
	)
	for len(b) > 0 {
		c, b, boundaries, state = uniseg.Step(b, state)
		fmt.Print(string(c))
		switch boundaries.Line() {
		case uniseg.LineCanBreak:
			fmt.Print("|")
		case uniseg.LineMustBreak:
			fmt.Print("‖")
		}
	}
	// Output:
	// First |line.
	// ‖Second |line.‖
}

func ExampleStepString_lineBreaking() {
	str := "First line.\nSecond line."
	var (
		c          string
		boundaries uniseg.Boundaries
		state      uniseg.State
	)
	for len(str) > 0 {
		c, str, boundaries, state = uniseg.StepString(str, state)
		fmt.Print(c)
		switch boundaries.Line() {
		case uniseg.LineCanBreak:
			fmt.Print("|")
		case uniseg.LineMustBreak:
			fmt.Print("‖")
		}
	}
	// Output: First |line.
	//‖Second |line.‖
}

func ExampleGraphemes_graphemes() {
	g := uniseg.NewGraphemes("🇩🇪🏳️\u200d🌈")
	for g.Next() {
		fmt.Println(g.Str())
	}
	// Output:
	// 🇩🇪
	// 🏳️‍🌈
}

func ExampleGraphemes_word() {
	g := uniseg.NewGraphemes("Hello, world!")
	for g.Next() {
		fmt.Print(g.Str())
		if g.IsWordBoundary() {
			fmt.Print("|")
		}
	}
	// Output: Hello|,| |world|!|
}

func ExampleGraphemes_sentence() {
	g := uniseg.NewGraphemes("This is sentence 1.0. And this is sentence two.")
	for g.Next() {
		fmt.Print(g.Str())
		if g.IsSentenceBoundary() {
			fmt.Print("|")
		}
	}
	// Output: This is sentence 1.0. |And this is sentence two.|
}

func ExampleGraphemes_lineBreaking() {
	g := uniseg.NewGraphemes("First line.\nSecond line.")
	for g.Next() {
		fmt.Print(g.Str())
		switch g.LineBreak() {
		case uniseg.LineCanBreak:
			fmt.Print("|")
		case uniseg.LineMustBreak:
			fmt.Print("‖")
		}
	}
	// Output: First |line.
	//‖Second |line.‖
}

func ExampleStringWidth() {
	fmt.Println(uniseg.StringWidth("Hello, 世界"))
	// Output: 11
}
