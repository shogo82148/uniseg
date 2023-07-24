package uniseg

import "unicode/utf8"

// LineBreakState is the type of the line break parser's states.
type LineBreakState int

// The states of the line break parser.
const (
	_ LineBreakState = iota // The zero value is reserved for the initial state.
	lbAny
	lbBK
	lbCR
	lbLF
	lbNL
	lbSP
	lbZW
	lbWJ
	lbGL
	lbBA
	lbHY
	lbCL
	lbCP
	lbEX
	lbIS
	lbSY
	lbOP
	lbQU
	lbQUSP
	lbNS
	lbCLCPSP
	lbB2
	lbB2SP
	lbCB
	lbBB
	lbLB21a
	lbHL
	lbAL
	lbNU
	lbPR
	lbEB
	lbIDEM
	lbNUNU
	lbNUSY
	lbNUIS
	lbNUCL
	lbNUCP
	lbPO
	lbJL
	lbJV
	lbJT
	lbH2
	lbH3
	lbOddRI
	lbEvenRI
	lbExtPicCn
	lbMax = iota

	lbZWJBit     LineBreakState = 64
	lbCPeaFWHBit LineBreakState = 128
)

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

type lbStateProperty struct {
	LineBreakState
	lbProperty
}

type lbTransitionResult struct {
	LineBreakState
	boundary   LineBreak
	ruleNumber int
}

// The line break parser's state transitions. It's analogous to grTransitions,
// see comments there for details. Unicode version 15.0.0.
var lbTransitions = map[lbStateProperty]lbTransitionResult{
	// LB4.
	{lbAny, lbprBK}: {lbBK, LineCanBreak, 310},
	{lbBK, lbprXX}:  {lbAny, LineMustBreak, 40},

	// LB5.
	{lbAny, lbprCR}: {lbCR, LineCanBreak, 310},
	{lbAny, lbprLF}: {lbLF, LineCanBreak, 310},
	{lbAny, lbprNL}: {lbNL, LineCanBreak, 310},
	{lbCR, lbprLF}:  {lbLF, LineDontBreak, 50},
	{lbCR, lbprXX}:  {lbAny, LineMustBreak, 50},
	{lbLF, lbprXX}:  {lbAny, LineMustBreak, 50},
	{lbNL, lbprXX}:  {lbAny, LineMustBreak, 50},

	// LB6.
	{lbAny, lbprBK}: {lbBK, LineDontBreak, 60},
	{lbAny, lbprCR}: {lbCR, LineDontBreak, 60},
	{lbAny, lbprLF}: {lbLF, LineDontBreak, 60},
	{lbAny, lbprNL}: {lbNL, LineDontBreak, 60},

	// LB7.
	{lbAny, lbprSP}: {lbSP, LineDontBreak, 70},
	{lbAny, lbprZW}: {lbZW, LineDontBreak, 70},

	// LB8.
	{lbZW, lbprSP}: {lbZW, LineDontBreak, 70},
	{lbZW, lbprXX}: {lbAny, LineCanBreak, 80},

	// LB11.
	{lbAny, lbprWJ}: {lbWJ, LineDontBreak, 110},
	{lbWJ, lbprXX}:  {lbAny, LineDontBreak, 110},

	// LB12.
	{lbAny, lbprGL}: {lbGL, LineCanBreak, 310},
	{lbGL, lbprXX}:  {lbAny, LineDontBreak, 120},

	// LB13 (simple transitions).
	{lbAny, lbprCL}: {lbCL, LineCanBreak, 310},
	{lbAny, lbprCP}: {lbCP, LineCanBreak, 310},
	{lbAny, lbprEX}: {lbEX, LineDontBreak, 130},
	{lbAny, lbprIS}: {lbIS, LineCanBreak, 310},
	{lbAny, lbprSY}: {lbSY, LineCanBreak, 310},

	// LB14.
	{lbAny, lbprOP}: {lbOP, LineCanBreak, 310},
	{lbOP, lbprSP}:  {lbOP, LineDontBreak, 70},
	{lbOP, lbprXX}:  {lbAny, LineDontBreak, 140},

	// LB15.
	{lbQU, lbprSP}:   {lbQUSP, LineDontBreak, 70},
	{lbQU, lbprOP}:   {lbOP, LineDontBreak, 150},
	{lbQUSP, lbprOP}: {lbOP, LineDontBreak, 150},

	// LB16.
	{lbCL, lbprSP}:     {lbCLCPSP, LineDontBreak, 70},
	{lbNUCL, lbprSP}:   {lbCLCPSP, LineDontBreak, 70},
	{lbCP, lbprSP}:     {lbCLCPSP, LineDontBreak, 70},
	{lbNUCP, lbprSP}:   {lbCLCPSP, LineDontBreak, 70},
	{lbCL, lbprNS}:     {lbNS, LineDontBreak, 160},
	{lbNUCL, lbprNS}:   {lbNS, LineDontBreak, 160},
	{lbCP, lbprNS}:     {lbNS, LineDontBreak, 160},
	{lbNUCP, lbprNS}:   {lbNS, LineDontBreak, 160},
	{lbCLCPSP, lbprNS}: {lbNS, LineDontBreak, 160},

	// LB17.
	{lbAny, lbprB2}:  {lbB2, LineCanBreak, 310},
	{lbB2, lbprSP}:   {lbB2SP, LineDontBreak, 70},
	{lbB2, lbprB2}:   {lbB2, LineDontBreak, 170},
	{lbB2SP, lbprB2}: {lbB2, LineDontBreak, 170},

	// LB18.
	{lbSP, lbprXX}:     {lbAny, LineCanBreak, 180},
	{lbQUSP, lbprXX}:   {lbAny, LineCanBreak, 180},
	{lbCLCPSP, lbprXX}: {lbAny, LineCanBreak, 180},
	{lbB2SP, lbprXX}:   {lbAny, LineCanBreak, 180},

	// LB19.
	{lbAny, lbprQU}: {lbQU, LineDontBreak, 190},
	{lbQU, lbprXX}:  {lbAny, LineDontBreak, 190},

	// LB20.
	{lbAny, lbprCB}: {lbCB, LineCanBreak, 200},
	{lbCB, lbprXX}:  {lbAny, LineCanBreak, 200},

	// LB21.
	{lbAny, lbprBA}: {lbBA, LineDontBreak, 210},
	{lbAny, lbprHY}: {lbHY, LineDontBreak, 210},
	{lbAny, lbprNS}: {lbNS, LineDontBreak, 210},
	{lbAny, lbprBB}: {lbBB, LineCanBreak, 310},
	{lbBB, lbprXX}:  {lbAny, LineDontBreak, 210},

	// LB21a.
	{lbAny, lbprHL}:   {lbHL, LineCanBreak, 310},
	{lbHL, lbprHY}:    {lbLB21a, LineDontBreak, 210},
	{lbHL, lbprBA}:    {lbLB21a, LineDontBreak, 210},
	{lbLB21a, lbprXX}: {lbAny, LineDontBreak, 211},

	// LB21b.
	{lbSY, lbprHL}:   {lbHL, LineDontBreak, 212},
	{lbNUSY, lbprHL}: {lbHL, LineDontBreak, 212},

	// LB22.
	{lbAny, lbprIN}: {lbAny, LineDontBreak, 220},

	// LB23.
	{lbAny, lbprAL}:  {lbAL, LineCanBreak, 310},
	{lbAny, lbprNU}:  {lbNU, LineCanBreak, 310},
	{lbAL, lbprNU}:   {lbNU, LineDontBreak, 230},
	{lbHL, lbprNU}:   {lbNU, LineDontBreak, 230},
	{lbNU, lbprAL}:   {lbAL, LineDontBreak, 230},
	{lbNU, lbprHL}:   {lbHL, LineDontBreak, 230},
	{lbNUNU, lbprAL}: {lbAL, LineDontBreak, 230},
	{lbNUNU, lbprHL}: {lbHL, LineDontBreak, 230},

	// LB23a.
	{lbAny, lbprPR}:  {lbPR, LineCanBreak, 310},
	{lbAny, lbprID}:  {lbIDEM, LineCanBreak, 310},
	{lbAny, lbprEB}:  {lbEB, LineCanBreak, 310},
	{lbAny, lbprEM}:  {lbIDEM, LineCanBreak, 310},
	{lbPR, lbprID}:   {lbIDEM, LineDontBreak, 231},
	{lbPR, lbprEB}:   {lbEB, LineDontBreak, 231},
	{lbPR, lbprEM}:   {lbIDEM, LineDontBreak, 231},
	{lbIDEM, lbprPO}: {lbPO, LineDontBreak, 231},
	{lbEB, lbprPO}:   {lbPO, LineDontBreak, 231},

	// LB24.
	{lbAny, lbprPO}: {lbPO, LineCanBreak, 310},
	{lbPR, lbprAL}:  {lbAL, LineDontBreak, 240},
	{lbPR, lbprHL}:  {lbHL, LineDontBreak, 240},
	{lbPO, lbprAL}:  {lbAL, LineDontBreak, 240},
	{lbPO, lbprHL}:  {lbHL, LineDontBreak, 240},
	{lbAL, lbprPR}:  {lbPR, LineDontBreak, 240},
	{lbAL, lbprPO}:  {lbPO, LineDontBreak, 240},
	{lbHL, lbprPR}:  {lbPR, LineDontBreak, 240},
	{lbHL, lbprPO}:  {lbPO, LineDontBreak, 240},

	// LB25 (simple transitions).
	{lbPR, lbprNU}:   {lbNU, LineDontBreak, 250},
	{lbPO, lbprNU}:   {lbNU, LineDontBreak, 250},
	{lbOP, lbprNU}:   {lbNU, LineDontBreak, 250},
	{lbHY, lbprNU}:   {lbNU, LineDontBreak, 250},
	{lbNU, lbprNU}:   {lbNUNU, LineDontBreak, 250},
	{lbNU, lbprSY}:   {lbNUSY, LineDontBreak, 250},
	{lbNU, lbprIS}:   {lbNUIS, LineDontBreak, 250},
	{lbNUNU, lbprNU}: {lbNUNU, LineDontBreak, 250},
	{lbNUNU, lbprSY}: {lbNUSY, LineDontBreak, 250},
	{lbNUNU, lbprIS}: {lbNUIS, LineDontBreak, 250},
	{lbNUSY, lbprNU}: {lbNUNU, LineDontBreak, 250},
	{lbNUSY, lbprSY}: {lbNUSY, LineDontBreak, 250},
	{lbNUSY, lbprIS}: {lbNUIS, LineDontBreak, 250},
	{lbNUIS, lbprNU}: {lbNUNU, LineDontBreak, 250},
	{lbNUIS, lbprSY}: {lbNUSY, LineDontBreak, 250},
	{lbNUIS, lbprIS}: {lbNUIS, LineDontBreak, 250},
	{lbNU, lbprCL}:   {lbNUCL, LineDontBreak, 250},
	{lbNU, lbprCP}:   {lbNUCP, LineDontBreak, 250},
	{lbNUNU, lbprCL}: {lbNUCL, LineDontBreak, 250},
	{lbNUNU, lbprCP}: {lbNUCP, LineDontBreak, 250},
	{lbNUSY, lbprCL}: {lbNUCL, LineDontBreak, 250},
	{lbNUSY, lbprCP}: {lbNUCP, LineDontBreak, 250},
	{lbNUIS, lbprCL}: {lbNUCL, LineDontBreak, 250},
	{lbNUIS, lbprCP}: {lbNUCP, LineDontBreak, 250},
	{lbNU, lbprPO}:   {lbPO, LineDontBreak, 250},
	{lbNUNU, lbprPO}: {lbPO, LineDontBreak, 250},
	{lbNUSY, lbprPO}: {lbPO, LineDontBreak, 250},
	{lbNUIS, lbprPO}: {lbPO, LineDontBreak, 250},
	{lbNUCL, lbprPO}: {lbPO, LineDontBreak, 250},
	{lbNUCP, lbprPO}: {lbPO, LineDontBreak, 250},
	{lbNU, lbprPR}:   {lbPR, LineDontBreak, 250},
	{lbNUNU, lbprPR}: {lbPR, LineDontBreak, 250},
	{lbNUSY, lbprPR}: {lbPR, LineDontBreak, 250},
	{lbNUIS, lbprPR}: {lbPR, LineDontBreak, 250},
	{lbNUCL, lbprPR}: {lbPR, LineDontBreak, 250},
	{lbNUCP, lbprPR}: {lbPR, LineDontBreak, 250},

	// LB26.
	{lbAny, lbprJL}: {lbJL, LineCanBreak, 310},
	{lbAny, lbprJV}: {lbJV, LineCanBreak, 310},
	{lbAny, lbprJT}: {lbJT, LineCanBreak, 310},
	{lbAny, lbprH2}: {lbH2, LineCanBreak, 310},
	{lbAny, lbprH3}: {lbH3, LineCanBreak, 310},
	{lbJL, lbprJL}:  {lbJL, LineDontBreak, 260},
	{lbJL, lbprJV}:  {lbJV, LineDontBreak, 260},
	{lbJL, lbprH2}:  {lbH2, LineDontBreak, 260},
	{lbJL, lbprH3}:  {lbH3, LineDontBreak, 260},
	{lbJV, lbprJV}:  {lbJV, LineDontBreak, 260},
	{lbJV, lbprJT}:  {lbJT, LineDontBreak, 260},
	{lbH2, lbprJV}:  {lbJV, LineDontBreak, 260},
	{lbH2, lbprJT}:  {lbJT, LineDontBreak, 260},
	{lbJT, lbprJT}:  {lbJT, LineDontBreak, 260},
	{lbH3, lbprJT}:  {lbJT, LineDontBreak, 260},

	// LB27.
	{lbJL, lbprPO}: {lbPO, LineDontBreak, 270},
	{lbJV, lbprPO}: {lbPO, LineDontBreak, 270},
	{lbJT, lbprPO}: {lbPO, LineDontBreak, 270},
	{lbH2, lbprPO}: {lbPO, LineDontBreak, 270},
	{lbH3, lbprPO}: {lbPO, LineDontBreak, 270},
	{lbPR, lbprJL}: {lbJL, LineDontBreak, 270},
	{lbPR, lbprJV}: {lbJV, LineDontBreak, 270},
	{lbPR, lbprJT}: {lbJT, LineDontBreak, 270},
	{lbPR, lbprH2}: {lbH2, LineDontBreak, 270},
	{lbPR, lbprH3}: {lbH3, LineDontBreak, 270},

	// LB28.
	{lbAL, lbprAL}: {lbAL, LineDontBreak, 280},
	{lbAL, lbprHL}: {lbHL, LineDontBreak, 280},
	{lbHL, lbprAL}: {lbAL, LineDontBreak, 280},
	{lbHL, lbprHL}: {lbHL, LineDontBreak, 280},

	// LB29.
	{lbIS, lbprAL}:   {lbAL, LineDontBreak, 290},
	{lbIS, lbprHL}:   {lbHL, LineDontBreak, 290},
	{lbNUIS, lbprAL}: {lbAL, LineDontBreak, 290},
	{lbNUIS, lbprHL}: {lbHL, LineDontBreak, 290},
}

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

	// Prepare.
	var forceNoBreak, isCPeaFWH bool
	if state > 0 && state&lbCPeaFWHBit != 0 {
		isCPeaFWH = true // LB30: CP but ea is not F, W, or H.
		state = state &^ lbCPeaFWHBit
	}
	if state > 0 && state&lbZWJBit != 0 {
		state = state &^ lbZWJBit // Extract zero-width joiner bit.
		forceNoBreak = true       // LB8a.
	}

	defer func() {
		// Transition into LB30.
		if newState == lbCP || newState == lbNUCP {
			ea := eastAsianWidth.search(r)
			if ea != eawprF && ea != eawprW && ea != eawprH {
				newState |= lbCPeaFWHBit
			}
		}

		// Override break.
		if forceNoBreak {
			lineBreak = LineDontBreak
		}
	}()

	// LB1.
	if nextProperty == lbprAI || nextProperty == lbprSG || nextProperty == lbprXX {
		nextProperty = lbprAL
	} else if nextProperty == lbprSA {
		if generalCategory == gcMn || generalCategory == gcMc {
			nextProperty = lbprCM
		} else {
			nextProperty = lbprAL
		}
	} else if nextProperty == lbprCJ {
		nextProperty = lbprNS
	}

	// Combining marks.
	if nextProperty == lbprZWJ || nextProperty == lbprCM {
		var bit LineBreakState
		if nextProperty == lbprZWJ {
			bit = lbZWJBit
		}
		mustBreakState := state <= 0 || state == lbBK || state == lbCR || state == lbLF || state == lbNL
		if !mustBreakState && state != lbSP && state != lbZW && state != lbQUSP && state != lbCLCPSP && state != lbB2SP {
			// LB9.
			return state | bit, LineDontBreak
		} else {
			// LB10.
			if mustBreakState {
				return lbAL | bit, LineMustBreak
			}
			return lbAL | bit, LineCanBreak
		}
	}

	// Find the applicable transition in the table.
	var rule int
	transition, ok := lbTransitions[lbStateProperty{state, nextProperty}]
	if ok {
		// We have a specific transition. We'll use it.
		newState, lineBreak, rule = transition.LineBreakState, transition.boundary, transition.ruleNumber
	} else {
		// No specific transition found. Try the less specific ones.
		transAnyProp, okAnyProp := lbTransitions[lbStateProperty{state, lbprXX}]
		transAnyState, okAnyState := lbTransitions[lbStateProperty{lbAny, nextProperty}]
		if okAnyProp && okAnyState {
			// Both apply. We'll use a mix (see comments for grTransitions).
			newState, lineBreak, rule = transAnyState.LineBreakState, transAnyState.boundary, transAnyState.ruleNumber
			if transAnyProp.ruleNumber < transAnyState.ruleNumber {
				lineBreak, rule = transAnyProp.boundary, transAnyProp.ruleNumber
			}
		} else if okAnyProp {
			// We only have a specific state.
			newState, lineBreak, rule = transAnyProp.LineBreakState, transAnyProp.boundary, transAnyProp.ruleNumber
			// This branch will probably never be reached because okAnyState will
			// always be true given the current transition map. But we keep it here
			// for future modifications to the transition map where this may not be
			// true anymore.
		} else if okAnyState {
			// We only have a specific property.
			newState, lineBreak, rule = transAnyState.LineBreakState, transAnyState.boundary, transAnyState.ruleNumber
		} else {
			// No known transition. LB31: ALL ÷ ALL.
			newState, lineBreak, rule = lbAny, LineCanBreak, 310
		}
	}

	// LB12a.
	if rule > 121 &&
		nextProperty == lbprGL &&
		(state != lbSP && state != lbBA && state != lbHY && state != lbLB21a && state != lbQUSP && state != lbCLCPSP && state != lbB2SP) {
		return lbGL, LineDontBreak
	}

	// LB13.
	if rule > 130 && state != lbNU && state != lbNUNU {
		switch nextProperty {
		case lbprCL:
			return lbCL, LineDontBreak
		case lbprCP:
			return lbCP, LineDontBreak
		case lbprIS:
			return lbIS, LineDontBreak
		case lbprSY:
			return lbSY, LineDontBreak
		}
	}

	// LB25 (look ahead).
	if rule > 250 &&
		(state == lbPR || state == lbPO) &&
		nextProperty == lbprOP || nextProperty == lbprHY {
		var r rune
		r, _ = decoder(str)
		if r != utf8.RuneError {
			pr := lineBreakCodePoints.search(r).lbProperty
			if pr == lbprNU {
				return lbNU, LineDontBreak
			}
		}
	}

	// LB30 (part one).
	if rule > 300 {
		if (state == lbAL || state == lbHL || state == lbNU || state == lbNUNU) && nextProperty == lbprOP {
			ea := eastAsianWidth.search(r)
			if ea != eawprF && ea != eawprW && ea != eawprH {
				return lbOP, LineDontBreak
			}
		} else if isCPeaFWH {
			switch nextProperty {
			case lbprAL:
				return lbAL, LineDontBreak
			case lbprHL:
				return lbHL, LineDontBreak
			case lbprNU:
				return lbNU, LineDontBreak
			}
		}
	}

	// LB30a.
	if newState == lbAny && nextProperty == lbprRI {
		if state != lbOddRI && state != lbEvenRI { // Includes state == -1.
			// Transition into the first RI.
			return lbOddRI, lineBreak
		}
		if state == lbOddRI {
			// Don't break pairs of Regional Indicators.
			return lbEvenRI, LineDontBreak
		}
		return lbOddRI, lineBreak
	}

	// LB30b.
	if rule > 302 {
		if nextProperty == lbprEM {
			if state == lbEB || state == lbExtPicCn {
				return lbAny, LineDontBreak
			}
		}
		graphemeProperty := graphemeCodePoints.search(r)
		if graphemeProperty == prExtendedPictographic && generalCategory == gcCn {
			return lbExtPicCn, LineCanBreak
		}
	}

	return
}
