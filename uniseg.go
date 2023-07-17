package uniseg

//go:generate go run ./internal/cmd/gen_breaktest GraphemeBreakTest graphemebreak_test.go graphemeBreakTestCases graphemes
//go:generate go run ./internal/cmd/gen_breaktest WordBreakTest wordbreak_test.go wordBreakTestCases words
//go:generate go run ./internal/cmd/gen_breaktest SentenceBreakTest sentencebreak_test.go sentenceBreakTestCases sentences
//go:generate go run ./internal/cmd/gen_breaktest LineBreakTest linebreak_test.go lineBreakTestCases lines

//go:generate go run ./internal/cmd/gen_properties auxiliary/GraphemeBreakProperty graphemeproperties.go graphemeCodePoints graphemes emojis=Extended_Pictographic
//go:generate go run ./internal/cmd/gen_properties auxiliary/WordBreakProperty wordproperties.go workBreakCodePoints words emojis=Extended_Pictographic
//go:generate go run ./internal/cmd/gen_properties auxiliary/SentenceBreakProperty sentenceproperties.go sentenceBreakCodePoints sentences
//go:generate go run ./internal/cmd/gen_properties LineBreak lineproperties.go lineBreakCodePoints lines gencat
//go:generate go run ./internal/cmd/gen_properties EastAsianWidth eastasianwidth.go eastAsianWidth eastasianwidth
//go:generate go run ./internal/cmd/gen_properties - emojipresentation.go emojiPresentation emojipresentation emojis=Emoji_Presentation
