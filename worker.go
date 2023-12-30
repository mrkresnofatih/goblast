package goblast

import (
	"github.com/google/uuid"
)

type BackgroundWorker struct {
	Subscribers []IBackgroundWorkerSubscriber
}

func (b *BackgroundWorker) Initialize() {
	for _, subscriber := range b.Subscribers {
		subscriber := subscriber
		go func() {
			for {
				msgs, err := subscriber.ReceiveMessages()
				if err != nil {
					LogError(uuid.NewString(), uuid.NewString(), "Failed to receive messages")
					return
				}
				for _, msg := range msgs {
					er := subscriber.HandleMessage(msg)
					if er != nil {
						LogErr(msg.Request.GetMetadata(), "Failed to handle message")
						return
					}
					_ = subscriber.DeleteMessage(msg.Handle)
				}
			}
		}()
	}
}

type IBackgroundWorkerSubscriber interface {
	GetSubscriberName() string
	ReceiveMessages() ([]BackgroundWorkerMessage[interface{}], error)
	HandleMessage(BackgroundWorkerMessage[interface{}]) error
	DeleteMessage(string) error
}

type IBackgroundWorkerPublisher interface {
	PublishMessage(message BackgroundWorkerMessage[interface{}]) error
}

type BackgroundWorkerMessage[T interface{}] struct {
	Topic   string
	Event   string
	Handle  string
	Request ContextfulReq[T]
}
