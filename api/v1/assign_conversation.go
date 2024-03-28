package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type AssignConversation struct {
	assignConversation service.IAssignConversation
}

func NewAssignConversation(r *gin.Engine, assignConversationService service.IAssignConversation) {
	handler := &AssignConversation{
		assignConversation: assignConversationService,
	}
	Group := r.Group("bss-message/v1/assign-conversation")
	{
		Group.GET("user-in-queue", handler.GetUserInQueue)
	}
}

func (handler *AssignConversation) GetUserInQueue(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	filter := model.UserInQueueFilter{
		AppId:          c.Query("app_id"),
		OaId:           c.Query("oa_id"),
		ConversationId: c.Query("conversation_id"),
	}

	code, result := handler.assignConversation.GetUserInQueue(c, res.Data, filter)
	c.JSON(code, result)
}
