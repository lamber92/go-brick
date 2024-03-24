package apollo

import (
	"go-brick/bconfig/bstorage"
	"go-brick/btrace"
	"go-brick/internal/json"

	"go.uber.org/zap/zapcore"
)

const (
	traceModule btrace.Module = "apollo_config"
)

func newMetadata(namespace, k string, v bstorage.Value) *defaultMD {
	return &defaultMD{
		ModuleName: traceModule,
		TypeName:   "yaml",
		Namespace:  namespace,
		Key:        k,
		// the value pointed by the pointer may change, here must be a mirror image
		Value: v.String(),
	}
}

type defaultMD struct {
	ModuleName btrace.Module `json:"module"`
	TypeName   string        `json:"type"`
	Namespace  string        `json:"namespace"`
	Key        string        `json:"key"`
	Value      string        `json:"value"`
}

func (m *defaultMD) Module() btrace.Module {
	return m.ModuleName
}

func (m *defaultMD) String() string {
	out, _ := json.MarshalToString(m)
	return out
}

func (m *defaultMD) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("module", string(m.ModuleName))
	enc.AddString("type", m.TypeName)
	enc.AddString("namespace", m.Namespace)
	enc.AddString("key", m.Key)
	enc.AddString("value", m.Value)
	return nil
}
