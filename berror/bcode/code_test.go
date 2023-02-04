package bcode_test

import (
	"go-brick/berror/bcode"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultCode(t *testing.T) {
	errStringOk := bcode.OK.ToString()
	assert.Equal(t, true, bcode.OK.Is(errStringOk))

	errIntUnknown := bcode.Unknown.ToInt()
	assert.Equal(t, true, bcode.Unknown.Is(errIntUnknown))
	assert.Equal(t, false, bcode.Unknown.Is(errStringOk))

	errFloat := 0.00
	assert.Equal(t, false, bcode.OK.Is(errFloat))

	errInt0 := bcode.OK.ToInt()
	errIntOk := bcode.New(errInt0)
	assert.Equal(t, true, errIntOk.Is(bcode.OK))
}

type myCode struct {
	code int32
}

func NewMyCode(code int32) bcode.Code {
	return &myCode{code: code}
}

func (c *myCode) ToInt() int {
	return int(c.code)
}

func (c *myCode) ToString() string {
	return strconv.Itoa(int(c.code))
}

func (c *myCode) Is(target any) bool {
	switch tmp := target.(type) {
	case bcode.Code:
		return c.ToInt() == tmp.ToInt()
	case int:
		return c.ToInt() == tmp
	case int8:
		return c.ToInt() == int(tmp)
	case int32:
		return c.ToInt() == int(tmp)
	case int64:
		return c.ToInt() == int(tmp)
	case uint:
		return c.ToInt() == int(tmp)
	case uint8:
		return c.ToInt() == int(tmp)
	case uint32:
		return c.ToInt() == int(tmp)
	case uint64:
		return c.ToInt() == int(tmp)
	case string:
		return c.ToString() == tmp
	}
	return false
}

func TestCustomizedCode(t *testing.T) {
	var (
		ok      = NewMyCode(0)
		unknown = NewMyCode(-1)
	)

	errStringOk := ok.ToString()
	assert.Equal(t, true, ok.Is(errStringOk))

	errIntUnknown := unknown.ToInt()
	assert.Equal(t, true, unknown.Is(errIntUnknown))
	assert.Equal(t, false, unknown.Is(errStringOk))

	errFloat := 0.00
	assert.Equal(t, false, ok.Is(errFloat))

	errInt0 := ok.ToInt()
	errIntOk := myCode{int32(errInt0)}
	assert.Equal(t, true, errIntOk.Is(ok))

	inIntOk := bcode.New(errInt0)
	assert.Equal(t, true, inIntOk.Is(ok))

	assert.Equal(t, true, bcode.OK.Is(ok))
	assert.Equal(t, true, ok.Is(bcode.OK))

	// because the myCode{}-Is() function receiver is a pointer.
	// so myCode{} is not an implementation of bcode.Code, but &myCode{} is.
	assert.Equal(t, false, bcode.OK.Is(myCode{0}))
	assert.Equal(t, true, bcode.OK.Is(&myCode{0}))
}
