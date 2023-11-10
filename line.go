package uniseg

import "unicode/utf8"

// FirstLineSegment returns the prefix of the given byte slice after which a
// decision to break the string over to the next line can or must be made,
// according to the rules of [Unicode Standard Annex #14]. This is used to
// implement line breaking.
//
// Line breaking, also known as word wrapping, is the process of breaking a
// section of text into lines such that it will fit in the available width of a
// page, window or other display area.
//
// The returned "segment" may not be broken into smaller parts, unless no other
// breaking opportunities present themselves, in which case you may break by
// grapheme clusters (using the [FirstGraphemeCluster] function to determine the
// grapheme clusters).
//
// The "mustBreak" flag indicates whether you MUST break the line after the
// given segment (true), for example after newline characters, or you MAY break
// the line after the given segment (false).
//
// This function can be called continuously to extract all non-breaking sub-sets
// from a byte slice, as illustrated in the example below.
//
// If you don't know the current state, for example when calling the function
// for the first time, you must pass -1. For consecutive calls, pass the state
// and rest slice returned by the previous call.
//
// The "rest" slice is the sub-slice of the original byte slice "b" starting
// after the last byte of the identified line segment. If the length of the
// "rest" slice is 0, the entire byte slice "b" has been processed. The
// "segment" byte slice is the sub-slice of the input slice containing the
// identified line segment.
//
// Given an empty byte slice "b", the function returns nil values.
//
// Note that in accordance with [UAX #14 LB3], the final segment will end with
// "mustBreak" set to true. You can choose to ignore this by checking if the
// length of the "rest" slice is 0 and calling [HasTrailingLineBreak] or
// [HasTrailingLineBreakInString] on the last rune.
//
// Note also that this algorithm may break within grapheme clusters. This is
// addressed in Section 8.2 Example 6 of UAX #14. To avoid this, you can use
// the [Step] function instead.
//
// [Unicode Standard Annex #14]: https://www.unicode.org/reports/tr14/tr14-49.html
// [UAX #14 LB3]: https://www.unicode.org/reports/tr14/tr14-49.html#Algorithm
func FirstLineSegment(b []byte, state LineBreakState) (segment, rest []byte, mustBreak bool, newState LineBreakState) {
	return firstLineSegment(b, state, utf8.DecodeRune)
}

// FirstLineSegment returns the prefix of the given byte slice after which a
// decision to break the string over to the next line can or must be made,
// according to the rules of [Unicode Standard Annex #14]. This is used to
// implement line breaking.
//
// Line breaking, also known as word wrapping, is the process of breaking a
// section of text into lines such that it will fit in the available width of a
// page, window or other display area.
//
// The returned "segment" may not be broken into smaller parts, unless no other
// breaking opportunities present themselves, in which case you may break by
// grapheme clusters (using the [FirstGraphemeCluster] function to determine the
// grapheme clusters).
//
// The "mustBreak" flag indicates whether you MUST break the line after the
// given segment (true), for example after newline characters, or you MAY break
// the line after the given segment (false).
//
// This function can be called continuously to extract all non-breaking sub-sets
// from a byte slice, as illustrated in the example below.
//
// If you don't know the current state, for example when calling the function
// for the first time, you must pass -1. For consecutive calls, pass the state
// and rest slice returned by the previous call.
//
// The "rest" slice is the sub-slice of the original byte slice "b" starting
// after the last byte of the identified line segment. If the length of the
// "rest" slice is 0, the entire byte slice "b" has been processed. The
// "segment" byte slice is the sub-slice of the input slice containing the
// identified line segment.
//
// Given an empty byte slice "b", the function returns nil values.
//
// Note that in accordance with [UAX #14 LB3], the final segment will end with
// "mustBreak" set to true. You can choose to ignore this by checking if the
// length of the "rest" slice is 0 and calling [HasTrailingLineBreak] or
// [HasTrailingLineBreakInString] on the last rune.
//
// Note also that this algorithm may break within grapheme clusters. This is
// addressed in Section 8.2 Example 6 of UAX #14. To avoid this, you can use
// the [Step] function instead.
//
// [Unicode Standard Annex #14]: https://www.unicode.org/reports/tr14/tr14-49.html
// [UAX #14 LB3]: https://www.unicode.org/reports/tr14/tr14-49.html#Algorithm
func (*Parser) FirstLineSegment(b []byte, state LineBreakState) (segment, rest []byte, mustBreak bool, newState LineBreakState) {
	return firstLineSegment(b, state, utf8.DecodeRune)
}

// FirstLineSegmentInString is like [FirstLineSegment] but its input and outputs
// are strings.
func FirstLineSegmentInString(str string, state LineBreakState) (segment, rest string, mustBreak bool, newState LineBreakState) {
	return firstLineSegment(str, state, utf8.DecodeRuneInString)
}

// FirstLineSegmentInString is like [Parser.FirstLineSegment] but its input and outputs
// are strings.
func (*Parser) FirstLineSegmentInString(str string, state LineBreakState) (segment, rest string, mustBreak bool, newState LineBreakState) {
	return firstLineSegment(str, state, utf8.DecodeRuneInString)
}

func firstLineSegment[T bytes](str T, state LineBreakState, decoder runeDecoder[T]) (segment, rest T, mustBreak bool, newState LineBreakState) {
	var zero T

	// An empty byte slice returns nothing.
	if len(str) == 0 {
		return
	}

	// Extract the first rune.
	r, length := decoder(str)
	if len(str) <= length { // If we're already past the end, there is nothing else to parse.
		return str, zero, true, 0 // LB3.
	}

	// If we don't know the state, determine it now.
	if state <= 0 {
		state, _ = transitionLineBreakState(state, r, str[length:], decoder)
	}

	// Transition until we find a boundary.
	var boundary LineBreak
	for {
		r, l := decoder(str[length:])
		state, boundary = transitionLineBreakState(state, r, str[length+l:], decoder)

		if boundary != LineDontBreak {
			return str[:length], str[length:], boundary == LineMustBreak, state
		}

		length += l
		if len(str) <= length {
			return str, zero, true, 0 // LB3.
		}
	}
}

// HasTrailingLineBreak returns true if the last rune in the given byte slice is
// one of the hard line break code points defined in LB4 and LB5 of [UAX #14].
//
// [UAX #14]: https://www.unicode.org/reports/tr14/tr14-49.html#Algorithm
func HasTrailingLineBreak(b []byte) bool {
	r, _ := utf8.DecodeLastRune(b)
	p := lineBreakCodePoints.search(r).lbProperty
	return p == lbprBK || p == lbprCR || p == lbprLF || p == lbprNL
}

// HasTrailingLineBreak returns true if the last rune in the given byte slice is
// one of the hard line break code points defined in LB4 and LB5 of [UAX #14].
//
// [UAX #14]: https://www.unicode.org/reports/tr14/tr14-49.html#Algorithm
func (*Parser) HasTrailingLineBreak(b []byte) bool {
	r, _ := utf8.DecodeLastRune(b)
	p := lineBreakCodePoints.search(r).lbProperty
	return p == lbprBK || p == lbprCR || p == lbprLF || p == lbprNL
}

// HasTrailingLineBreakInString is like [HasTrailingLineBreak] but for a string.
func HasTrailingLineBreakInString(str string) bool {
	r, _ := utf8.DecodeLastRuneInString(str)
	p := lineBreakCodePoints.search(r).lbProperty
	return p == lbprBK || p == lbprCR || p == lbprLF || p == lbprNL
}

// HasTrailingLineBreakInString is like [HasTrailingLineBreak] but for a string.
func (*Parser) HasTrailingLineBreakInString(str string) bool {
	r, _ := utf8.DecodeLastRuneInString(str)
	p := lineBreakCodePoints.search(r).lbProperty
	return p == lbprBK || p == lbprCR || p == lbprLF || p == lbprNL
}
