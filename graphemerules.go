package uniseg

// grState is a state of the grapheme cluster parser.
type grState int

// The states of the grapheme cluster parser.
const (
	_ grState = iota // The zero value is reserved for the initial state.
	grAny
	grCR
	grControlLF
	grL
	grLVV
	grLVTT
	grPrepend
	grExtendedPictographic
	grExtendedPictographicZWJ
	grRIOdd
	grRIEven
	grMax = iota
)

type grTransitionResult struct {
	grState
	boundary   bool
	ruleNumber int
}

// The grapheme cluster parser's state transitions. Maps (state, property) to
// (new state, breaking instruction, rule number). The breaking instruction
// always refers to the boundary between the last and next code point.
//
// This map is queried as follows:
//
//  1. Find specific state + specific property. Stop if found.
//  2. Find specific state + any property.
//  3. Find any state + specific property.
//  4. If only (2) or (3) (but not both) was found, stop.
//  5. If both (2) and (3) were found, use state from (3) and breaking instruction
//     from the transition with the lower rule number, prefer (3) if rule numbers
//     are equal. Stop.
//  6. Assume grAny and grBoundary.
//
// Unicode version 15.0.0.
var grTransitions = [grMax * prMax]grTransitionResult{
	// GB5
	int(grAny)*prMax + int(prCR):      {grCR, true, 50},
	int(grAny)*prMax + int(prLF):      {grControlLF, true, 50},
	int(grAny)*prMax + int(prControl): {grControlLF, true, 50},

	// GB4
	int(grCR)*prMax + int(prAny):        {grAny, true, 40},
	int(grControlLF)*prMax + int(prAny): {grAny, true, 40},

	// GB3.
	int(grCR)*prMax + int(prLF): {grControlLF, false, 30},

	// GB6.
	int(grAny)*prMax + int(prL): {grL, true, 9990},
	int(grL)*prMax + int(prL):   {grL, false, 60},
	int(grL)*prMax + int(prV):   {grLVV, false, 60},
	int(grL)*prMax + int(prLV):  {grLVV, false, 60},
	int(grL)*prMax + int(prLVT): {grLVTT, false, 60},

	// GB7.
	int(grAny)*prMax + int(prLV): {grLVV, true, 9990},
	int(grAny)*prMax + int(prV):  {grLVV, true, 9990},
	int(grLVV)*prMax + int(prV):  {grLVV, false, 70},
	int(grLVV)*prMax + int(prT):  {grLVTT, false, 70},

	// GB8.
	int(grAny)*prMax + int(prLVT): {grLVTT, true, 9990},
	int(grAny)*prMax + int(prT):   {grLVTT, true, 9990},
	int(grLVTT)*prMax + int(prT):  {grLVTT, false, 80},

	// GB9.
	int(grAny)*prMax + int(prExtend): {grAny, false, 90},
	int(grAny)*prMax + int(prZWJ):    {grAny, false, 90},

	// GB9a.
	int(grAny)*prMax + int(prSpacingMark): {grAny, false, 91},

	// GB9b.
	int(grAny)*prMax + int(prPrepend): {grPrepend, true, 9990},
	int(grPrepend)*prMax + int(prAny): {grAny, false, 92},

	// GB11.
	int(grAny)*prMax + int(prExtendedPictographic):                     {grExtendedPictographic, true, 9990},
	int(grExtendedPictographic)*prMax + int(prExtend):                  {grExtendedPictographic, false, 110},
	int(grExtendedPictographic)*prMax + int(prZWJ):                     {grExtendedPictographicZWJ, false, 110},
	int(grExtendedPictographicZWJ)*prMax + int(prExtendedPictographic): {grExtendedPictographic, false, 110},

	// GB12 / GB13.
	int(grAny)*prMax + int(prRegionalIndicator):    {grRIOdd, true, 9990},
	int(grRIOdd)*prMax + int(prRegionalIndicator):  {grRIEven, false, 120},
	int(grRIEven)*prMax + int(prRegionalIndicator): {grRIOdd, true, 120},
}

// transitionGraphemeState determines the new state of the grapheme cluster
// parser given the current state and the next code point. It also returns the
// code point's grapheme property (the value mapped by the [graphemeCodePoints]
// table) and whether a cluster boundary was detected.
func transitionGraphemeState(state grState, r rune) (newState grState, prop property, boundary bool) {
	// Determine the property of the next character.
	prop = graphemeCodePoints.search(r)

	// Find the applicable transition.
	transition := grTransitions[int(state)*prMax+int(prop)]
	if transition.ruleNumber > 0 {
		// We have a specific transition. We'll use it.
		return transition.grState, prop, transition.boundary
	}

	// No specific transition found. Try the less specific ones.
	transAnyProp := grTransitions[int(state)*prMax+int(prAny)]
	transAnyState := grTransitions[int(grAny)*prMax+int(prop)]
	if transAnyProp.ruleNumber > 0 && transAnyState.ruleNumber > 0 {
		// Both apply. We'll use a mix (see comments for grTransitions).
		newState = transAnyState.grState
		boundary = transAnyState.boundary
		if transAnyProp.ruleNumber < transAnyState.ruleNumber {
			boundary = transAnyProp.boundary
		}
		return
	}

	if transAnyProp.ruleNumber > 0 {
		// We only have a specific state.
		return transAnyProp.grState, prop, transAnyProp.boundary
		// This branch will probably never be reached because okAnyState will
		// always be true given the current transition map. But we keep it here
		// for future modifications to the transition map where this may not be
		// true anymore.
	}

	if transAnyState.ruleNumber > 0 {
		// We only have a specific property.
		return transAnyState.grState, prop, transAnyState.boundary
	}

	// No known transition. GB999: Any ÷ Any.
	return grAny, prop, true
}
