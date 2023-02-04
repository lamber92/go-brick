package bcode_test

import (
	"go-brick/berror/bcode"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestDefaultCodeConverter(t *testing.T) {
	grpcCode := bcode.ToGRPCCode(bcode.NotFound)
	assert.Equal(t, codes.NotFound, grpcCode)

	inCode := bcode.FromGRPCCode(grpcCode)
	assert.Equal(t, bcode.NotFound, inCode)

	httpCode := bcode.ToHTTPStatusCode(bcode.NotFound)
	assert.Equal(t, http.StatusNotFound, httpCode)

	inCode = bcode.FromHTTPStatusCode(http.StatusNotFound)
	assert.Equal(t, bcode.NotFound, inCode)
}

func TestExpandDefaultConverter(t *testing.T) {
	var (
		newInCode      = bcode.New(99999)
		grpcCode       = codes.DataLoss
		httpStatusCode = http.StatusExpectationFailed
	)
	bcode.RegisterMapToGRPCCode(newInCode, grpcCode)
	bcode.RegisterMapFromGRPCCode(grpcCode, newInCode)
	bcode.RegisterMapToHTTPStatusCode(newInCode, httpStatusCode)
	bcode.RegisterMapFromHTTPStatusCode(httpStatusCode, newInCode)

	assert.Equal(t, grpcCode, bcode.ToGRPCCode(newInCode))
	assert.Equal(t, newInCode, bcode.FromGRPCCode(grpcCode))
	assert.Equal(t, httpStatusCode, bcode.ToHTTPStatusCode(newInCode))
	assert.Equal(t, newInCode, bcode.FromHTTPStatusCode(httpStatusCode))
}

type myConverter struct {
}

func (d *myConverter) ToGRPCCode(code bcode.Code) codes.Code {
	return codes.Code(code.ToInt())
}

func (d *myConverter) FromGRPCCode(code codes.Code) bcode.Code {
	return NewMyCode(int32(code))
}

func (d *myConverter) ToHTTPStatusCode(code bcode.Code) int {
	return code.ToInt()
}

func (d *myConverter) FromHTTPStatusCode(code int) bcode.Code {
	return NewMyCode(int32(code))
}

func (d *myConverter) RegisterMapToGRPCCode(code bcode.Code, grpcCode codes.Code) {
	return
}

func (d *myConverter) RegisterMapFromGRPCCode(grpcCode codes.Code, code bcode.Code) {
	return
}

func (d *myConverter) RegisterMapToHTTPStatusCode(code bcode.Code, statusCode int) {
	return
}

func (d *myConverter) RegisterMapFromHTTPStatusCode(statusCode int, code bcode.Code) {
	return
}

func TestCustomizedConverter(t *testing.T) {
	bcode.ReplaceCodeConverter(&myConverter{})

	var codeValue int32 = 1000

	xcode := bcode.ToGRPCCode(NewMyCode(codeValue))
	assert.Equal(t, xcode.String(), "Code(1000)")

	ycode := bcode.FromGRPCCode(xcode)
	assert.Equal(t, ycode.ToInt(), int(codeValue))
}
