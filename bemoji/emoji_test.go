package bemoji_test

import (
	"testing"

	"github.com/lamber92/go-brick/bemoji"
	"github.com/stretchr/testify/assert"
)

var testDataGroup = []string{
	"👩",
	"™️我真的会谢",
	"这本书🈶️一些问题",
	"This Package emoji is designed to recognize and parse every indivisual Unicode Emoji characters from a string. \n这个包表情符号旨在识别和解析字符串中的每个单独的 Unicode 表情符号字符",
	"1234567890-=!@#$%^&*()_+,./<>?;'[]:\"{}`~",
	"11️⃣1111",
	"O(∩_∩)O哈哈~",
	"👩‍👩‍👦🇨🇳",
	"是吗？🛢️",
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
		{true, []rune("👩")},
		{true, []rune("™️")},
		{true, []rune("🈶️")}, // There is no way to fully match this expression
		{false, []rune("")},
		{false, []rune("")},
		{true, []rune("1️⃣")},
		{false, []rune("")},
		{true, []rune("👩‍")}, // There is no way to fully match this expression
		{true, []rune("🛢️")},
	}
	for i, v := range testDataGroup {
		emoji, ok := bemoji.FindEmojiPrefix(v)
		expected := testDataGroupResult[i]
		assert.Equal(t, expected.OK, ok, "Expected results do not match actual results. [%v]", v)
		t.Logf("expected: %x, actual: %x", expected.Emoji, emoji)
	}
}
