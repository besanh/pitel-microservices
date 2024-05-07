package v1

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"golang.org/x/exp/slices"
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
		Group.PUT(":app_id/:oa_id/:id", handler.UpdateConversation)
		Group.POST("status", handler.UpdateStatusConversation)
		Group.PATCH(":id/reassign", handler.ReassignConversation)
		Group.GET(":app_id/:oa_id/:id", handler.GetConversationById)
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

	isDone := sql.NullBool{}
	if len(c.Query("is_done")) > 0 {
		isDone.Valid = true
		isDone.Bool, _ = strconv.ParseBool(c.Query("is_done"))
	}

	filter := model.ConversationFilter{
		AppId:          util.ParseQueryArray(c.QueryArray("app_id")),
		ConversationId: util.ParseQueryArray(c.QueryArray("conversation_id")),
		Username:       c.Query("username"),
		PhoneNumber:    c.Query("phone_number"),
		Email:          c.Query("email"),
		IsDone:         isDone,
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

	appId := c.Param("app_id")
	if len(appId) < 1 {
		c.JSON(response.BadRequestMsg("app_id is required"))
	}

	oaId := c.Param("oa_id")
	if len(oaId) < 1 {
		c.JSON(response.BadRequestMsg("oa_id is required"))
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

	code, result := handler.conversationService.UpdateConversationById(c, res.Data, appId, oaId, id, shareInfo)
	c.JSON(code, result)
}

func (handler *Conversation) UpdateStatusConversation(c *gin.Context) {
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
	status, _ := jsonBody["status"].(string)

	log.Info("update status conversation payload -> ", jsonBody)

	if !slices.Contains([]string{"done", "reopen"}, status) {
		c.JSON(response.BadRequestMsg("status is invalid"))
		return
	}

	err := handler.conversationService.UpdateStatusConversation(c, res.Data, appId, conversationId, res.Data.UserId, status)
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
	isDone := sql.NullBool{}
	if len(c.Query("is_done")) > 0 {
		isDone.Valid = true
		isDone.Bool, _ = strconv.ParseBool(c.Query("is_done"))
	}
	filter := model.ConversationFilter{
		AppId:          util.ParseQueryArray(c.QueryArray("app_id")),
		ConversationId: util.ParseQueryArray(c.QueryArray("conversation_id")),
		Username:       c.Query("username"),
		PhoneNumber:    c.Query("phone_number"),
		Email:          c.Query("email"),
		IsDone:         isDone,
	}

	code, result := handler.conversationService.GetConversationsByManage(c, res.Data, filter, limit, offset)
	c.JSON(code, result)
}

func (hanlder *Conversation) ReassignConversation(c *gin.Context) {
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

}

func (handler *Conversation) GetConversationById(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	appId := c.Param("app_id")
	if len(appId) < 1 {
		c.JSON(response.BadRequestMsg("app_id is required"))
		return
	}

	oaId := c.Param("oa_id")
	if len(oaId) < 1 {
		c.JSON(response.BadRequestMsg("oa_id is required"))
		return
	}

	conversationId := c.Param("id")
	if len(conversationId) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	code, result := handler.conversationService.GetConversationById(c, res.Data, appId, conversationId)
	c.JSON(code, result)
}
