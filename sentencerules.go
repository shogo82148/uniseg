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
)

type sbStateProperty struct {
	SentenceBreakState
	property
}

type sbTransitionResult struct {
	SentenceBreakState
	boundary   bool
	ruleNumber int
}

// The sentence break parser's state transitions. It's analogous to
// grTransitions, see comments there for details. Unicode version 15.0.0.
var sbTransitions = map[sbStateProperty]sbTransitionResult{
	// SB3.
	{sbAny, prCR}: {sbCR, false, 9990},
	{sbCR, prLF}:  {sbParaSep, false, 30},

	// SB4.
	{sbAny, prSep}:     {sbParaSep, false, 9990},
	{sbAny, prLF}:      {sbParaSep, false, 9990},
	{sbParaSep, prAny}: {sbAny, true, 40},
	{sbCR, prAny}:      {sbAny, true, 40},

	// SB6.
	{sbAny, prATerm}:     {sbATerm, false, 9990},
	{sbATerm, prNumeric}: {sbAny, false, 60},
	{sbSB7, prNumeric}:   {sbAny, false, 60}, // Because ATerm also appears in SB7.

	// SB7.
	{sbAny, prUpper}:   {sbUpper, false, 9990},
	{sbAny, prLower}:   {sbLower, false, 9990},
	{sbUpper, prATerm}: {sbSB7, false, 70},
	{sbLower, prATerm}: {sbSB7, false, 70},
	{sbSB7, prUpper}:   {sbUpper, false, 70},

	// SB8a.
	{sbAny, prSTerm}:           {sbSTerm, false, 9990},
	{sbATerm, prSContinue}:     {sbAny, false, 81},
	{sbATerm, prATerm}:         {sbATerm, false, 81},
	{sbATerm, prSTerm}:         {sbSTerm, false, 81},
	{sbSB7, prSContinue}:       {sbAny, false, 81},
	{sbSB7, prATerm}:           {sbATerm, false, 81},
	{sbSB7, prSTerm}:           {sbSTerm, false, 81},
	{sbSB8Close, prSContinue}:  {sbAny, false, 81},
	{sbSB8Close, prATerm}:      {sbATerm, false, 81},
	{sbSB8Close, prSTerm}:      {sbSTerm, false, 81},
	{sbSB8Sp, prSContinue}:     {sbAny, false, 81},
	{sbSB8Sp, prATerm}:         {sbATerm, false, 81},
	{sbSB8Sp, prSTerm}:         {sbSTerm, false, 81},
	{sbSTerm, prSContinue}:     {sbAny, false, 81},
	{sbSTerm, prATerm}:         {sbATerm, false, 81},
	{sbSTerm, prSTerm}:         {sbSTerm, false, 81},
	{sbSB8aClose, prSContinue}: {sbAny, false, 81},
	{sbSB8aClose, prATerm}:     {sbATerm, false, 81},
	{sbSB8aClose, prSTerm}:     {sbSTerm, false, 81},
	{sbSB8aSp, prSContinue}:    {sbAny, false, 81},
	{sbSB8aSp, prATerm}:        {sbATerm, false, 81},
	{sbSB8aSp, prSTerm}:        {sbSTerm, false, 81},

	// SB9.
	{sbATerm, prClose}:     {sbSB8Close, false, 90},
	{sbSB7, prClose}:       {sbSB8Close, false, 90},
	{sbSB8Close, prClose}:  {sbSB8Close, false, 90},
	{sbATerm, prSp}:        {sbSB8Sp, false, 90},
	{sbSB7, prSp}:          {sbSB8Sp, false, 90},
	{sbSB8Close, prSp}:     {sbSB8Sp, false, 90},
	{sbSTerm, prClose}:     {sbSB8aClose, false, 90},
	{sbSB8aClose, prClose}: {sbSB8aClose, false, 90},
	{sbSTerm, prSp}:        {sbSB8aSp, false, 90},
	{sbSB8aClose, prSp}:    {sbSB8aSp, false, 90},
	{sbATerm, prSep}:       {sbParaSep, false, 90},
	{sbATerm, prCR}:        {sbParaSep, false, 90},
	{sbATerm, prLF}:        {sbParaSep, false, 90},
	{sbSB7, prSep}:         {sbParaSep, false, 90},
	{sbSB7, prCR}:          {sbParaSep, false, 90},
	{sbSB7, prLF}:          {sbParaSep, false, 90},
	{sbSB8Close, prSep}:    {sbParaSep, false, 90},
	{sbSB8Close, prCR}:     {sbParaSep, false, 90},
	{sbSB8Close, prLF}:     {sbParaSep, false, 90},
	{sbSTerm, prSep}:       {sbParaSep, false, 90},
	{sbSTerm, prCR}:        {sbParaSep, false, 90},
	{sbSTerm, prLF}:        {sbParaSep, false, 90},
	{sbSB8aClose, prSep}:   {sbParaSep, false, 90},
	{sbSB8aClose, prCR}:    {sbParaSep, false, 90},
	{sbSB8aClose, prLF}:    {sbParaSep, false, 90},

	// SB10.
	{sbSB8Sp, prSp}:  {sbSB8Sp, false, 100},
	{sbSB8aSp, prSp}: {sbSB8aSp, false, 100},
	{sbSB8Sp, prSep}: {sbParaSep, false, 100},
	{sbSB8Sp, prCR}:  {sbParaSep, false, 100},
	{sbSB8Sp, prLF}:  {sbParaSep, false, 100},

	// SB11.
	{sbATerm, prAny}:     {sbAny, true, 110},
	{sbSB7, prAny}:       {sbAny, true, 110},
	{sbSB8Close, prAny}:  {sbAny, true, 110},
	{sbSB8Sp, prAny}:     {sbAny, true, 110},
	{sbSTerm, prAny}:     {sbAny, true, 110},
	{sbSB8aClose, prAny}: {sbAny, true, 110},
	{sbSB8aSp, prAny}:    {sbAny, true, 110},
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
	if nextProperty == prExtend || nextProperty == prFormat {
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
	transition, ok := sbTransitions[sbStateProperty{state, nextProperty}]
	if ok {
		// We have a specific transition. We'll use it.
		newState, sentenceBreak, rule = transition.SentenceBreakState, transition.boundary, transition.ruleNumber
	} else {
		// No specific transition found. Try the less specific ones.
		transAnyProp, okAnyProp := sbTransitions[sbStateProperty{state, prAny}]
		transAnyState, okAnyState := sbTransitions[sbStateProperty{sbAny, nextProperty}]
		if okAnyProp && okAnyState {
			// Both apply. We'll use a mix (see comments for grTransitions).
			newState, sentenceBreak, rule = transAnyState.SentenceBreakState, transAnyState.boundary, transAnyState.ruleNumber
			if transAnyProp.ruleNumber < transAnyState.ruleNumber {
				sentenceBreak, rule = transAnyProp.boundary, transAnyProp.ruleNumber
			}
		} else if okAnyProp {
			// We only have a specific state.
			newState, sentenceBreak, rule = transAnyProp.SentenceBreakState, transAnyProp.boundary, transAnyProp.ruleNumber
			// This branch will probably never be reached because okAnyState will
			// always be true given the current transition map. But we keep it here
			// for future modifications to the transition map where this may not be
			// true anymore.
		} else if okAnyState {
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
		for nextProperty != prOLetter &&
			nextProperty != prUpper &&
			nextProperty != prLower &&
			nextProperty != prSep &&
			nextProperty != prCR &&
			nextProperty != prLF &&
			nextProperty != prATerm &&
			nextProperty != prSTerm {
			// Move on to the next rune.
			r, length = decoder(str)
			str = str[length:]
			if r == utf8.RuneError {
				break
			}
			nextProperty = sentenceBreakCodePoints.search(r)
		}
		if nextProperty == prLower {
			return sbLower, false
		}
	}

	return
}
