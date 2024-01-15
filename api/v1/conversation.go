package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type Conversation struct {
	conversationService service.IConversation
}

func NewConversation(engine *gin.Engine, conversationService service.IConversation, crmAuthUrl string) {
	handler := &Conversation{
		conversationService: conversationService,
	}
	Group := engine.Group("bss-message/v1/conversation")
	{
		Group.GET("", func(ctx *gin.Context) {
			handler.GetConversations(ctx, crmAuthUrl)
		})
	}
}

func (handler *Conversation) GetConversations(c *gin.Context, crmAuthUrl string) {
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

	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))

	filter := model.ConversationFilter{}

	code, result := handler.conversationService.GetConversations(c, res.Data, filter, limit, offset)
	c.JSON(code, result)
}
