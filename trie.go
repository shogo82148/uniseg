package uniseg

// The number of bits of a code point used for indexing each stage of a [trie].
const (
	trieStage2Bits  = 6
	trieStage3Bits  = 6
	trieStage1Shift = trieStage2Bits + trieStage3Bits
	trieStage2Mask  = 1<<trieStage2Bits - 1
	trieStage3Mask  = 1<<trieStage3Bits - 1
)

// trie is a three-stage lookup table mapping code points to values. It is
// generated from a [dictionary] by "go generate" (see trie_gen_test.go).
//
// A code point r is looked up as follows:
//
//  1. stage1[r >> 12] is the index of a block of stage2.
//  2. That block, indexed by bits 6 to 11 of r, is the index of a block of
//     stage3.
//  3. That block, indexed by bits 0 to 5 of r, is the index into values.
//
// Identical blocks are shared, which keeps the tables small.
type trie[T any] struct {
	stage1 []uint16
	stage2 []uint16
	stage3 []uint8
	values []T
}

// search returns the value associated with the given rune. It returns the zero
// value for runes outside the Unicode code space.
func (t *trie[T]) search(r rune) T {
	if uint32(r) > maxRune {
		var zero T
		return zero
	}
	i := int(t.stage1[r>>trieStage1Shift])<<trieStage2Bits | int(r>>trieStage3Bits)&trieStage2Mask
	i = int(t.stage2[i])<<trieStage3Bits | int(r)&trieStage3Mask
	return t.values[t.stage3[i]]
}

// maxRune is the maximum valid Unicode code point.
const maxRune = 0x10FFFF
