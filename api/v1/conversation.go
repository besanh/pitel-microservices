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
		Group.GET("manager", handler.GetConversationsByManager)
		Group.PUT(":id", handler.UpdateConversation)
		Group.POST("make-done", handler.UpdateMakeDoneConversation)
	}
}

func (handler *Conversation) GetConversations(c *gin.Context) {
	res := api.AuthMiddleware(c)
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
		PhoneNumber:    c.Query("phone_number"),
		Email:          c.Query("email"),
	}

	code, result := handler.conversationService.GetConversations(c, res.Data, filter, limit, offset)
	c.JSON(code, result)
}

func (handler *Conversation) UpdateConversation(c *gin.Context) {
	res := api.AuthMiddleware(c)
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

	code, result := handler.conversationService.UpdateConversationById(c, res.Data, "", id, shareInfo)
	c.JSON(code, result)
}

func (handler *Conversation) UpdateMakeDoneConversation(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	jsonBody := make(map[string]any, 0)
	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	appId, _ := jsonBody["app_id"].(string)
	conversationId, _ := jsonBody["conversation_id"].(string)

	err := handler.conversationService.UpdateMakeDoneConversation(c, res.Data, appId, conversationId, res.Data.UserId)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}

func (handler *Conversation) GetConversationsByManager(c *gin.Context) {
	res := api.AuthMiddleware(c)
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
		PhoneNumber:    c.Query("phone_number"),
		Email:          c.Query("email"),
	}

	code, result := handler.conversationService.GetConversationsByManager(c, res.Data, filter, limit, offset)
	c.JSON(code, result)
}
