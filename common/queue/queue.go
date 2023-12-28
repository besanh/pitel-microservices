package queue

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sys/unix"
)

type Rcfg struct {
	Address  string
	Username string
	Password string
	DB       int
}

var RMQ *RMQConnection

const (
	pollDuration  = 100 * time.Millisecond
	prefetchLimit = 10
)

type RMQConnection struct {
	RedisClient *redis.Client
	Config      Rcfg
	Conn        rmq.Connection
	Queues      map[string]rmq.Queue
	Server      *RMQServer
}

type RMQServer struct {
	conn   rmq.Connection
	Queues map[string]rmq.Queue
}

type RMQClient struct {
	conn rmq.Connection
}

func NewRMQ(config Rcfg) *RMQConnection {
	poolSize := runtime.NumCPU() * 5
	errChan := make(chan error, 10)
	go logErrors(errChan)
	client := redis.NewClient(&redis.Options{
		Addr:         config.Address,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     poolSize,
		PoolTimeout:  time.Duration(20) * time.Second,
		ReadTimeout:  time.Duration(20) * time.Second,
		WriteTimeout: time.Duration(20) * time.Second,
	})
	connection, err := rmq.OpenConnectionWithRedisClient("rmq", client, errChan)
	if err != nil {
		log.Fatal(err)
	}
	return &RMQConnection{
		RedisClient: client,
		Config:      config,
		Conn:        connection,
	}
}
func logErrors(errChan <-chan error) {
	for err := range errChan {
		switch err := err.(type) {
		case *rmq.HeartbeatError:
			if err.Count == rmq.HeartbeatErrorLimit {
				log.Print("heartbeat error (limit): ", err)
			} else {
				log.Print("heartbeat error: ", err)
			}
		case *rmq.ConsumeError:
			log.Print("consume error: ", err)
		case *rmq.DeliveryError:
			log.Print("delivery error: ", err.Delivery, err)
		default:
			log.Print("other error: ", err)
		}
	}
}

func (c *RMQConnection) NewServer() error {
	c.Server = &RMQServer{
		conn:   c.Conn,
		Queues: make(map[string]rmq.Queue),
	}
	return nil
}

func (c *RMQConnection) NewClient() (*RMQClient, error) {
	return &RMQClient{
		conn: c.Conn,
	}, nil
}

var (
	ERR_QUEUE_IS_EXIST     = errors.New("queue is existed")
	ERR_QUEUE_IS_NOT_EXIST = errors.New("queue is not existed")
)

func (srv *RMQServer) AddQueue(name string, handler rmq.ConsumerFunc) error {
	if _, ok := srv.Queues[name]; ok {
		return ERR_QUEUE_IS_EXIST
	}
	queue, err := srv.conn.OpenQueue(name)
	if err != nil {
		return err
	}
	srv.Queues[name] = queue
	if _, err := queue.AddConsumerFunc(name, handler); err != nil {
		return err
	}
	if err := queue.StartConsuming(prefetchLimit, pollDuration); err != nil {
		return err
	}
	return nil
}

func (srv *RMQServer) Close() {
	srv.conn.StopAllConsuming()
}

func (srv *RMQServer) Run() {
	srv.waitForSignals()
}

func (srv *RMQServer) waitForSignals() {
	log.Println("Waiting for signals...")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT, unix.SIGTSTP)
	for {
		sig := <-sigs
		if sig == unix.SIGTSTP {
			srv.conn.StopAllConsuming()
			continue
		}
		break
	}
}

func (srv *RMQServer) RemoveQueue(name string) error {
	if _, ok := srv.Queues[name]; !ok {
		return ERR_QUEUE_IS_NOT_EXIST
	}
	srv.Queues[name].StopConsuming()
	srv.Queues[name].PurgeReady()
	srv.Queues[name].PurgeRejected()
	delete(srv.Queues, name)
	return nil
}

func (c *RMQClient) Publish(queueName string, payload []byte) error {
	queue, err := c.conn.OpenQueue(queueName)
	if err != nil {
		return err
	}
	return queue.PublishBytes(payload)
}
