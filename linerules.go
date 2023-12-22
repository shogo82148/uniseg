package uniseg

import (
	"unicode/utf8"
)

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
	lbQUSP // (sot | BK | CR | LF | NL | OP | QU | GL | SP | ZW) [\p{Pi}&QU] SP*
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
	lbAP
	lbAK
	lbAS
	lbVF
	lbVI  // (AK | ◌ | AS) VI
	lbMax = iota

	lbZWJBit          LineBreakState = 64
	lbCPeaFWHBit      LineBreakState = 128
	lbDottedCircleBit LineBreakState = 256
	lb15Bit           LineBreakState = 512
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

type lbTransitionResult struct {
	LineBreakState
	boundary   LineBreak
	ruleNumber int
}

// The line break parser's state transitions. It's analogous to grTransitions,
// see comments there for details. Unicode version 15.0.0.
var lbTransitions = [lbMax * lbprMax]lbTransitionResult{
	// LB4.
	int(lbBK)*lbprMax + int(lbprXX): {lbAny, LineMustBreak, 40},

	// LB5.
	int(lbCR)*lbprMax + int(lbprLF): {lbLF, LineDontBreak, 50},
	int(lbCR)*lbprMax + int(lbprXX): {lbAny, LineMustBreak, 50},
	int(lbLF)*lbprMax + int(lbprXX): {lbAny, LineMustBreak, 50},
	int(lbNL)*lbprMax + int(lbprXX): {lbAny, LineMustBreak, 50},

	// LB6.
	int(lbAny)*lbprMax + int(lbprBK): {lbBK, LineDontBreak, 60},
	int(lbAny)*lbprMax + int(lbprCR): {lbCR, LineDontBreak, 60},
	int(lbAny)*lbprMax + int(lbprLF): {lbLF, LineDontBreak, 60},
	int(lbAny)*lbprMax + int(lbprNL): {lbNL, LineDontBreak, 60},

	// LB7.
	int(lbAny)*lbprMax + int(lbprSP): {lbSP, LineDontBreak, 70},
	int(lbAny)*lbprMax + int(lbprZW): {lbZW, LineDontBreak, 70},

	// LB8.
	int(lbZW)*lbprMax + int(lbprSP): {lbZW, LineDontBreak, 70},
	int(lbZW)*lbprMax + int(lbprXX): {lbAny, LineCanBreak, 80},

	// LB11.
	int(lbAny)*lbprMax + int(lbprWJ): {lbWJ, LineDontBreak, 110},
	int(lbWJ)*lbprMax + int(lbprXX):  {lbAny, LineDontBreak, 110},

	// LB12.
	int(lbAny)*lbprMax + int(lbprGL): {lbGL, LineCanBreak, 310},
	int(lbGL)*lbprMax + int(lbprXX):  {lbAny, LineDontBreak, 120},

	// LB13 (simple transitions).
	int(lbAny)*lbprMax + int(lbprCL): {lbCL, LineCanBreak, 310},
	int(lbAny)*lbprMax + int(lbprCP): {lbCP, LineCanBreak, 310},
	int(lbAny)*lbprMax + int(lbprEX): {lbEX, LineDontBreak, 130},
	int(lbAny)*lbprMax + int(lbprIS): {lbIS, LineCanBreak, 310},
	int(lbAny)*lbprMax + int(lbprSY): {lbSY, LineCanBreak, 310},

	// LB14.
	int(lbAny)*lbprMax + int(lbprOP): {lbOP, LineCanBreak, 310},
	int(lbOP)*lbprMax + int(lbprSP):  {lbOP, LineDontBreak, 70},
	int(lbOP)*lbprMax + int(lbprXX):  {lbAny, LineDontBreak, 140},

	// LB16.
	int(lbCL)*lbprMax + int(lbprSP):     {lbCLCPSP, LineDontBreak, 70},
	int(lbNUCL)*lbprMax + int(lbprSP):   {lbCLCPSP, LineDontBreak, 70},
	int(lbCP)*lbprMax + int(lbprSP):     {lbCLCPSP, LineDontBreak, 70},
	int(lbNUCP)*lbprMax + int(lbprSP):   {lbCLCPSP, LineDontBreak, 70},
	int(lbCL)*lbprMax + int(lbprNS):     {lbNS, LineDontBreak, 160},
	int(lbNUCL)*lbprMax + int(lbprNS):   {lbNS, LineDontBreak, 160},
	int(lbCP)*lbprMax + int(lbprNS):     {lbNS, LineDontBreak, 160},
	int(lbNUCP)*lbprMax + int(lbprNS):   {lbNS, LineDontBreak, 160},
	int(lbCLCPSP)*lbprMax + int(lbprNS): {lbNS, LineDontBreak, 160},

	// LB17.
	int(lbAny)*lbprMax + int(lbprB2):  {lbB2, LineCanBreak, 310},
	int(lbB2)*lbprMax + int(lbprSP):   {lbB2SP, LineDontBreak, 70},
	int(lbB2)*lbprMax + int(lbprB2):   {lbB2, LineDontBreak, 170},
	int(lbB2SP)*lbprMax + int(lbprB2): {lbB2, LineDontBreak, 170},

	// LB18.
	int(lbSP)*lbprMax + int(lbprXX):     {lbAny, LineCanBreak, 180},
	int(lbQUSP)*lbprMax + int(lbprXX):   {lbAny, LineCanBreak, 180},
	int(lbCLCPSP)*lbprMax + int(lbprXX): {lbAny, LineCanBreak, 180},
	int(lbB2SP)*lbprMax + int(lbprXX):   {lbAny, LineCanBreak, 180},

	// LB19.
	int(lbAny)*lbprMax + int(lbprQU): {lbQU, LineDontBreak, 190},
	int(lbQU)*lbprMax + int(lbprXX):  {lbAny, LineDontBreak, 190},

	// LB20.
	int(lbAny)*lbprMax + int(lbprCB): {lbCB, LineCanBreak, 200},
	int(lbCB)*lbprMax + int(lbprXX):  {lbAny, LineCanBreak, 200},

	// LB21.
	int(lbAny)*lbprMax + int(lbprBA): {lbBA, LineDontBreak, 210},
	int(lbAny)*lbprMax + int(lbprHY): {lbHY, LineDontBreak, 210},
	int(lbAny)*lbprMax + int(lbprNS): {lbNS, LineDontBreak, 210},
	int(lbAny)*lbprMax + int(lbprBB): {lbBB, LineCanBreak, 310},
	int(lbBB)*lbprMax + int(lbprXX):  {lbAny, LineDontBreak, 210},

	// LB21a.
	int(lbAny)*lbprMax + int(lbprHL):   {lbHL, LineCanBreak, 310},
	int(lbHL)*lbprMax + int(lbprHY):    {lbLB21a, LineDontBreak, 210},
	int(lbHL)*lbprMax + int(lbprBA):    {lbLB21a, LineDontBreak, 210},
	int(lbLB21a)*lbprMax + int(lbprXX): {lbAny, LineDontBreak, 211},

	// LB21b.
	int(lbSY)*lbprMax + int(lbprHL):   {lbHL, LineDontBreak, 212},
	int(lbNUSY)*lbprMax + int(lbprHL): {lbHL, LineDontBreak, 212},

	// LB22.
	int(lbAny)*lbprMax + int(lbprIN): {lbAny, LineDontBreak, 220},

	// LB23.
	int(lbAny)*lbprMax + int(lbprAL):  {lbAL, LineCanBreak, 310},
	int(lbAny)*lbprMax + int(lbprNU):  {lbNU, LineCanBreak, 310},
	int(lbAL)*lbprMax + int(lbprNU):   {lbNU, LineDontBreak, 230},
	int(lbHL)*lbprMax + int(lbprNU):   {lbNU, LineDontBreak, 230},
	int(lbNU)*lbprMax + int(lbprAL):   {lbAL, LineDontBreak, 230},
	int(lbNU)*lbprMax + int(lbprHL):   {lbHL, LineDontBreak, 230},
	int(lbNUNU)*lbprMax + int(lbprAL): {lbAL, LineDontBreak, 230},
	int(lbNUNU)*lbprMax + int(lbprHL): {lbHL, LineDontBreak, 230},

	// LB23a.
	int(lbAny)*lbprMax + int(lbprPR):  {lbPR, LineCanBreak, 310},
	int(lbAny)*lbprMax + int(lbprID):  {lbIDEM, LineCanBreak, 310},
	int(lbAny)*lbprMax + int(lbprEB):  {lbEB, LineCanBreak, 310},
	int(lbAny)*lbprMax + int(lbprEM):  {lbIDEM, LineCanBreak, 310},
	int(lbPR)*lbprMax + int(lbprID):   {lbIDEM, LineDontBreak, 231},
	int(lbPR)*lbprMax + int(lbprEB):   {lbEB, LineDontBreak, 231},
	int(lbPR)*lbprMax + int(lbprEM):   {lbIDEM, LineDontBreak, 231},
	int(lbIDEM)*lbprMax + int(lbprPO): {lbPO, LineDontBreak, 231},
	int(lbEB)*lbprMax + int(lbprPO):   {lbPO, LineDontBreak, 231},

	// LB24.
	int(lbAny)*lbprMax + int(lbprPO): {lbPO, LineCanBreak, 310},
	int(lbPR)*lbprMax + int(lbprAL):  {lbAL, LineDontBreak, 240},
	int(lbPR)*lbprMax + int(lbprHL):  {lbHL, LineDontBreak, 240},
	int(lbPO)*lbprMax + int(lbprAL):  {lbAL, LineDontBreak, 240},
	int(lbPO)*lbprMax + int(lbprHL):  {lbHL, LineDontBreak, 240},
	int(lbAL)*lbprMax + int(lbprPR):  {lbPR, LineDontBreak, 240},
	int(lbAL)*lbprMax + int(lbprPO):  {lbPO, LineDontBreak, 240},
	int(lbHL)*lbprMax + int(lbprPR):  {lbPR, LineDontBreak, 240},
	int(lbHL)*lbprMax + int(lbprPO):  {lbPO, LineDontBreak, 240},

	// LB25 (simple transitions).
	int(lbPR)*lbprMax + int(lbprNU):   {lbNU, LineDontBreak, 250},
	int(lbPO)*lbprMax + int(lbprNU):   {lbNU, LineDontBreak, 250},
	int(lbOP)*lbprMax + int(lbprNU):   {lbNU, LineDontBreak, 250},
	int(lbHY)*lbprMax + int(lbprNU):   {lbNU, LineDontBreak, 250},
	int(lbNU)*lbprMax + int(lbprNU):   {lbNUNU, LineDontBreak, 250},
	int(lbNU)*lbprMax + int(lbprSY):   {lbNUSY, LineDontBreak, 250},
	int(lbNU)*lbprMax + int(lbprIS):   {lbNUIS, LineDontBreak, 250},
	int(lbNUNU)*lbprMax + int(lbprNU): {lbNUNU, LineDontBreak, 250},
	int(lbNUNU)*lbprMax + int(lbprSY): {lbNUSY, LineDontBreak, 250},
	int(lbNUNU)*lbprMax + int(lbprIS): {lbNUIS, LineDontBreak, 250},
	int(lbNUSY)*lbprMax + int(lbprNU): {lbNUNU, LineDontBreak, 250},
	int(lbNUSY)*lbprMax + int(lbprSY): {lbNUSY, LineDontBreak, 250},
	int(lbNUSY)*lbprMax + int(lbprIS): {lbNUIS, LineDontBreak, 250},
	int(lbNUIS)*lbprMax + int(lbprNU): {lbNUNU, LineDontBreak, 250},
	int(lbNUIS)*lbprMax + int(lbprSY): {lbNUSY, LineDontBreak, 250},
	int(lbNUIS)*lbprMax + int(lbprIS): {lbNUIS, LineDontBreak, 250},
	int(lbNU)*lbprMax + int(lbprCL):   {lbNUCL, LineDontBreak, 250},
	int(lbNU)*lbprMax + int(lbprCP):   {lbNUCP, LineDontBreak, 250},
	int(lbNUNU)*lbprMax + int(lbprCL): {lbNUCL, LineDontBreak, 250},
	int(lbNUNU)*lbprMax + int(lbprCP): {lbNUCP, LineDontBreak, 250},
	int(lbNUSY)*lbprMax + int(lbprCL): {lbNUCL, LineDontBreak, 250},
	int(lbNUSY)*lbprMax + int(lbprCP): {lbNUCP, LineDontBreak, 250},
	int(lbNUIS)*lbprMax + int(lbprCL): {lbNUCL, LineDontBreak, 250},
	int(lbNUIS)*lbprMax + int(lbprCP): {lbNUCP, LineDontBreak, 250},
	int(lbNU)*lbprMax + int(lbprPO):   {lbPO, LineDontBreak, 250},
	int(lbNUNU)*lbprMax + int(lbprPO): {lbPO, LineDontBreak, 250},
	int(lbNUSY)*lbprMax + int(lbprPO): {lbPO, LineDontBreak, 250},
	int(lbNUIS)*lbprMax + int(lbprPO): {lbPO, LineDontBreak, 250},
	int(lbNUCL)*lbprMax + int(lbprPO): {lbPO, LineDontBreak, 250},
	int(lbNUCP)*lbprMax + int(lbprPO): {lbPO, LineDontBreak, 250},
	int(lbNU)*lbprMax + int(lbprPR):   {lbPR, LineDontBreak, 250},
	int(lbNUNU)*lbprMax + int(lbprPR): {lbPR, LineDontBreak, 250},
	int(lbNUSY)*lbprMax + int(lbprPR): {lbPR, LineDontBreak, 250},
	int(lbNUIS)*lbprMax + int(lbprPR): {lbPR, LineDontBreak, 250},
	int(lbNUCL)*lbprMax + int(lbprPR): {lbPR, LineDontBreak, 250},
	int(lbNUCP)*lbprMax + int(lbprPR): {lbPR, LineDontBreak, 250},

	// LB26.
	int(lbAny)*lbprMax + int(lbprJL): {lbJL, LineCanBreak, 310},
	int(lbAny)*lbprMax + int(lbprJV): {lbJV, LineCanBreak, 310},
	int(lbAny)*lbprMax + int(lbprJT): {lbJT, LineCanBreak, 310},
	int(lbAny)*lbprMax + int(lbprH2): {lbH2, LineCanBreak, 310},
	int(lbAny)*lbprMax + int(lbprH3): {lbH3, LineCanBreak, 310},
	int(lbJL)*lbprMax + int(lbprJL):  {lbJL, LineDontBreak, 260},
	int(lbJL)*lbprMax + int(lbprJV):  {lbJV, LineDontBreak, 260},
	int(lbJL)*lbprMax + int(lbprH2):  {lbH2, LineDontBreak, 260},
	int(lbJL)*lbprMax + int(lbprH3):  {lbH3, LineDontBreak, 260},
	int(lbJV)*lbprMax + int(lbprJV):  {lbJV, LineDontBreak, 260},
	int(lbJV)*lbprMax + int(lbprJT):  {lbJT, LineDontBreak, 260},
	int(lbH2)*lbprMax + int(lbprJV):  {lbJV, LineDontBreak, 260},
	int(lbH2)*lbprMax + int(lbprJT):  {lbJT, LineDontBreak, 260},
	int(lbJT)*lbprMax + int(lbprJT):  {lbJT, LineDontBreak, 260},
	int(lbH3)*lbprMax + int(lbprJT):  {lbJT, LineDontBreak, 260},

	// LB27.
	int(lbJL)*lbprMax + int(lbprPO): {lbPO, LineDontBreak, 270},
	int(lbJV)*lbprMax + int(lbprPO): {lbPO, LineDontBreak, 270},
	int(lbJT)*lbprMax + int(lbprPO): {lbPO, LineDontBreak, 270},
	int(lbH2)*lbprMax + int(lbprPO): {lbPO, LineDontBreak, 270},
	int(lbH3)*lbprMax + int(lbprPO): {lbPO, LineDontBreak, 270},
	int(lbPR)*lbprMax + int(lbprJL): {lbJL, LineDontBreak, 270},
	int(lbPR)*lbprMax + int(lbprJV): {lbJV, LineDontBreak, 270},
	int(lbPR)*lbprMax + int(lbprJT): {lbJT, LineDontBreak, 270},
	int(lbPR)*lbprMax + int(lbprH2): {lbH2, LineDontBreak, 270},
	int(lbPR)*lbprMax + int(lbprH3): {lbH3, LineDontBreak, 270},

	// LB28.
	int(lbAL)*lbprMax + int(lbprAL): {lbAL, LineDontBreak, 280},
	int(lbAL)*lbprMax + int(lbprHL): {lbHL, LineDontBreak, 280},
	int(lbHL)*lbprMax + int(lbprAL): {lbAL, LineDontBreak, 280},
	int(lbHL)*lbprMax + int(lbprHL): {lbHL, LineDontBreak, 280},

	// LB28a.
	int(lbAny)*lbprMax + int(lbprAP): {lbAP, LineCanBreak, 310},
	int(lbAny)*lbprMax + int(lbprAK): {lbAK, LineCanBreak, 310},
	int(lbAny)*lbprMax + int(lbprAS): {lbAS, LineCanBreak, 310},
	int(lbAP)*lbprMax + int(lbprAK):  {lbAK, LineDontBreak, 281},
	int(lbAP)*lbprMax + int(lbprAS):  {lbAS, LineDontBreak, 281},
	int(lbAK)*lbprMax + int(lbprVF):  {lbVF, LineDontBreak, 281},
	int(lbAK)*lbprMax + int(lbprVI):  {lbVI, LineDontBreak, 281},
	int(lbAS)*lbprMax + int(lbprVF):  {lbVF, LineDontBreak, 281},
	int(lbAS)*lbprMax + int(lbprVI):  {lbVI, LineDontBreak, 281},
	int(lbVI)*lbprMax + int(lbprAK):  {lbAK, LineDontBreak, 281},

	// LB29.
	int(lbIS)*lbprMax + int(lbprAL):   {lbAL, LineDontBreak, 290},
	int(lbIS)*lbprMax + int(lbprHL):   {lbHL, LineDontBreak, 290},
	int(lbNUIS)*lbprMax + int(lbprAL): {lbAL, LineDontBreak, 290},
	int(lbNUIS)*lbprMax + int(lbprHL): {lbHL, LineDontBreak, 290},
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
	var forceNoBreak, isCPeaFWH, isLB15, isDottedCircle bool
	if state > 0 && state&lbCPeaFWHBit != 0 {
		isCPeaFWH = true // LB30: CP but ea is not F, W, or H.
		state = state &^ lbCPeaFWHBit
	}
	if state > 0 && state&lbZWJBit != 0 {
		state = state &^ lbZWJBit // Extract zero-width joiner bit.
		forceNoBreak = true       // LB8a.
	}
	if state > 0 && state&lb15Bit != 0 {
		state = state &^ lb15Bit // Extract LB15 bit.
		isLB15 = true            // LB15.
	}
	if state > 0 && state&lbDottedCircleBit != 0 {
		state = state &^ lbDottedCircleBit // Extract dotted circle bit.
		isDottedCircle = true              // is Dotted Circle ◌.
	}

	defer func() {
		// Transition into LB30.
		if newState == lbCP || newState == lbNUCP {
			ea := eastAsianWidth.search(r)
			if ea != eawprF && ea != eawprW && ea != eawprH {
				newState |= lbCPeaFWHBit
			}
		}

		// Transition into LB15a.
		// (sot | BK | CR | LF | NL | OP | QU | GL | SP | ZW) [\p{Pi}&QU]
		if (state <= 0 || state == lbBK || state == lbCR || state == lbLF || state == lbNL || state == lbOP || state == lbQU || state == lbGL || state == lbSP || state == lbZW) && generalCategory == gcPi && nextProperty == lbprQU {
			newState |= lb15Bit
		}
		if isLB15 && nextProperty == lbprSP {
			newState |= lb15Bit
		}

		// Transition into LB28a.
		if r == '\u25CC' {
			newState |= lbDottedCircleBit
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
		if isDottedCircle {
			bit |= lbDottedCircleBit
		}
		if isLB15 {
			// LB15a.
			bit |= lb15Bit
			return state | bit, LineDontBreak
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
	transition := lbTransitions[int(state)*lbprMax+int(nextProperty)]
	if transition.ruleNumber > 0 {
		// We have a specific transition. We'll use it.
		newState, lineBreak, rule = transition.LineBreakState, transition.boundary, transition.ruleNumber
	} else {
		// No specific transition found. Try the less specific ones.
		transAnyProp := lbTransitions[int(state)*lbprMax+int(lbprXX)]
		transAnyState := lbTransitions[int(lbAny)*lbprMax+int(nextProperty)]
		if transAnyProp.ruleNumber > 0 && transAnyState.ruleNumber > 0 {
			// Both apply. We'll use a mix (see comments for grTransitions).
			newState, lineBreak, rule = transAnyState.LineBreakState, transAnyState.boundary, transAnyState.ruleNumber
			if transAnyProp.ruleNumber < transAnyState.ruleNumber {
				lineBreak, rule = transAnyProp.boundary, transAnyProp.ruleNumber
			}
		} else if transAnyProp.ruleNumber > 0 {
			// We only have a specific state.
			newState, lineBreak, rule = transAnyProp.LineBreakState, transAnyProp.boundary, transAnyProp.ruleNumber
			// This branch will probably never be reached because okAnyState will
			// always be true given the current transition map. But we keep it here
			// for future modifications to the transition map where this may not be
			// true anymore.
		} else if transAnyState.ruleNumber > 0 {
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

	// LB15a.
	if rule > 150 && isLB15 && state == lbSP {
		return lbAny, LineDontBreak
	}

	// LB15b.
	if rule > 151 && generalCategory == gcPf && nextProperty == lbprQU {
		// ( SP | GL | WJ | CL | QU | CP | EX | IS | SY | BK | CR | LF | NL | ZW | eot)
		var r rune
		if len(str) == 0 {
			return lbQU, LineDontBreak
		}
		r, _ = decoder(str)
		if r != utf8.RuneError {
			pr := lineBreakCodePoints.search(r).lbProperty
			if pr == lbprSP || pr == lbprGL || pr == lbprWJ || pr == lbprCL ||
				pr == lbprQU || pr == lbprCP || pr == lbprEX || pr == lbprIS ||
				pr == lbprSY || pr == lbprBK || pr == lbprCR || pr == lbprLF ||
				pr == lbprNL || pr == lbprZW {

				return lbQU, LineDontBreak
			}
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

	// LB28a.
	if rule > 280 {
		// AP × ◌
		if state == lbAP && r == '\u25CC' {
			return lbAL, LineDontBreak
		}

		// ◌ × (VF | VI)
		if isDottedCircle {
			if nextProperty == lbprVF {
				return lbVF, LineDontBreak
			}
			if nextProperty == lbprVI {
				return lbVI, LineDontBreak
			}
		}

		// (AK | ◌ | AS) VI × ◌
		if state == lbVI && r == '\u25CC' {
			return lbAL, LineDontBreak
		}

		// (AK | ◌ | AS) × (AK | ◌ | AS) VF
		if (state == lbAK || state == lbAS || isDottedCircle) &&
			(nextProperty == lbprAK || r == '\u25CC' || nextProperty == lbprAS) {

			// look ahead
			var r rune
			r, _ = decoder(str)
			if r != utf8.RuneError {
				pr := lineBreakCodePoints.search(r).lbProperty
				if pr == lbprVF {
					if nextProperty == lbprAK {
						return lbAK, LineDontBreak
					} else if nextProperty == lbprAS {
						return lbAS, LineDontBreak
					} else {
						return lbAL, LineDontBreak
					}
				}
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
