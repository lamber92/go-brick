// Package official
// Declaration:
// This file originates from https://github.com/Andrew-M-C/go.emoji/tree/v1.0.1
// Although there are some changes, the underlying implementation is consistent
package official

// sequences is the collection of Sequence type
type sequences map[rune]*Sequence

// Sequence shows unicode Emoji sequences
type Sequence struct {
	Rune    rune
	End     bool
	Nexts   sequences
	Comment string
}

// newSequence returns a sequence object
func newSequence(r rune) *Sequence {
	return &Sequence{
		Rune:    r,
		End:     false,
		Nexts:   sequences{},
		Comment: "",
	}
}

func init() {
	initSequences()
}

// AllSequences indicates all specified unicode emoji sequences (including single basic emojis)
var AllSequences = sequences{}

// AddSequence add a sequence identified by unicode slice. Notice: this function is NOT goroutine-safe.
func (seq sequences) AddSequence(s []rune, comment string) {
	parentSeq := seq
	total := len(s)
	for i, r := range s {
		node, exist := parentSeq[r]
		if false == exist {
			node = newSequence(r)
			parentSeq[r] = node
		}

		if i == total-1 {
			node.End = true
			node.Comment = comment
		}

		parentSeq = node.Nexts
	}

	return
}

// HasEmoji Check if emoji exists in the string
func (seq sequences) HasEmoji(s string) bool {
	r := []rune(s[:])
	for i := range r {
		// Substring-by-substring traversal inspection
		if seq.checkSub(r[i:]) {
			return true
		}
		// If there is no match, continue to traverse and check backwards
	}
	return false
}

// checkSub Check if substring matches remaining emoji encoding
func (seq sequences) checkSub(r []rune) bool {
	if len(r) == 0 {
		return false
	}
	if sub, exist := seq[r[0]]; exist {
		if sub.End {
			// Proof of match is emoji
			return true
		}

		if len(r) == 1 {
			// There is no next character in the substring to be checked, and the expression library has not been matched to prove that it is not an emoji
			return false
		}
		return sub.Nexts.checkSub(r[1:])
	}
	// Not in the expression library to prove that it is not an emoji
	return false
}

// FindEmojiPrefix Find and return the first emoji if there is an emoji in the string
func (seq sequences) FindEmojiPrefix(s string) ([]rune, bool) {
	r := []rune(s[:])
	for i := range r {
		// Substring-by-substring traversal inspection
		offset := 0
		if seq.markEmoji(r[i:], &offset) {
			return r[i : i+offset], true
		}
		// If there is no match, continue to traverse and check backwards
	}
	return []rune{}, false
}

// markEmoji mark offset if substring matches remaining emoji encoding
func (seq sequences) markEmoji(r []rune, offset *int) bool {
	if len(r) == 0 || offset == nil {
		return false
	}
	if sub, exist := seq[r[0]]; exist {
		*offset++
		if sub.End {
			//if sub.Nexts == nil {
			//	// Proof of match is emoji
			//	return true
			//}
			//// If there are still extensions, continue to check, but all have been found
			//sub.Nexts.markEmoji(r[1:], offset)
			return true
		}

		if len(r) == 1 {
			// There is no next character in the substring to be checked, and the expression library has not been matched to prove that it is not an emoji
			return false
		}
		return sub.Nexts.markEmoji(r[1:], offset)
	}
	// Not in the expression library to prove that it is not an emoji
	return false
}
