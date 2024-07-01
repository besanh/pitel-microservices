package streamclient

import "github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"

type (
	IStreamClient interface {
		NewEnvironment() (env *stream.Environment, err error)
	}

	StreamClient struct {
		Config
	}
	Config struct {
		Host string
		Port int
		User string
		Pass string
	}
)

var RabbitMQStreamClient IStreamClient

func NewStreamClient(cfg Config) IStreamClient {
	return &StreamClient{
		Config: cfg,
	}
}

func (c *StreamClient) NewEnvironment() (env *stream.Environment, err error) {
	// rabbitmq-stream://tel4vn:Tel4vn%40PsWrd%23202399@localhost:5552/
	streamResolver := stream.AddressResolver{
		Host: c.Host,
		Port: c.Port,
	}
	env, err = stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			// SetUri("rabbitmq-stream://tel4vn:Tel4vn%40PsWrd%23202399@localhost:5552/").
			SetHost(c.Host).SetPort(c.Port).SetUser(c.User).SetPassword(c.Pass).
			SetVHost("/").
			SetAddressResolver(streamResolver).
			SetMaxConsumersPerClient(2),
	)
	return
}
