package apollo

import (
	"go-brick/internal/json"

	"github.com/apolloconfig/agollo/v4/storage"
	"go-brick/bconfig/bstorage"
)

type defaultListener struct {
	eventHook bstorage.OnChangeFunc
}

func newDefaultListener(eventHook bstorage.OnChangeFunc) *defaultListener {
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

// OnChange existing config has been modified callback method
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

// OnNewestChange new config has been added callback method
// TODO:
func (listener *defaultListener) OnNewestChange(event *storage.FullChangeEvent) {
	return
}
