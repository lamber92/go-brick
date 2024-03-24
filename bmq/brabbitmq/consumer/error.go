package consumer

import (
	"go-brick/berror"
	"go-brick/berror/bcode"
	"go-brick/berror/bstatus"
)

var (
	EventCodeAckFail       = bcode.New(5010001)
	EventCodeNackFail      = bcode.New(5010002)
	EventCodeInfiniteRetry = bcode.New(5010003)
)

var (
	EventAckFail       = berror.New(bstatus.New(EventCodeAckFail, "rabbitmq delivery ACK fail", nil))
	EventNackFail      = berror.New(bstatus.New(EventCodeNackFail, "rabbitmq delivery NACK fail", nil))
	EventInfiniteRetry = berror.New(bstatus.New(EventCodeInfiniteRetry, "rabbitmq delivery infinite retry", nil))
)
