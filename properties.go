package uniseg

// property is the Unicode property type.
type property int

// The Unicode properties as used in the various parsers. Only the ones needed
// in the context of this package are included.
const (
	prXX      property = 0    // Same as prAny.
	prAny     property = iota // prAny must be 0.
	prPrepend                 // Grapheme properties must come first, to reduce the number of bits stored in the state vector.
	prCR
	prLF
	prControl
	prExtend
	prRegionalIndicator
	prSpacingMark
	prL
	prV
	prT
	prLV
	prLVT
	prZWJ
	prExtendedPictographic
	prNewline
	prWSegSpace
	prDoubleQuote
	prSingleQuote
	prMidNumLet
	prNumeric
	prMidLetter
	prMidNum
	prExtendNumLet
	prALetter
	prFormat
	prHebrewLetter
	prKatakana
	prSp
	prSTerm
	prClose
	prSContinue
	prATerm
	prUpper
	prLower
	prSep
	prOLetter
	prCM
	prBA
	prBK
	prSP
	prEX
	prQU
	prAL
	prPR
	prPO
	prOP
	prCP
	prIS
	prHY
	prSY
	prNU
	prCL
	prNL
	prGL
	prAI
	prBB
	prHL
	prSA
	prJL
	prJV
	prJT
	prNS
	prZW
	prB2
	prIN
	prWJ
	prID
	prEB
	prCJ
	prH2
	prH3
	prSG
	prCB
	prRI
	prEM

	// Emoji
	prEmoji
	prEmojiPresentation
)

// East-Asian Width properties.
type eawProperty int8

// East-Asian Width properties.
const (
	eawprN  eawProperty = iota // Neutral (Not East Asian): https://www.unicode.org/reports/tr11/tr11-40.html#ED7
	eawprNa                    // East Asian Narrow (Na): https://www.unicode.org/reports/tr11/tr11-40.html#ED5
	eawprA                     // East Asian Ambiguous (A): https://www.unicode.org/reports/tr11/tr11-40.html#ED6
	eawprW                     // East Asian Wide (W): https://www.unicode.org/reports/tr11/tr11-40.html#ED4
	eawprH                     // East Asian Halfwidth (H): https://www.unicode.org/reports/tr11/tr11-40.html#ED3
	eawprF                     // East Asian Fullwidth (F): https://www.unicode.org/reports/tr11/tr11-40.html#ED2
)

// Word break properties.
type wbProperty int8

// Word break properties.
const (
	wbprAny wbProperty = iota // wbprAny must be 0.
	wbprCR
	wbprLF
	wbprNewline
	wbprExtend
	wbprZWJ
	wbprRegionalIndicator
	wbprFormat
	wbprKatakana
	wbprHebrewLetter
	wbprALetter
	wbprSingleQuote
	wbprDoubleQuote
	wbprMidNumLet
	wbprMidLetter
	wbprMidNum
	wbprNumeric
	wbprExtendNumLet
	wbprWSegSpace
	wbprExtendedPictographic
	wbprMax = iota
)

// Sentence break properties.
type sbProperty int8

// Sentence break properties.
const (
	sbprAny sbProperty = iota // sbprAny must be 0.
	sbprCR
	sbprLF
	sbprExtend
	sbprSep
	sbprFormat
	sbprSp
	sbprLower
	sbprUpper
	sbprOLetter
	sbprNumeric
	sbprATerm
	sbprSContinue
	sbprSTerm
	sbprClose
	sbprMax = iota
)

// Line break properties.
type lbProperty int8

// Line break properties.
const (
	lbprXX  lbProperty = iota // Unknown. lbprXX must be 0.
	lbprBK                    // Mandatory Break
	lbprCR                    // Carriage Return
	lbprLF                    // Line Feed
	lbprCM                    // Combining Mark
	lbprNL                    // Next Line
	lbprSG                    // Surrogate
	lbprWJ                    // Word Joiner
	lbprZW                    // Zero Width Space
	lbprGL                    // Non-breaking ("Glue")
	lbprSP                    // Space
	lbprZWJ                   // Zero Width Joiner
	lbprB2                    // Break Opportunity Before and After
	lbprBA                    // Break After
	lbprBB                    // Break Before
	lbprHY                    // Hyphen
	lbprCB                    // Contingent Break Opportunity
	lbprCL                    // Close Punctuation
	lbprCP                    // Close Parenthesis
	lbprEX                    // Exclamation/Interrogation
	lbprIN                    // Inseparable
	lbprNS                    // Nonstarter
	lbprOP                    // Open Punctuation
	lbprQU                    // Quotation
	lbprIS                    // Infix Separator
	lbprNU                    // Numeric
	lbprPO                    // Postfix Numeric
	lbprPR                    // Prefix Numeric
	lbprSY                    // Symbols Allowing Break After
	lbprAI                    // Ambiguous (Alphabetic or Ideograph)
	lbprAL                    // Alphabetic
	lbprCJ                    // Conditional Japanese Starter
	lbprEB                    // Emoji Base
	lbprEM                    // Emoji Modifier
	lbprH2                    // Hangul LV Syllable
	lbprH3                    // Hangul LVT Syllable
	lbprHL                    // Hebrew Letter
	lbprID                    // Ideographic
	lbprJL                    // Hangul L Jamo
	lbprJV                    // Hangul V Jamo
	lbprJT                    // Hangul T Jamo
	lbprRI                    // Regional Indicator
	lbprSA                    // Complex Context Dependent
	lbprMax = iota
)

// generalCategory is the Unicode General Categories.
type generalCategory int

// Unicode General Categories. Only the ones needed in the context of this
// package are included.
const (
	gcNone generalCategory = iota // gcNone must be 0.
	gcCc
	gcZs
	gcPo
	gcSc
	gcPs
	gcPe
	gcSm
	gcPd
	gcNd
	gcLu
	gcSk
	gcPc
	gcLl
	gcSo
	gcLo
	gcPi
	gcCf
	gcNo
	gcPf
	gcLC
	gcLm
	gcMn
	gcMe
	gcMc
	gcNl
	gcZl
	gcZp
	gcCn
	gcCs
	gcCo
)

type propertyGeneralCategory struct {
	lbProperty
	generalCategory
}

// Special code points.
const (
	vs15 = 0xfe0e // Variation Selector-15 (text presentation)
	vs16 = 0xfe0f // Variation Selector-16 (emoji presentation)
)

// runeRange represents of a range of Unicode code points.
// The range runs from Lo to Hi inclusive.
type runeRange struct {
	Lo rune
	Hi rune
}

type dictionaryEntry[T any] struct {
	runeRange runeRange
	value     T
}

type dictionary[T any] []dictionaryEntry[T]

// search returns the value associated with the given rune in the dictionary.
func (d dictionary[T]) search(r rune) T {
	from := 0
	to := len(d)
	for to > from {
		middle := int(uint(from+to) >> 1) // avoid overflow when computing middle
		entry := d[middle]
		if r < entry.runeRange.Lo {
			to = middle
			continue
		}
		if r > entry.runeRange.Hi {
			from = middle + 1
			continue
		}
		return entry.value
	}

	var zero T
	return zero
}
