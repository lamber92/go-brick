package brabbitmq

import (
	"context"
	"go-brick/btrace"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// BuildTextMsg4Publish build a simple rabbitmq-producer-message of text
func BuildTextMsg4Publish(ctx context.Context, body []byte, persistent bool, priorities ...uint8) *amqp.Publishing {
	var (
		priority     uint8 = 0
		deliveryMode uint8 = 0
	)
	if len(priorities) > 0 {
		priority = priorities[0]
	}
	if persistent {
		deliveryMode = 2
	}
	return &amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "text/plain",
		ContentEncoding: "",
		Body:            body,
		DeliveryMode:    deliveryMode, // example: amqp.Transient, 1=non-persistent, 2=persistent
		Priority:        priority,     // 0-9
		MessageId:       btrace.GetTraceID(ctx),
		Timestamp:       time.Now(),
	}
}
