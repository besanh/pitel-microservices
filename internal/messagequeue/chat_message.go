package messagequeue

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/internal/rabbitmq"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IChatMessage interface {
		ProduceChatMsgExcute(ctx context.Context, data *model.Message, messageId *uuid.UUID) (err error)
		ConsumeChatMsgExcute(ctx context.Context) (err error)
	}

	ChatMessage struct {
	}
)

var ChatMessageRMQ IChatMessage

func NewChatMessage() IChatMessage {
	return &ChatMessage{}
}

func (c *ChatMessage) ProduceChatMsgExcute(ctx context.Context, data *model.Message, messageId *uuid.UUID) (err error) {

	if data == nil && len(data.Id) > 0 {
		return errors.New("data message is not exist")
	}

	if messageId == nil {
		return errors.New("message id is not exist")
	}

	payload := data

	log.Info("msg payload: ", payload)

	bodyProduce, err := util.ParseStructToByte(payload)
	if err != nil {
		log.Error(err)
		return err
	}

	channelMQ := rabbitmq.ClientMQ.GetChannel()

	if channelMQ.IsClosed() {
		log.Println("ProduceChatMsgExcute channelMQ close connection")
		rabbitmq.ClientMQ.Connect(rabbitmq.ConfigMQ)
		channelMQ = rabbitmq.ClientMQ.GetChannel()
	}

	err = channelMQ.PublishWithContext(ctx,
		"",                          // exchange
		rabbitmq.ConfigMQ.QueueName, // routing key
		false,                       // mandatory
		false,                       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bodyProduce,
		})
	if err != nil {
		log.Error(err)
		return err
	}
	return
}

func (c *ChatMessage) ConsumeChatMsgExcute(ctx context.Context) (err error) {

	channelMQ := rabbitmq.ClientMQ.GetChannel()
	if channelMQ.IsClosed() {
		log.Println("ConsumeChatMsgExcute channelMQ close connection")
		rabbitmq.ClientMQ.Connect(rabbitmq.ConfigMQ)
		channelMQ = rabbitmq.ClientMQ.GetChannel()
	}

	msgs, err := channelMQ.Consume(
		rabbitmq.ConfigMQ.QueueName, // queue
		"",                          // consumer
		false,                       // auto-ack
		false,                       // exclusive
		false,                       // no-local
		false,                       // no-wait
		nil,                         // args
	)
	if err != nil {
		log.Fatal(err)
	}

	for msg := range msgs {
		time.Sleep(time.Millisecond * 500)
		go func(msg amqp.Delivery) {
			log.Info("Consumer from queue --> Received message value:", string(msg.Body))
			log.Println("")
			log.Println("--------------------------------------------------------------")
			log.Println("---------------PROCESSING CONSUME MESSAGE QUEUE---------------")
			log.Println("--------------------------------------------------------------")
			log.Println("")

			var data model.MessageValueInQueue
			if err := util.ParseStringToAny(string(msg.Body), &data); err != nil {
				log.Error("ParseStringToAny message failed:", err)
				return
			}

			// Check Acknowledge message after execute success
			if err := msg.Ack(false); err != nil {
				log.Error("Error acknowledging message:", err)
				return
			}
		}(msg)
	}

	return
}
