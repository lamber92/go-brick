package consumer

import (
	"go-brick/berror"
	"go-brick/berror/bcode"
	"go-brick/berror/bstatus"
)

var (
	EventCodeAckFail         = bcode.New(5010001)
	EventCodeNackFail        = bcode.New(5010002)
	EventCodeRetryInfinitely = bcode.New(5010003)
)

var (
	EventAckFail         = berror.New(bstatus.New(EventCodeAckFail, "rabbitmq delivery ACK fail", nil))
	EventNackFail        = berror.New(bstatus.New(EventCodeNackFail, "rabbitmq delivery NACK fail", nil))
	EventRetryInfinitely = berror.New(bstatus.New(EventCodeRetryInfinitely, "rabbitmq delivery RETRY infinitely", nil))
)
