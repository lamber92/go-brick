package btrace

import (
	"context"
	"sync"

	"github.com/lamber92/go-brick/bcontext"
	"github.com/lamber92/go-brick/internal/bufferpool"
	bsync "github.com/lamber92/go-brick/internal/sync"
	"go.uber.org/zap/zapcore"
)

type Module string

type Chain interface {
	// Append append metadata to the end of the chain
	Append(...Metadata)
	// Get get metadata chain
	Get() MetadataList
	// Clear clear the metadata chain
	Clear()
	// String formatted output
	String() string
}

type Metadata interface {
	// Module return the name of the module to which the metadata belongs
	Module() Module
	// String formatted output
	String() string
}

type MetadataList []Metadata

func (mdl MetadataList) MarshalLogArray(enc zapcore.ArrayEncoder) (err error) {
	for _, s := range mdl {
		if x, ok := s.(interface {
			MarshalLogObject(zapcore.ObjectEncoder) error
		}); ok {
			_ = enc.AppendObject(x)
			continue
		}
		_ = enc.AppendReflected(s)
	}
	return
}

func NewChain() Chain {
	return &defaultChain{
		chain:  make([]Metadata, 0),
		Locker: bsync.NewSpinLock(),
	}
}

func AppendMDIntoCtx(ctx context.Context, md Metadata) bool {
	switch tmp := ctx.(type) {
	case bcontext.Context:
		var chain Chain
		if ptr, ok := tmp.Get(bcontext.TraceChain); !ok || ptr == nil {
			chain = NewChain()
			chain.Append(md)
		} else {
			if chain, ok = ptr.(Chain); !ok {
				chain = NewChain()
				chain.Append(md)
			} else {
				chain.Append(md)
			}
		}
		tmp.Set(bcontext.TraceChain, chain)
		return true
	default:
		// do nothing...
		return false
	}
}

func GetMDFromCtx(ctx context.Context) (chain Chain, ok bool) {
	ptr := ctx.Value(bcontext.TraceChain)
	if ptr == nil {
		return
	}
	chain, ok = ptr.(Chain)
	return
}

type defaultChain struct {
	chain []Metadata
	sync.Locker
}

func (d *defaultChain) Append(metadata ...Metadata) {
	d.Lock()
	d.chain = append(d.chain, metadata...)
	d.Unlock()
}

func (d *defaultChain) Get() MetadataList {
	d.Lock()
	out := make([]Metadata, 0, len(d.chain))
	out = append(out, d.chain...)
	d.Unlock()
	return out
}

func (d *defaultChain) Clear() {
	d.Lock()
	d.chain = make([]Metadata, 0)
	d.Unlock()
}

func (d *defaultChain) String() string {
	d.Lock()
	buff := bufferpool.Get()
	for idx, v := range d.chain {
		buff.AppendByte('[')
		buff.AppendInt(int64(idx + 1))
		buff.AppendByte(']')
		buff.AppendByte('{')
		buff.AppendString(v.String())
		buff.AppendByte('}')
		if idx+1 < len(d.chain) {
			buff.AppendString(" --> ")
		}
	}
	out := buff.String()
	buff.Free()
	d.Unlock()
	return out
}

func NewMD(mod Module, val string) Metadata {
	return &defaultMD{
		module: mod,
		value:  val,
	}
}

type defaultMD struct {
	module Module
	value  string
}

func (d *defaultMD) Module() Module {
	return d.module
}

func (d *defaultMD) String() string {
	buff := bufferpool.Get()
	buff.AppendString("module: ")
	buff.AppendString(string(d.module))
	buff.AppendByte(',')
	buff.AppendByte(' ')
	buff.AppendString("value: ")
	buff.AppendString(d.value)
	out := buff.String()
	buff.Free()
	return out
}

func (d *defaultMD) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("module", string(d.module))
	enc.AddString("value", d.value)
	return nil
}
