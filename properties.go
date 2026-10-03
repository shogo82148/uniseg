package uniseg

import (
	"slices"
	"unicode/utf8"
)

// property is the Unicode property type.
type property int

// The Unicode properties as used in the various parsers. Only the ones needed
// in the context of this package are included.
const (
	prAny property = iota // prAny must be 0.
	prCR
	prLF
	prControl
	prExtend
	prZWJ
	prRegionalIndicator
	prPrepend
	prSpacingMark
	prL
	prV
	prT
	prLV
	prLVT
	prExtendedPictographic
	prMax = iota
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
	lbprXX lbProperty = iota // Unknown. lbprXX must be 0.

	// Non-tailorable Line Breaking Classes
	lbprBK  // Mandatory Break
	lbprCR  // Carriage Return
	lbprLF  // Line Feed
	lbprCM  // Combining Mark
	lbprNL  // Next Line
	lbprSG  // Surrogate
	lbprWJ  // Word Joiner
	lbprZW  // Zero Width Space
	lbprGL  // Non-breaking ("Glue")
	lbprSP  // Space
	lbprZWJ // Zero Width Joiner

	// Break Opportunities
	lbprB2 // Break Opportunity Before and After
	lbprBA // Break After
	lbprBB // Break Before
	lbprHY // Hyphen
	lbprHH // Hyphen and Hebrew hyphen
	lbprCB // Contingent Break Opportunity

	// Characters Prohibiting Certain Breaks
	lbprCL // Close Punctuation
	lbprCP // Close Parenthesis
	lbprEX // Exclamation/Interrogation
	lbprIN // Inseparable
	lbprNS // Nonstarter
	lbprOP // Open Punctuation
	lbprQU // Quotation

	// Numeric Context
	lbprIS // Infix Separator
	lbprNU // Numeric
	lbprPO // Postfix Numeric
	lbprPR // Prefix Numeric
	lbprSY // Symbols Allowing Break After

	// Other Characters
	lbprAI  // Ambiguous (Alphabetic or Ideograph)
	lbprAK  // Aksara
	lbprAL  // Alphabetic
	lbprAP  // Aksara Pre-Base
	lbprAS  // Aksara Start
	lbprCJ  // Conditional Japanese Starter
	lbprEB  // Emoji Base
	lbprEM  // Emoji Modifier
	lbprH2  // Hangul LV Syllable
	lbprH3  // Hangul LVT Syllable
	lbprHL  // Hebrew Letter
	lbprID  // Ideographic
	lbprJL  // Hangul L Jamo
	lbprJT  // Hangul T Jamo
	lbprJV  // Hangul V Jamo
	lbprRI  // Regional Indicator
	lbprSA  // Complex Context Dependent
	lbprVF  // Virama Final
	lbprVI  // Virama
	lbprMax = iota
)

type emojiProperty int8

const (
	_ emojiProperty = iota // reserved for the zero value
	prEmoji
	prEmojiPresentation
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

// incbProperty is the Indic_Conjunct_Break property.
type incbProperty int8

const (
	incbNone incbProperty = iota // incbNone must be 0.
	incbLinker
	incbConsonant
	incbExtend
)

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
	k := 0
	for k < len(d) {
		entry := d[k]
		if r < entry.runeRange.Lo {
			k = 2*k + 1
			continue
		}
		if r > entry.runeRange.Hi {
			k = 2*k + 2
			continue
		}
		return entry.value
	}

	var zero T
	return zero
}

// lookupTable wraps a [dictionary] with a precomputed table for ASCII code
// points, which are by far the most common ones in typical text. This avoids
// the tree search for them.
type lookupTable[T any] struct {
	ascii [utf8.RuneSelf]T
	dict  dictionary[T]
}

// newLookupTable returns a lookup table for the given dictionary.
func newLookupTable[T any](d dictionary[T]) *lookupTable[T] {
	t := &lookupTable[T]{dict: d}
	for r := range rune(utf8.RuneSelf) {
		t.ascii[r] = d.search(r)
	}
	return t
}

// search returns the value associated with the given rune.
func (t *lookupTable[T]) search(r rune) T {
	if uint32(r) < utf8.RuneSelf {
		return t.ascii[r]
	}
	return t.dict.search(r)
}

// Lookup tables used by the parsers.
var (
	graphemeLookup          = newLookupTable(graphemeCodePoints)
	incbLookup              = newLookupTable(incb)
	wordBreakLookup         = newLookupTable(workBreakCodePoints)
	sentenceBreakLookup     = newLookupTable(sentenceBreakCodePoints)
	lineBreakLookup         = newLookupTable(lineBreakCodePoints)
	eastAsianWidthLookup    = newLookupTable(eastAsianWidth)
	emojiLookup             = newLookupTable(emoji)
	emojiPresentationLookup = newLookupTable(emojiPresentation)
)

// transitionResult is an entry of the state transition tables of the parsers.
// It holds the new state, the breaking instruction, and the rule number.
// A rule number of 0 means that no transition is defined.
type transitionResult[S, B any] struct {
	state      S
	boundary   B
	ruleNumber int
}

// resolveTransitions fills the undefined entries of the transition table with
// the less specific transitions, so that the parsers can find the applicable
// transition with a single lookup. The table is queried as follows:
//
//  1. Find specific state + specific property. Stop if found.
//  2. Find specific state + any property.
//  3. Find any state + specific property.
//  4. If only (2) or (3) (but not both) was found, stop.
//  5. If both (2) and (3) were found, use state from (3) and breaking instruction
//     from the transition with the lower rule number, prefer (3) if rule numbers
//     are equal. Stop.
//  6. Use the default transition.
func resolveTransitions[S, B any](table []transitionResult[S, B], numProps, anyState, anyProp int, def transitionResult[S, B]) {
	original := slices.Clone(table)
	for i := range table {
		if original[i].ruleNumber > 0 {
			continue
		}
		state, prop := i/numProps, i%numProps
		transAnyProp := original[state*numProps+anyProp]
		transAnyState := original[anyState*numProps+prop]
		switch {
		case transAnyProp.ruleNumber > 0 && transAnyState.ruleNumber > 0:
			table[i] = transAnyState
			if transAnyProp.ruleNumber < transAnyState.ruleNumber {
				table[i].boundary = transAnyProp.boundary
				table[i].ruleNumber = transAnyProp.ruleNumber
			}
		case transAnyProp.ruleNumber > 0:
			table[i] = transAnyProp
		case transAnyState.ruleNumber > 0:
			table[i] = transAnyState
		default:
			table[i] = def
		}
	}
}
