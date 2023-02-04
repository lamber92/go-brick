package bemoji

import "go-brick/bemoji/official"

// HasEmoji Check if emoji exists in the string
func HasEmoji(s string) bool {
	return official.AllSequences.HasEmoji(s)
}

// FindEmojiPrefix Find and return the first emoji if there is an emoji in the string
func FindEmojiPrefix(s string) ([]rune, bool) {
	return official.AllSequences.FindEmojiPrefix(s)
}
