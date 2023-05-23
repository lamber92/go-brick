package yaml

import (
	"go-brick/bconfig"
	"go-brick/btrace"
	"go-brick/internal/json"

	"go.uber.org/zap/zapcore"
)

const (
	traceModule btrace.Module = "config"
)

func newMetadata(namespace, k string, v bconfig.Value) *defaultMD {
	return &defaultMD{
		ModuleName: traceModule,
		TypeName:   "yaml",
		Namespace:  namespace,
		Key:        k,
		Value:      v,
	}
}

type defaultMD struct {
	ModuleName btrace.Module `json:"module"`
	TypeName   string        `json:"type"`
	Namespace  string        `json:"namespace"`
	Key        string        `json:"key"`
	Value      bconfig.Value `json:"value"`
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
	enc.AddString("value", m.Value.String())
	return nil
}
