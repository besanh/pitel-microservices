package v1

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/service"
	"nhooyr.io/websocket"
)

type WebSocket struct {
	subscriber service.ISubscriber
}

var (
	WebSocketHandler *WebSocket
	ORIGIN_LIST      = []string{"localhost:*"}
)

func NewWebSocket(r *gin.Engine, subscriberService service.ISubscriber) {
	WebSocketHandler = &WebSocket{
		subscriber: subscriberService,
	}
	Group := r.Group("bss/v1/chat")
	{
		Group.GET("subscriber", api.MoveTokenToHeader(), WebSocketHandler.Subscribe)
	}
}

func (handler *WebSocket) Subscribe(ctx *gin.Context) {
	wsCon, err := websocket.Accept(ctx.Writer, ctx.Request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		CompressionMode:    websocket.CompressionContextTakeover,
		OriginPatterns:     ORIGIN_LIST,
	})
	if err != nil {
		log.Error(err)
		ctx.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	defer wsCon.Close(websocket.StatusInternalError, "close connection error")

	err = handler.subscribe(ctx, wsCon)
	if err != nil {
		log.Error(err)
		return
	}
}

func (handler *WebSocket) subscribe(c *gin.Context, wsCon *websocket.Conn) error {
	if len(c.Query("source")) < 1 {
		return errors.New("source is required")
	}
	if len(c.Query("token")) < 1 {
		return errors.New("token is required")
	}

	ctx := wsCon.CloseRead(c)
	s := &service.Subscriber{
		Message: make(chan []byte, service.WsSubscribers.SubscriberMessageBuffer),
		CloseSlow: func() {
			wsCon.Close(websocket.StatusPolicyViolation, "connection too slow")
		},
	}
	res := api.AAAMiddleware(c)
	if res == nil {
		return errors.New("token is invalid")
	}
	err := handler.subscriber.AddSubscriber(ctx, res.Data, s)
	if err != nil {
		return err
	}

	defer service.WsSubscribers.DeleteSubscriber(s)
	for {
		select {
		case msg := <-s.Message:
			err := api.WriteTimeout(ctx, 5*time.Second, wsCon, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
