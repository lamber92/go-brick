package apollo

import (
	"go-brick/bconfig"
	"go-brick/internal/json"

	"github.com/apolloconfig/agollo/v4/storage"
)

type defaultListener struct {
	eventHook bconfig.OnChangeFunc
}

func newDefaultListener(eventHook bconfig.OnChangeFunc) *defaultListener {
	return &defaultListener{eventHook: eventHook}
}

type ChangeEvent struct {
	Namespace      string                   `json:"Namespace"`
	NotificationID int64                    `json:"NotificationID"`
	Changes        map[string]*ConfigChange `json:"Changes"`
}

type ConfigChange struct {
	ChangeType string      `json:"ChangeType"`
	OldValue   interface{} `json:"OldValue"`
	NewValue   interface{} `json:"NewValue"`
}

var mapChangeTypeToDesc = map[storage.ConfigChangeType]string{
	storage.ADDED:    "ADD",
	storage.MODIFIED: "MODIFY",
	storage.DELETED:  "DELETE",
}

// OnChange 增加变更监控
func (listener *defaultListener) OnChange(event *storage.ChangeEvent) {
	innerEvent := &ChangeEvent{
		Namespace:      event.Namespace,
		NotificationID: event.NotificationID,
		Changes:        make(map[string]*ConfigChange, len(event.Changes)),
	}
	if event.Changes != nil {
		for k, v := range event.Changes {
			innerEvent.Changes[k] = &ConfigChange{
				ChangeType: mapChangeTypeToDesc[v.ChangeType],
				OldValue:   v.OldValue,
				NewValue:   v.NewValue,
			}
		}
	}
	eventSrc, _ := json.MarshalToString(innerEvent)
	listener.eventHook(eventSrc)
}

// OnNewestChange 监控最新变更
func (listener *defaultListener) OnNewestChange(event *storage.FullChangeEvent) {
	return
}
