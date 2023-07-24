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
	wbprWSegSpace
	wbprHebrewLetter
	wbprALetter
	wbprWB7
	wbprWB7c
	wbprNumeric
	wbprWB11
	wbprKatakana
	wbprExtendNumLet
	wbprOddRI
	wbprEvenRI
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
	property
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
