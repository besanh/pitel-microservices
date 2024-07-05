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

type ChatPolicySetting struct {
	chatPolicySettingService service.IChatPolicySetting
}

func NewChatPolicySetting(engine *gin.Engine, chatPolicySettingService service.IChatPolicySetting) {
	handler := ChatPolicySetting{
		chatPolicySettingService: chatPolicySettingService,
	}

	group := engine.Group("bss-message/v1/chat-policy-setting")
	{
		group.GET("", handler.GetChatPolicySettings)
		group.GET(":id", handler.GetChatPolicySettingById)
		group.POST("", handler.InsertChatPolicySetting)
		group.PUT(":id", handler.UpdateChatPolicySetting)
		group.DELETE(":id", handler.DeleteChatPolicySettingById)
	}
}

func (handler *ChatPolicySetting) GetChatPolicySettings(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("invalid token"))
		return
	}

	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))

	filter := model.ChatPolicyFilter{
		ConnectionType: c.Query("connection_type"),
	}

	total, result, err := handler.chatPolicySettingService.GetChatPolicySettings(c, res.Data, filter, limit, offset)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.Pagination(result, total, limit, offset))
}

func (handler *ChatPolicySetting) GetChatPolicySettingById(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("invalid token"))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	chatPolicySetting, err := handler.chatPolicySettingService.GetChatPolicySettingById(c, res.Data, id)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OK(chatPolicySetting))
}

func (handler *ChatPolicySetting) InsertChatPolicySetting(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("invalid token"))
		return
	}

	var request model.ChatPolicyConfigRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	if err := request.Validate(); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	id, err := handler.chatPolicySettingService.InsertChatPolicySetting(c, res.Data, request)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.OK(map[string]any{
		"id": id,
	}))
}

func (handler *ChatPolicySetting) UpdateChatPolicySetting(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("invalid token"))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	var request model.ChatPolicyConfigRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	if err := request.Validate(); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	err = handler.chatPolicySettingService.UpdateChatPolicySettingById(c, res.Data, id, request)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	c.JSON(response.OKResponse())
}

func (handler *ChatPolicySetting) DeleteChatPolicySettingById(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("invalid token"))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	err := handler.chatPolicySettingService.DeleteChatPolicySettingById(c, res.Data, id)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}
