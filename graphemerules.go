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

const (
	grStateMask     grState = 0x0f
	grGB9cStateMask grState = 0xf0

	// GB9c states.
	grGB9c1 grState = 0x10 // seen \p{InCB=Consonant}
	grGB9c2 grState = 0x20 // seen \p{InCB=Linker}
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
// Unicode version 15.1.0.
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
	incbProp := incb.search(r)

	// Find the applicable transition.
	gb9cState := state & grGB9cStateMask
	state &= grStateMask
	transition := grTransitions[int(state)*prMax+int(prop)]
	ruleNumber := 0
	if transition.ruleNumber > 0 {
		// We have a specific transition.
		ruleNumber = transition.ruleNumber
		newState = transition.grState
		boundary = transition.boundary
	} else {
		// No specific transition found. Try the less specific ones.
		transAnyProp := grTransitions[int(state)*prMax+int(prAny)]
		transAnyState := grTransitions[int(grAny)*prMax+int(prop)]
		if transAnyProp.ruleNumber > 0 && transAnyState.ruleNumber > 0 {
			// Both apply. We'll use a mix (see comments for grTransitions).
			ruleNumber = transAnyState.ruleNumber
			newState = transAnyState.grState
			boundary = transAnyState.boundary
			if transAnyProp.ruleNumber < transAnyState.ruleNumber {
				ruleNumber = transAnyProp.ruleNumber
				boundary = transAnyProp.boundary
			}
		} else if transAnyProp.ruleNumber > 0 {
			ruleNumber = transAnyProp.ruleNumber
			newState = transAnyProp.grState
			boundary = transAnyProp.boundary
		} else if transAnyState.ruleNumber > 0 {
			ruleNumber = transAnyState.ruleNumber
			newState = transAnyState.grState
			boundary = transAnyState.boundary
		} else {
			// No known transition. GB999: Any ÷ Any.
			ruleNumber = 9990
			newState = grAny
			boundary = true
		}
	}

	// GB9c: \p{InCB=Consonant} [ \p{InCB=Extend} \p{InCB=Linker} ]* \p{InCB=Linker} [ \p{InCB=Extend} \p{InCB=Linker} ]* 	× 	\p{InCB=Consonant}
	if ruleNumber >= 93 {
		if gb9cState == grGB9c2 && incbProp == incbConsonant {
			boundary = false
		}
	}

	var newGBcState grState

	// GB9c: state transition
	switch gb9cState {
	case grGB9c1:
		switch incbProp {
		case incbLinker:
			newGBcState = grGB9c2
		case incbExtend:
			newGBcState = grGB9c1
		}
	case grGB9c2:
		switch incbProp {
		case incbLinker, incbExtend:
			newGBcState = grGB9c2
		}
	}
	if incbProp == incbConsonant {
		newGBcState = grGB9c1
	}

	newState |= newGBcState

	return
}
