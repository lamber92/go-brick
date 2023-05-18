package berror_test

import (
	"errors"
	"fmt"
	"go-brick/berror"
	"go-brick/berror/bcode"
	"go-brick/berror/bstatus"
	"testing"

	xerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	testErr1reason = "xxxx"
	testErr4reason = "yyyy"
	testErr4detail = "jjjj"
)

func generateTestError() (err1, err2, err3, err4 error) {
	err1 = errors.New(testErr1reason)
	err2 = fmt.Errorf("err2. %w", err1)
	err3 = berror.New(bstatus.InternalError, err2)
	err4 = berror.New(bstatus.New(bcode.NotFound, testErr4reason, testErr4detail), err3)
	return
}

func TestDefaultError_Error(t *testing.T) {
	_, _, _, err4 := generateTestError()
	t.Log(err4)
	// {"code":404,"reason":"yyyy","detail":null,"next":{"code":500,"reason":"Internal Server Error","detail":null,"next":"err2. xxxx"}}
}

func TestDefaultError_Status(t *testing.T) {
	_, _, err3, err4 := generateTestError()
	status3 := err3.(berror.Error).Status()
	assert.Equal(t, bstatus.InternalError.Code(), status3.Code())
	assert.Equal(t, bstatus.InternalError.Reason(), status3.Reason())

	status4 := err4.(berror.Error).Status()
	assert.Equal(t, bstatus.NotFound.Code(), status4.Code(), err4)
	assert.Equal(t, testErr4reason, status4.Reason())
	assert.Equal(t, testErr4detail, status4.Detail())
}

func TestDefaultError_Stack(t *testing.T) {
	_, _, err3, err4 := generateTestError()
	t.Log(err3.(berror.Error).Stack())
	t.Log(err4.(berror.Error).Stack())
	// [{"func":"go-brick/berror_test.generateTestError","file":"D:/GitHub/go-brick/berror/error_test.go","line":23},{"func":"go-brick/berror_test.TestDefaultError_Stack","file":"D:/GitHub/go-brick/berror/error_test.go","line":47},{"func":"testing.tRunner","file":"D:/Programs/go1.19.1/go/src/testing/testing.go","line":1446}]
	// [{"func":"go-brick/berror_test.generateTestError","file":"D:/GitHub/go-brick/berror/error_test.go","line":23},{"func":"go-brick/berror_test.TestDefaultError_Stack","file":"D:/GitHub/go-brick/berror/error_test.go","line":47},{"func":"testing.tRunner","file":"D:/Programs/go1.19.1/go/src/testing/testing.go","line":1446}]
}

func TestDefaultError_Unwrap(t *testing.T) {
	err1, err2, err3, err4 := generateTestError()
	targetErrors := []error{err3, err2, err1, nil}
	var tmp = err4
	for i := 0; tmp != nil; i++ {
		tmp = errors.Unwrap(tmp)
		assert.Equal(t, targetErrors[i], tmp, i)
	}
}

func TestDefaultError_Cause(t *testing.T) {
	_, err2, err3, err4 := generateTestError()
	targetErrors := []error{err4, err3, err2}
	for _, tmp := range targetErrors {
		tmp = xerrors.Cause(tmp)
		assert.Equal(t, err2, tmp, tmp.Error())
	}
}

func TestDefaultError_Is(t *testing.T) {
	_, _, err3, err4 := generateTestError()
	assert.ErrorIs(t, err4, err3)
}

func TestDefaultError_As(t *testing.T) {
	err1, err2, err3, err4 := generateTestError()
	assert.ErrorAs(t, err4, &err3)
	assert.ErrorAs(t, err4, &err2)
	assert.ErrorAs(t, err4, &err1)
	assert.Equal(t, err1, err4)
	assert.Equal(t, err2, err4)
	assert.Equal(t, err3, err4)
}
