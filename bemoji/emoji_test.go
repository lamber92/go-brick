package bemoji_test

import (
	"go-brick/bemoji"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDataGroup = []string{
	"ð©",
	"â¢ï¸æççä¼è°¢",
	"è¿æ¬ä¹¦ð¶ï¸ä¸äºé®é¢",
	"This Package emoji is designed to recognize and parse every indivisual Unicode Emoji characters from a string. \nè¿ä¸ªåè¡¨æç¬¦å·æ¨å¨è¯å«åè§£æå­ç¬¦ä¸²ä¸­çæ¯ä¸ªåç¬ç Unicode è¡¨æç¬¦å·å­ç¬¦",
	"1234567890-=!@#$%^&*()_+,./<>?;'[]:\"{}`~",
	"11ï¸â£1111",
	"O(â©_â©)Oåå~",
	"ð©âð©âð¦ð¨ð³",
	"æ¯åï¼ð¢ï¸",
}

func TestHasEmoji(t *testing.T) {
	testDataGroupResult := []bool{
		true,
		true,
		true,
		false,
		false,
		true,
		false,
		true,
		true,
	}
	for i, v := range testDataGroup {
		assert.Equal(t, testDataGroupResult[i], bemoji.HasEmoji(v), "Expected results do not match actual results. [%v]", v)
	}
}

func TestFindEmojiPrefix(t *testing.T) {
	type Result struct {
		OK    bool
		Emoji []rune
	}
	testDataGroupResult := []*Result{
		{true, []rune("ð©")},
		{true, []rune("â¢ï¸")},
		{true, []rune("ð¶ï¸")}, // There is no way to fully match this expression
		{false, []rune("")},
		{false, []rune("")},
		{true, []rune("1ï¸â£")},
		{false, []rune("")},
		{true, []rune("ð©â")}, // There is no way to fully match this expression
		{true, []rune("ð¢ï¸")},
	}
	for i, v := range testDataGroup {
		emoji, ok := bemoji.FindEmojiPrefix(v)
		expected := testDataGroupResult[i]
		assert.Equal(t, expected.OK, ok, "Expected results do not match actual results. [%v]", v)
		t.Logf("expected: %x, actual: %x", expected.Emoji, emoji)
	}
}
