package berror_test

import (
	"errors"
	"fmt"
	"go-brick/berror"
	"go-brick/berror/bstatus"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapAndUnwrap(t *testing.T) {
	err1 := errors.New("xxxx")
	err2 := fmt.Errorf("err2. %w", err1)
	err3 := berror.New(bstatus.InternalError, err2)
	err4 := berror.New(bstatus.NotFound, err3)

	targetErrors := []error{err3, err2, err1, nil}

	var tmp error = err4
	for i := 0; tmp != nil; i++ {
		tmp = errors.Unwrap(tmp)
		assert.Equal(t, targetErrors[i], tmp, i)
	}
}
