package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type Message struct {
	messageService service.IMessage
}

func NewMessage(r *gin.Engine, messageService service.IMessage, crmAuthUrl string) {
	handler := &Message{
		messageService: messageService,
	}

	Group := r.Group("bss-message/v1/message")
	{
		Group.POST("send", api.MoveTokenToHeader(), func(ctx *gin.Context) {
			handler.SendMessage(ctx, crmAuthUrl)
		})
	}
}

func (h *Message) SendMessage(c *gin.Context, crmAuthUrl string) {
	bssAuthRequest := model.BssAuthRequest{
		Token:   c.Query("token"),
		AuthUrl: c.Query("auth_url"),
		Source:  c.Query("source"),
	}
	res := api.AAAMiddleware(c, crmAuthUrl, bssAuthRequest)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	jsonBody := make(map[string]any, 0)
	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	log.Info("send message body: ", jsonBody)

	conversationId, _ := jsonBody["conversation_id"].(string)
	parentMessageId, _ := jsonBody["parent_message_id"].(string)
	userIdByApp, _ := jsonBody["user_id_by_app"].(string)
	content, _ := jsonBody["content"].(string)
	// attachment, _ := jsonBody["attachment"].(string)
	message := model.MessageRequest{
		ConversationId:  conversationId,
		ParentMessageId: parentMessageId,
		UserIdByApp:     userIdByApp,
		Content:         content,
		// Attachments:     attachment,
	}

	code, result := h.messageService.SendMessageToOTT(c, res.Data, message)
	c.JSON(code, result)
}
