package uniseg

import "testing"

// testTrie checks that the trie returns the same values as lookup for all
// code points. If it fails, run "go generate" to regenerate tries.go.
func testTrie[T comparable](t *testing.T, name string, tr *trie[T], lookup func(rune) T) {
	t.Helper()
	for r := rune(-1); r <= maxRune+1; r++ {
		if got, want := tr.search(r), lookup(r); got != want {
			t.Errorf("%s: search(%U) = %v, want %v", name, r, got, want)
			return
		}
	}
}

func TestTries(t *testing.T) {
	testTrie(t, "graphemeTrie", graphemeTrie, lookupGraphemeProperties)
	testTrie(t, "wordBreakTrie", wordBreakTrie, workBreakCodePoints.search)
	testTrie(t, "sentenceBreakTrie", sentenceBreakTrie, sentenceBreakCodePoints.search)
	testTrie(t, "lineBreakTrie", lineBreakTrie, lineBreakCodePoints.search)
	testTrie(t, "eastAsianWidthTrie", eastAsianWidthTrie, eastAsianWidth.search)
	testTrie(t, "emojiTrie", emojiTrie, emoji.search)
	testTrie(t, "emojiPresentationTrie", emojiPresentationTrie, emojiPresentation.search)
}
