package uniseg

import "unicode/utf8"

// SentenceBreakState is the state of the sentence break parser.
type SentenceBreakState int

// The states of the sentence break parser.
const (
	_ SentenceBreakState = iota // The zero value is reserved for the initial state.
	sbAny
	sbCR
	sbParaSep
	sbATerm
	sbUpper
	sbLower
	sbSB7
	sbSB8Close
	sbSB8Sp
	sbSTerm
	sbSB8aClose
	sbSB8aSp
	sbMax = iota
)

type sbTransitionResult struct {
	SentenceBreakState
	boundary   bool
	ruleNumber int
}

// The sentence break parser's state transitions. It's analogous to
// grTransitions, see comments there for details. Unicode version 15.0.0.
var sbTransitions = [sbMax * sbprMax]sbTransitionResult{
	// SB3.
	int(sbAny)*sbprMax + int(sbprCR): {sbCR, false, 9990},
	int(sbCR)*sbprMax + int(sbprLF):  {sbParaSep, false, 30},

	// SB4.
	int(sbAny)*sbprMax + int(sbprSep):     {sbParaSep, false, 9990},
	int(sbAny)*sbprMax + int(sbprLF):      {sbParaSep, false, 9990},
	int(sbParaSep)*sbprMax + int(sbprAny): {sbAny, true, 40},
	int(sbCR)*sbprMax + int(sbprAny):      {sbAny, true, 40},

	// SB6.
	int(sbAny)*sbprMax + int(sbprATerm):     {sbATerm, false, 9990},
	int(sbATerm)*sbprMax + int(sbprNumeric): {sbAny, false, 60},
	int(sbSB7)*sbprMax + int(sbprNumeric):   {sbAny, false, 60}, // Because ATerm also appears in SB7.

	// SB7.
	int(sbAny)*sbprMax + int(sbprUpper):   {sbUpper, false, 9990},
	int(sbAny)*sbprMax + int(sbprLower):   {sbLower, false, 9990},
	int(sbUpper)*sbprMax + int(sbprATerm): {sbSB7, false, 70},
	int(sbLower)*sbprMax + int(sbprATerm): {sbSB7, false, 70},
	int(sbSB7)*sbprMax + int(sbprUpper):   {sbUpper, false, 70},

	// SB8a.
	int(sbAny)*sbprMax + int(sbprSTerm):           {sbSTerm, false, 9990},
	int(sbATerm)*sbprMax + int(sbprSContinue):     {sbAny, false, 81},
	int(sbATerm)*sbprMax + int(sbprATerm):         {sbATerm, false, 81},
	int(sbATerm)*sbprMax + int(sbprSTerm):         {sbSTerm, false, 81},
	int(sbSB7)*sbprMax + int(sbprSContinue):       {sbAny, false, 81},
	int(sbSB7)*sbprMax + int(sbprATerm):           {sbATerm, false, 81},
	int(sbSB7)*sbprMax + int(sbprSTerm):           {sbSTerm, false, 81},
	int(sbSB8Close)*sbprMax + int(sbprSContinue):  {sbAny, false, 81},
	int(sbSB8Close)*sbprMax + int(sbprATerm):      {sbATerm, false, 81},
	int(sbSB8Close)*sbprMax + int(sbprSTerm):      {sbSTerm, false, 81},
	int(sbSB8Sp)*sbprMax + int(sbprSContinue):     {sbAny, false, 81},
	int(sbSB8Sp)*sbprMax + int(sbprATerm):         {sbATerm, false, 81},
	int(sbSB8Sp)*sbprMax + int(sbprSTerm):         {sbSTerm, false, 81},
	int(sbSTerm)*sbprMax + int(sbprSContinue):     {sbAny, false, 81},
	int(sbSTerm)*sbprMax + int(sbprATerm):         {sbATerm, false, 81},
	int(sbSTerm)*sbprMax + int(sbprSTerm):         {sbSTerm, false, 81},
	int(sbSB8aClose)*sbprMax + int(sbprSContinue): {sbAny, false, 81},
	int(sbSB8aClose)*sbprMax + int(sbprATerm):     {sbATerm, false, 81},
	int(sbSB8aClose)*sbprMax + int(sbprSTerm):     {sbSTerm, false, 81},
	int(sbSB8aSp)*sbprMax + int(sbprSContinue):    {sbAny, false, 81},
	int(sbSB8aSp)*sbprMax + int(sbprATerm):        {sbATerm, false, 81},
	int(sbSB8aSp)*sbprMax + int(sbprSTerm):        {sbSTerm, false, 81},

	// SB9.
	int(sbATerm)*sbprMax + int(sbprClose):     {sbSB8Close, false, 90},
	int(sbSB7)*sbprMax + int(sbprClose):       {sbSB8Close, false, 90},
	int(sbSB8Close)*sbprMax + int(sbprClose):  {sbSB8Close, false, 90},
	int(sbATerm)*sbprMax + int(sbprSp):        {sbSB8Sp, false, 90},
	int(sbSB7)*sbprMax + int(sbprSp):          {sbSB8Sp, false, 90},
	int(sbSB8Close)*sbprMax + int(sbprSp):     {sbSB8Sp, false, 90},
	int(sbSTerm)*sbprMax + int(sbprClose):     {sbSB8aClose, false, 90},
	int(sbSB8aClose)*sbprMax + int(sbprClose): {sbSB8aClose, false, 90},
	int(sbSTerm)*sbprMax + int(sbprSp):        {sbSB8aSp, false, 90},
	int(sbSB8aClose)*sbprMax + int(sbprSp):    {sbSB8aSp, false, 90},
	int(sbATerm)*sbprMax + int(sbprSep):       {sbParaSep, false, 90},
	int(sbATerm)*sbprMax + int(sbprCR):        {sbParaSep, false, 90},
	int(sbATerm)*sbprMax + int(sbprLF):        {sbParaSep, false, 90},
	int(sbSB7)*sbprMax + int(sbprSep):         {sbParaSep, false, 90},
	int(sbSB7)*sbprMax + int(sbprCR):          {sbParaSep, false, 90},
	int(sbSB7)*sbprMax + int(sbprLF):          {sbParaSep, false, 90},
	int(sbSB8Close)*sbprMax + int(sbprSep):    {sbParaSep, false, 90},
	int(sbSB8Close)*sbprMax + int(sbprCR):     {sbParaSep, false, 90},
	int(sbSB8Close)*sbprMax + int(sbprLF):     {sbParaSep, false, 90},
	int(sbSTerm)*sbprMax + int(sbprSep):       {sbParaSep, false, 90},
	int(sbSTerm)*sbprMax + int(sbprCR):        {sbParaSep, false, 90},
	int(sbSTerm)*sbprMax + int(sbprLF):        {sbParaSep, false, 90},
	int(sbSB8aClose)*sbprMax + int(sbprSep):   {sbParaSep, false, 90},
	int(sbSB8aClose)*sbprMax + int(sbprCR):    {sbParaSep, false, 90},
	int(sbSB8aClose)*sbprMax + int(sbprLF):    {sbParaSep, false, 90},

	// SB10.
	int(sbSB8Sp)*sbprMax + int(sbprSp):  {sbSB8Sp, false, 100},
	int(sbSB8aSp)*sbprMax + int(sbprSp): {sbSB8aSp, false, 100},
	int(sbSB8Sp)*sbprMax + int(sbprSep): {sbParaSep, false, 100},
	int(sbSB8Sp)*sbprMax + int(sbprCR):  {sbParaSep, false, 100},
	int(sbSB8Sp)*sbprMax + int(sbprLF):  {sbParaSep, false, 100},

	// SB11.
	int(sbATerm)*sbprMax + int(sbprAny):     {sbAny, true, 110},
	int(sbSB7)*sbprMax + int(sbprAny):       {sbAny, true, 110},
	int(sbSB8Close)*sbprMax + int(sbprAny):  {sbAny, true, 110},
	int(sbSB8Sp)*sbprMax + int(sbprAny):     {sbAny, true, 110},
	int(sbSTerm)*sbprMax + int(sbprAny):     {sbAny, true, 110},
	int(sbSB8aClose)*sbprMax + int(sbprAny): {sbAny, true, 110},
	int(sbSB8aSp)*sbprMax + int(sbprAny):    {sbAny, true, 110},
	// We'll always break after ParaSep due to SB4.
}

// transitionSentenceBreakState determines the new state of the sentence break
// parser given the current state and the next code point. It also returns
// whether a sentence boundary was detected. If more than one code point is
// needed to determine the new state, the byte slice or the string starting
// after rune "r" can be used (whichever is not nil or empty) for further
// lookups.
func transitionSentenceBreakState[T bytes](state SentenceBreakState, r rune, str T, decoder runeDecoder[T]) (newState SentenceBreakState, sentenceBreak bool) {
	// Determine the property of the next character.
	nextProperty := sentenceBreakCodePoints.search(r)

	// SB5 (Replacing Ignore Rules).
	if nextProperty == sbprExtend || nextProperty == sbprFormat {
		if state == sbParaSep || state == sbCR {
			return sbAny, true // Make sure we don't apply SB5 to SB3 or SB4.
		}
		if state < 0 {
			return sbAny, true // SB1.
		}
		return state, false
	}

	// Find the applicable transition in the table.
	var rule int
	transition := sbTransitions[int(state)*sbprMax+int(nextProperty)]
	if transition.ruleNumber > 0 {
		// We have a specific transition. We'll use it.
		newState, sentenceBreak, rule = transition.SentenceBreakState, transition.boundary, transition.ruleNumber
	} else {
		// No specific transition found. Try the less specific ones.
		transAnyProp := sbTransitions[int(state)*sbprMax+int(sbprAny)]
		transAnyState := sbTransitions[int(sbAny)*sbprMax+int(nextProperty)]
		if transAnyProp.ruleNumber > 0 && transAnyState.ruleNumber > 0 {
			// Both apply. We'll use a mix (see comments for grTransitions).
			newState, sentenceBreak, rule = transAnyState.SentenceBreakState, transAnyState.boundary, transAnyState.ruleNumber
			if transAnyProp.ruleNumber < transAnyState.ruleNumber {
				sentenceBreak, rule = transAnyProp.boundary, transAnyProp.ruleNumber
			}
		} else if transAnyProp.ruleNumber > 0 {
			// We only have a specific state.
			newState, sentenceBreak, rule = transAnyProp.SentenceBreakState, transAnyProp.boundary, transAnyProp.ruleNumber
			// This branch will probably never be reached because okAnyState will
			// always be true given the current transition map. But we keep it here
			// for future modifications to the transition map where this may not be
			// true anymore.
		} else if transAnyState.ruleNumber > 0 {
			// We only have a specific property.
			newState, sentenceBreak, rule = transAnyState.SentenceBreakState, transAnyState.boundary, transAnyState.ruleNumber
		} else {
			// No known transition. SB999: Any × Any.
			newState, sentenceBreak, rule = sbAny, false, 9990
		}
	}

	// SB8.
	if rule > 80 && (state == sbATerm || state == sbSB8Close || state == sbSB8Sp || state == sbSB7) {
		// Check the right side of the rule.
		var length int
		for nextProperty != sbprOLetter &&
			nextProperty != sbprUpper &&
			nextProperty != sbprLower &&
			nextProperty != sbprSep &&
			nextProperty != sbprCR &&
			nextProperty != sbprLF &&
			nextProperty != sbprATerm &&
			nextProperty != sbprSTerm {
			// Move on to the next rune.
			r, length = decoder(str)
			str = str[length:]
			if r == utf8.RuneError {
				break
			}
			nextProperty = sentenceBreakCodePoints.search(r)
		}
		if nextProperty == sbprLower {
			return sbLower, false
		}
	}

	return
}
