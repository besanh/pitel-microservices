package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type Conversation struct {
	conversationService service.IConversation
}

func NewConversation(engine *gin.Engine, conversationService service.IConversation) {
	handler := &Conversation{
		conversationService: conversationService,
	}
	Group := engine.Group("bss-message/v1/conversation")
	{
		Group.GET("", handler.GetConversations)
		Group.PUT(":id", handler.UpdateConversation)
	}
}

func (handler *Conversation) GetConversations(c *gin.Context) {
	bssAuthRequest := model.BssAuthRequest{
		Token:   c.Query("token"),
		AuthUrl: c.Query("auth_url"),
		Source:  c.Query("source"),
	}
	res := api.AAAMiddleware(c, service.CRM_AUTH_URL, bssAuthRequest)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))

	filter := model.ConversationFilter{
		AppId:          util.ParseQueryArray(c.QueryArray("app_id")),
		ConversationId: util.ParseQueryArray(c.QueryArray("conversation_id")),
		Username:       c.Query("username"),
		PhoneNumber:    util.ParseQueryArray(c.QueryArray("phone_number")),
		Email:          util.ParseQueryArray(c.QueryArray("email")),
	}

	code, result := handler.conversationService.GetConversations(c, res.Data, filter, limit, offset)
	c.JSON(code, result)
}

func (handler *Conversation) UpdateConversation(c *gin.Context) {
	bssAuthRequest := model.BssAuthRequest{
		Token:   c.Query("token"),
		AuthUrl: c.Query("auth_url"),
		Source:  c.Query("source"),
	}
	res := api.AAAMiddleware(c, service.CRM_AUTH_URL, bssAuthRequest)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	var shareInfo model.ShareInfo
	if err := c.ShouldBindJSON(&shareInfo); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("update conversation payload -> ", shareInfo)

	code, result := handler.conversationService.UpdateConversationById(c, res.Data, id, shareInfo)
	c.JSON(code, result)
}
