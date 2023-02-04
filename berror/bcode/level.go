package bcode

// Level
// error code level, non-essential function, so it is independent and used when needed
// for example, it can be used to judge the log or alarm level according to the error level
type Level int8

const (
	LvInfo Level = 0x01 << iota
	LvNotice
	LvWarning
	LvCritical
)

// codeToLevel error code to level default mapping relationship
var codeToLevel = map[Code]Level{
	Unknown:            LvCritical,
	OK:                 LvInfo,
	InvalidArgument:    LvNotice,
	Unauthorized:       LvNotice,
	Forbidden:          LvWarning,
	NotFound:           LvNotice,
	RequestTimeout:     LvWarning,
	ClientClosed:       LvWarning,
	InternalError:      LvCritical,
	ServiceUnavailable: LvCritical,
	GatewayTimeout:     LvWarning,
	AlreadyExists:      LvNotice,
}

var defLevel = LvCritical

// GetLevel get error code level
func GetLevel(code Code) Level {
	i, ok := code.(interface {
		GetLevel() Level
	})
	if !ok {
		return 0
	}
	return i.GetLevel()
}

// GetLevel get error code level
func (c defaultCode) GetLevel() Level {
	if v, ok := codeToLevel[c]; ok {
		return v
	}
	return defLevel
}

// ReplaceCodeLevelMapping Replace the default code-level-mapping with a custom one
func ReplaceCodeLevelMapping(m map[Code]Level, defaultLevel Level) {
	codeToLevel = m
	defLevel = defaultLevel
}
