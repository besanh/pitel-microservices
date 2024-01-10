package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tel4vn/fins-microservices/common/log"
)

type RabbitMQClient struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

type Config struct {
	Addr         string
	ExchangeName string
	QueueName    string
}

var ConfigMQ Config
var ClientMQ RabbitMQClient

func NewRabbitMQ(config Config) error {
	ConfigMQ = config
	err := ClientMQ.Connect(config)
	if err != nil {
		log.Error("Error NewRabbitMQ: ", err)
		return err
	}
	return nil
}

func (r *RabbitMQClient) Connect(config Config) error {
	conn, err := amqp.Dial(config.Addr)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	queue, err := channel.QueueDeclare(
		config.QueueName, // queue name
		true,             // durable
		false,            // auto delete
		false,            // exclusive
		false,            // no wait
		nil,              // arguments
	)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	ClientMQ = RabbitMQClient{
		Channel: channel,
		Queue:   queue,
	}
	return err
}

func (r *RabbitMQClient) GetChannel() *amqp.Channel {
	return r.Channel
}
