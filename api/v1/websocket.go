package v1

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

type WebSocket struct {
	subscriber service.ISubscriber
}

var (
	WebSocketHandler *WebSocket
)

func NewWebSocket(engine *gin.Engine, subscriberService service.ISubscriber) {
	handler := &WebSocket{
		subscriber: subscriberService,
	}
	service.WsSubscribers = &service.Subscribers{
		SubscriberMessageBuffer: 16,
		Subscribers:             make(map[*service.Subscriber]struct{}),
		PublishLimiter:          rate.NewLimiter(rate.Every(time.Millisecond*100), 100),
	}
	Group := engine.Group("bss-message/v1/wss")
	{
		Group.GET("subscriber", handler.Subscribe)
	}
}

func (handler *WebSocket) Subscribe(ctx *gin.Context) {
	wsCon, err := websocket.Accept(ctx.Writer, ctx.Request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		CompressionMode:    websocket.CompressionContextTakeover,
		OriginPatterns:     service.ORIGIN_LIST,
	})
	if err != nil {
		log.Error(err)
		ctx.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	defer func() {
		if err := wsCon.Close(websocket.StatusInternalError, "close connection error"); err != nil {
			log.Error(err)
		}
	}()

	err = handler.subscribe(ctx, wsCon)
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
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

	s := &service.Subscriber{
		Message: make(chan []byte, service.WsSubscribers.SubscriberMessageBuffer),
		CloseSlow: func() {
			wsCon.Close(websocket.StatusPolicyViolation, "connection too slow")
		},
	}

	bssAuthRequest := model.BssAuthRequest{
		Token:   c.Query("token"),
		AuthUrl: c.Query("auth_url"),
		Source:  c.Query("source"),
	}

	res := api.AAAMiddleware(c, bssAuthRequest)
	if res == nil {
		return errors.New("token is invalid")
	}

	if err := handler.subscriber.AddSubscriber(c, res.Data, s); err != nil {
		log.Error(err)
		return err
	}
	defer service.WsSubscribers.DeleteSubscriber(s)

	// go func(conn *websocket.Conn) {
	// 	defer func() {
	// 		if err := recover(); err != nil {
	// 			log.Error(err)
	// 		}
	// 	}()

	// 	for {
	// 		messageType, p, err := conn.Read(c)
	// 		if err != nil {
	// 			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
	// 				websocket.CloseStatus(err) == websocket.StatusGoingAway {
	// 				// writeActionChan(callId, CLOSE)
	// 			} else if e, ok := err.(websocket.CloseError); ok {
	// 				log.Errorf("close socket code: %d reason: %s", e.Code, e.Reason)
	// 			} else {
	// 				log.Errorf("read: %v", err)
	// 			}
	// 			break
	// 		}
	// 		log.Infof("read: %s", string(p))
	// 		log.Infof("message type: %d", messageType)
	// 	}

	// }(wsCon)

	ctx := wsCon.CloseRead(c)
	for {
		select {
		case msg := <-s.Message:
			if err := api.WriteTimeout(c, 5*time.Second, wsCon, msg); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
