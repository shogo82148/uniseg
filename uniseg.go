package uniseg

import "os"

//go:generate go run ./internal/cmd/gen_breaktest GraphemeBreakTest graphemebreak_test.go graphemeBreakTestCases graphemes
//go:generate go run ./internal/cmd/gen_breaktest WordBreakTest wordbreak_test.go wordBreakTestCases words
//go:generate go run ./internal/cmd/gen_breaktest SentenceBreakTest sentencebreak_test.go sentenceBreakTestCases sentences
//go:generate go run ./internal/cmd/gen_breaktest LineBreakTest linebreak_test.go lineBreakTestCases lines

//go:generate go run ./internal/cmd/gen_properties -logprefix=graphemes -property=auxiliary/GraphemeBreakProperty -emojis=Extended_Pictographic graphemeproperties.go graphemeCodePoints
//go:generate go run ./internal/cmd/gen_properties -logprefix=words -property=auxiliary/WordBreakProperty -emojis=Extended_Pictographic -prefix=wbpr -type=wbProperty wordproperties.go workBreakCodePoints
//go:generate go run ./internal/cmd/gen_properties -logprefix=sentences -property=auxiliary/SentenceBreakProperty -prefix=sbpr -type=sbProperty sentenceproperties.go sentenceBreakCodePoints
//go:generate go run ./internal/cmd/gen_properties -logprefix=lines -property=LineBreak -gencat -prefix=lbpr -type=propertyGeneralCategory lineproperties.go lineBreakCodePoints
//go:generate go run ./internal/cmd/gen_properties -logprefix=eastasianwidth -property=EastAsianWidth -prefix=eawpr -type eawProperty eastasianwidth.go eastAsianWidth
//go:generate go run ./internal/cmd/gen_properties -logprefix=emojipresentation -emojis=Emoji_Presentation -type=emojiProperty emojipresentation.go emojiPresentation
//go:generate go run ./internal/cmd/gen_properties -logprefix=emoji -emojis=Emoji -type=emojiProperty emoji.go emoji

// Parser is a parser for Unicode text.
type Parser struct {
	// EastAsianWidth controls the width of characters
	// with the East Asian Width Ambiguous attribute.
	//
	// It it is true, the parser treats Unicode text
	// in the context of East Asian traditional character encodings.
	// The width of characters with the East Asian Width Ambiguous attribute is 2.
	//
	// It it is false, the parser treats Unicode text
	// in the context of non-East Asian traditional character encodings.
	// The width of characters with the East Asian Width Ambiguous attribute is 1.
	EastAsianWidth bool

	// WideEmoji controls the width of Emoji characters.
	// [UAX #11] recommends that Emoji characters should be rendered
	// with a width of 1, however some fonts render Emoji characters
	// wider than other characters.
	// WideEmoji is used to maintain compatibility with such fonts.
	//
	// If it is true, the width of Emoji characters is 2.
	// Otherwise, the width of Emoji characters is 1.
	// It is effective only when EastAsianWidth is true.
	//
	// [UAX #11]: https://www.unicode.org/reports/tr11/tr11-40.html
	WideEmoji bool
}

var DefaultParser = defaultParser()

func defaultParser() *Parser {
	p := &Parser{}
	// it is compatible with https://github.com/mattn/go-runewidth.
	env := os.Getenv("RUNEWIDTH_EASTASIAN")
	if env == "" {
		p.EastAsianWidth = IsEastAsian()
	} else {
		p.EastAsianWidth = env == "1"
	}
	return p
}

// bytes is a type that is either []byte or string.
type bytes interface {
	~string | ~[]byte
}

// runeDecoder is a function that decodes a rune from bytes.
type runeDecoder[T any] func(s T) (r rune, size int)
