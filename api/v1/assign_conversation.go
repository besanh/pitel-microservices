package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type AssignConversation struct {
	assignConversation service.IAssignConversation
}

func NewAssignConversation(engine *gin.Engine, assignConversationService service.IAssignConversation) {
	handler := &AssignConversation{
		assignConversation: assignConversationService,
	}
	Group := engine.Group("bss-message/v1/assign-conversation")
	{
		Group.GET("user-assigned/:id", handler.GetUserAssigned)
		Group.GET("user-in-queue", handler.GetUserInQueue)
		Group.POST("user-in-queue", handler.InsertUserInQueue)
	}
}

func (handler *AssignConversation) GetUserInQueue(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	filter := model.UserInQueueFilter{
		AppId:            c.Query("app_id"),
		OaId:             c.Query("oa_id"),
		ConversationId:   c.Query("conversation_id"),
		ConversationType: c.Query("conversation_type"),
		Status:           c.Query("status"),
	}

	code, result := handler.assignConversation.GetUserInQueue(c, res.Data, filter)
	c.JSON(code, result)
}

func (handler *AssignConversation) GetUserAssigned(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	conversationId := c.Param("id")
	status := c.Query("status")

	code, result := handler.assignConversation.GetUserAssigned(c, res.Data, conversationId, status)
	c.JSON(code, result)
}

func (handler *AssignConversation) InsertUserInQueue(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	var data model.AssignConversation
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	if err := data.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	code, result := handler.assignConversation.AllocateConversation(c, res.Data, &data)

	c.JSON(code, result)
}
