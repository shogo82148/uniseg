package uniseg

// LineBreakState is the type of the line break parser's states.
type LineBreakState struct {
	notSOT bool
	lb4    int
	lb5    int
}

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

	// LB31.
	// ALL ÷
	// ÷ ALL
	ruleNumber := 310
	lineBreak = LineCanBreak

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

	// LB2: sot ×
	if !state.notSOT {
		ruleNumber = 20
		lineBreak = LineDontBreak
	}

	// LB3: ! eot
	if len(str) <= 0 {
		if ruleNumber > 30 {
			ruleNumber = 30
			lineBreak = LineMustBreak
		}
	}

	// LB4: BK !
	if state.lb4 != 0 {
		if ruleNumber > 40 {
			ruleNumber = 40
			lineBreak = LineDontBreak
		}
	}
	if nextProperty == lbprBK {
		state.lb4 = 1
	} else {
		state.lb4 = 0
	}

	// LB5:
	// CR × LF
	if nextProperty == lbprLF && state.lb5 == 1 {
		if ruleNumber > 50 {
			ruleNumber = 50
			lineBreak = LineDontBreak
		}
	} else if state.lb5 != 0 {
		if ruleNumber > 50 {
			ruleNumber = 50
			lineBreak = LineMustBreak
		}
	}
	if nextProperty == lbprCR {
		state.lb5 = 1
	} else if nextProperty == lbprLF {
		state.lb5 = 2
	} else if nextProperty == lbprNL {
		state.lb5 = 3
	} else {
		state.lb5 = 0
	}

	// LB6: × ( BK | CR | LF | NL )
	if nextProperty == lbprBK || nextProperty == lbprCR || nextProperty == lbprLF || nextProperty == lbprNL {
		if ruleNumber > 60 {
			ruleNumber = 60
			lineBreak = LineDontBreak
		}
	}

	// LB7:
	// × SP
	// × ZW
	if nextProperty == lbprSP || nextProperty == lbprZW {
		if ruleNumber > 70 {
			ruleNumber = 70
			lineBreak = LineDontBreak
		}
	}

	// LB8:
	// ZW SP* ÷

	// LB8a:
	// ZWJ ×

	// LB9:
	// Treat X (CM | ZWJ)* as if it were X.

	// LB10:
	// Treat any remaining CM or ZWJ as it if were AL.

	// LB11:
	// × WJ
	// WJ ×
	if nextProperty == lbprWJ {
		if ruleNumber > 110 {
			ruleNumber = 110
			lineBreak = LineDontBreak
		}
	}

	// LB12:
	// GL ×

	// LB12a:
	// [^SP BA HY] × GL

	// LB13:
	// × CL
	// × CP
	// × EX
	// × IS
	// × SY
	if nextProperty == lbprCL || nextProperty == lbprCP || nextProperty == lbprEX || nextProperty == lbprIS || nextProperty == lbprSY {
		if ruleNumber > 130 {
			ruleNumber = 130
			lineBreak = LineDontBreak
		}
	}

	// LB14:
	// OP SP* ×

	// LB15a:
	// (sot | BK | CR | LF | NL | OP | QU | GL | SP | ZW) [\p{Pi}&QU] SP* ×

	// LB15b:
	// × [\p{Pf}&QU] ( SP | GL | WJ | CL | QU | CP | EX | IS | SY | BK | CR | LF | NL | ZW | eot)

	// LB16:
	// (CL | CP) SP* × NS

	// LB17:
	// B2 SP* × B2

	// LB18:
	// SP ÷

	// LB19:
	// × QU
	// QU ×
	if nextProperty == lbprQU {
		if ruleNumber > 190 {
			ruleNumber = 190
			lineBreak = LineDontBreak
		}
	}

	// LB20:
	// ÷ CB
	// CB ÷
	if nextProperty == lbprCB {
		if ruleNumber > 200 {
			ruleNumber = 200
			lineBreak = LineCanBreak
		}
	}

	// LB21:
	// SY × HL

	// LB22:
	// × IN
	if nextProperty == lbprIN {
		if ruleNumber > 220 {
			ruleNumber = 220
			lineBreak = LineDontBreak
		}
	}

	// LB23:
	// (AL | HL) × NU
	// NU × (AL | HL)

	// LB23a:
	// PR × (ID | EB | EM)
	// (ID | EB | EM) × PO

	// LB24:
	// (PR | PO) × (AL | HL)
	// (AL | HL) × (PR | PO)

	// LB25:
	// CL × PO
	// CP × PO
	// CL × PR
	// CP × PR
	// NU × PO
	// NU × PR
	// PO × OP
	// PO × NU
	// PR × OP
	// PR × NU
	// HY × NU
	// IS × NU
	// NU × NU
	// SY × NU

	// LB26:
	// JL × (JL | JV | H2 | H3)
	// (JV | H2) × (JV | JT)
	// (JT | H3) × JT

	// LB27:
	// (JL | JV | JT | H2 | H3) × PO
	// PR × (JL | JV | JT | H2 | H3)

	// LB28:
	// (AL | HL) × (AL | HL)

	// LB28a:
	// AP × (AK | ◌ | AS)
	// (AK | ◌ | AS) × (VF | VI)
	// (AK | ◌ | AS) VI × (AK | ◌)
	// (AK | ◌ | AS) × (AK | ◌ | AS) VF

	// LB29:
	// IS × (AL | HL)

	// LB30:
	// (AL | HL | NU) × [OP-[\p{ea=F}\p{ea=W}\p{ea=H}]]
	// [CP-[\p{ea=F}\p{ea=W}\p{ea=H}]] × (AL | HL | NU)

	// LB30a:
	// 	sot (RI RI)* RI × RI
	// [^RI] (RI RI)* RI × RI

	// LB30b:
	// 	EB × EM
	// [\p{Extended_Pictographic}&\p{Cn}] × EM

	state.notSOT = true
	return state, lineBreak
}
