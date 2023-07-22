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
)

type grStateProperty struct {
	grState
	property
}

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
var grTransitions = map[grStateProperty]grTransitionResult{
	// GB5
	{grAny, prCR}:      {grCR, true, 50},
	{grAny, prLF}:      {grControlLF, true, 50},
	{grAny, prControl}: {grControlLF, true, 50},

	// GB4
	{grCR, prAny}:        {grAny, true, 40},
	{grControlLF, prAny}: {grAny, true, 40},

	// GB3.
	{grCR, prLF}: {grControlLF, false, 30},

	// GB6.
	{grAny, prL}: {grL, true, 9990},
	{grL, prL}:   {grL, false, 60},
	{grL, prV}:   {grLVV, false, 60},
	{grL, prLV}:  {grLVV, false, 60},
	{grL, prLVT}: {grLVTT, false, 60},

	// GB7.
	{grAny, prLV}: {grLVV, true, 9990},
	{grAny, prV}:  {grLVV, true, 9990},
	{grLVV, prV}:  {grLVV, false, 70},
	{grLVV, prT}:  {grLVTT, false, 70},

	// GB8.
	{grAny, prLVT}: {grLVTT, true, 9990},
	{grAny, prT}:   {grLVTT, true, 9990},
	{grLVTT, prT}:  {grLVTT, false, 80},

	// GB9.
	{grAny, prExtend}: {grAny, false, 90},
	{grAny, prZWJ}:    {grAny, false, 90},

	// GB9a.
	{grAny, prSpacingMark}: {grAny, false, 91},

	// GB9b.
	{grAny, prPrepend}: {grPrepend, true, 9990},
	{grPrepend, prAny}: {grAny, false, 92},

	// GB11.
	{grAny, prExtendedPictographic}:                     {grExtendedPictographic, true, 9990},
	{grExtendedPictographic, prExtend}:                  {grExtendedPictographic, false, 110},
	{grExtendedPictographic, prZWJ}:                     {grExtendedPictographicZWJ, false, 110},
	{grExtendedPictographicZWJ, prExtendedPictographic}: {grExtendedPictographic, false, 110},

	// GB12 / GB13.
	{grAny, prRegionalIndicator}:    {grRIOdd, true, 9990},
	{grRIOdd, prRegionalIndicator}:  {grRIEven, false, 120},
	{grRIEven, prRegionalIndicator}: {grRIOdd, true, 120},
}

// transitionGraphemeState determines the new state of the grapheme cluster
// parser given the current state and the next code point. It also returns the
// code point's grapheme property (the value mapped by the [graphemeCodePoints]
// table) and whether a cluster boundary was detected.
func transitionGraphemeState(state grState, r rune) (newState grState, prop property, boundary bool) {
	// Determine the property of the next character.
	prop = graphemeCodePoints.search(r)

	// Find the applicable transition.
	transition, ok := grTransitions[grStateProperty{state, prop}]
	if ok {
		// We have a specific transition. We'll use it.
		return transition.grState, prop, transition.boundary
	}

	// No specific transition found. Try the less specific ones.
	transAnyProp, okAnyProp := grTransitions[grStateProperty{state, prAny}]
	transAnyState, okAnyState := grTransitions[grStateProperty{grAny, prop}]
	if okAnyProp && okAnyState {
		// Both apply. We'll use a mix (see comments for grTransitions).
		newState = transAnyState.grState
		boundary = transAnyState.boundary
		if transAnyProp.ruleNumber < transAnyState.ruleNumber {
			boundary = transAnyProp.boundary
		}
		return
	}

	if okAnyProp {
		// We only have a specific state.
		return transAnyProp.grState, prop, transAnyProp.boundary
		// This branch will probably never be reached because okAnyState will
		// always be true given the current transition map. But we keep it here
		// for future modifications to the transition map where this may not be
		// true anymore.
	}

	if okAnyState {
		// We only have a specific property.
		return transAnyState.grState, prop, transAnyState.boundary
	}

	// No known transition. GB999: Any ÷ Any.
	return grAny, prop, true
}
