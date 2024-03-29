package berror_test

import (
	"testing"

	"github.com/lamber92/go-brick/berror"
	"github.com/lamber92/go-brick/berror/bcode"
	"github.com/lamber92/go-brick/berror/bstatus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDefaultConverter_Convert(t *testing.T) {
	// convert and wrap berror.Error
	err1 := berror.NewInvalidArgument(nil, "test err1")
	err2 := berror.Convert(err1, "test err2", err1.Error())

	assert.Equal(t, err1.(berror.Error).Status().Code(), err2.(berror.Error).Status().Code())
	assert.Equal(t, "test err2", err2.(berror.Error).Status().Reason())
	assert.Equal(t, err1.Error(), err2.(berror.Error).Status().Detail())

}

func TestDefaultConverter_Convert2(t *testing.T) {
	err1 := berror.NewInvalidArgument(nil, "test err1")
	err2 := berror.Convert(err1, "test err2", err1.Error())

	assert.ErrorIs(t, err2, err1)
	assert.ErrorAs(t, err2, &err1)
	assert.Equal(t, err1, err2)
}

func TestDefaultConverter_Convert3(t *testing.T) {
	// convert and wrap berror.Error
	err1 := berror.NewInvalidArgument(nil, "test err1")
	err2 := berror.ConvertWithOption(err1, "test err2", nil, berror.IgnoreWrapError())

	assert.Equal(t, err1, err2)
}

func TestDefaultConverter_Convert4(t *testing.T) {
	err1 := gorm.ErrRecordNotFound
	err2 := berror.Convert(err1, "test convert err1")

	assert.Equal(t, bcode.NotFound, err2.(berror.Error).Status().Code())
	assert.Equal(t, "test convert err1", err2.(berror.Error).Status().Reason())

	err3 := gorm.ErrDuplicatedKey
	err4 := berror.Convert(err3, "test convert err3")

	assert.Equal(t, bcode.Unknown, err4.(berror.Error).Status().Code())
}

func TestDefaultConverter_Hook(t *testing.T) {
	berror.RegisterConvHook(func(err error, reason string, detail any, options ...berror.ConvOption) error {
		switch err {
		case gorm.ErrRecordNotFound:
			return berror.NewWithSkip(err, bstatus.New(bcode.InternalError, reason, detail), 1)
		default:
			return err
		}
	})
	err := berror.Convert(gorm.ErrRecordNotFound, "xxx", gorm.ErrRecordNotFound)
	assert.Equal(t, bcode.InternalError, err.(berror.Error).Status().Code())
}
