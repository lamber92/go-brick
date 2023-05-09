package btrace

import (
	"context"
	"encoding/hex"
	"go-brick/bconst"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cast"
)

// GenTraceID generate Tracking-ID
func GenTraceID() string {
	return hex.EncodeToString(uuid.NewV4().Bytes())
}

// GetTraceID get Tracking-ID
func GetTraceID(ctx context.Context) string {
	tmp := ctx.Value(bconst.KeyTraceID)
	traceID, ok := tmp.(string)
	if ok {
		return traceID
	}
	return cast.ToString(traceID)
}
