package messagequeue

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/rabbitmq"
)

var ChannelMQ *amqp.Channel

func NewMQConn(config rabbitmq.Config) (err error) {
	err = rabbitmq.NewRabbitMQ(config)
	if err != nil {
		log.Error(err)
	}
	return
}

func InitMessageQueue() {

}

func HandlerStartConsumerMQ(key string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	switch key {
	case "chat-message-execute":
		messageRMQ := NewChatMessage()
		if err = messageRMQ.ConsumeChatMsgExcute(ctx); err != nil {
			log.Error(err)
			return
		}
	default:
		{
			log.Fatalf("RUN_CONSUMER_FAILED: key \"%s\" is not exist, please check again", key)
			return
		}
	}

	return
}
