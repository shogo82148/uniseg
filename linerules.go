package uniseg

// LineBreakState is the type of the line break parser's states.
type LineBreakState int

// LineBreak defines whether a given text may be broken into the next line.
type LineBreak int

// These constants define whether a given text may be broken into the next line.
// If the break is optional (LineCanBreak), you may choose to break or not based
// on your own criteria, for example, if the text has reached the available
// width.
const (
	LineDontBreak LineBreak = iota // You may not break the line here.
	LineCanBreak                   // You may or may not break the line here.
	LineMustBreak                  // You must break the line here.
)

// transitionLineBreakState determines the new state of the line break parser
// given the current state and the next code point. It also returns the type of
// line break: LineDontBreak, LineCanBreak, or LineMustBreak. If more than one
// code point is needed to determine the new state, the byte slice or the string
// starting after rune "r" can be used (whichever is not nil or empty) for
// further lookups.
func transitionLineBreakState[T bytes](state LineBreakState, r rune, str T, decoder runeDecoder[T]) (newState LineBreakState, lineBreak LineBreak) {
	// Determine the property of the next character.
	lbProp := lineBreakCodePoints.search(r)
	nextProperty := lbProp.lbProperty
	generalCategory := lbProp.generalCategory

	_ = nextProperty
	_ = generalCategory
	return
}
