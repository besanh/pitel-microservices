package v1

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type ChatEmail struct {
	chatEmail service.IChatEmail
}

func NewChatEmail(engine *gin.Engine, chatEmail service.IChatEmail) {
	handler := &ChatEmail{
		chatEmail: chatEmail,
	}

	Group := engine.Group("bss-message/v1/chat-email")
	{
		Group.POST("", handler.InsertChatEmail)
		Group.GET("", handler.GetChatEmails)
		Group.GET(":id", handler.GetChatEmailById)
		Group.PUT(":id", handler.UpdateChatEmailById)
		Group.DELETE(":id", handler.DeleteChatEmailById)
	}
}

func (h *ChatEmail) GetChatEmails(c *gin.Context) {
	res := api.AuthMiddleware(c)

	statusTmp := c.Query("status")
	var status sql.NullBool
	if len(statusTmp) > 0 {
		statusTmp, _ := strconv.ParseBool(statusTmp)
		status.Valid = true
		status.Bool = statusTmp
	}

	filter := model.ChatEmailFilter{
		TenantId: c.Query("tenant_id"),
		OaId:     c.Query("oa_id"),
		Status:   status,
	}

	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))

	total, data, err := h.chatEmail.GetChatEmails(c, res.Data, filter, limit, offset)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.Pagination(data, total, limit, offset))
}

func (h *ChatEmail) InsertChatEmail(c *gin.Context) {
	res := api.AuthMiddleware(c)

	var chatEmail model.ChatEmailRequest
	err := c.ShouldBindJSON(&chatEmail)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	if err := chatEmail.Validate(); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	id, err := h.chatEmail.InsertChatEmail(c, res.Data, chatEmail)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OK(map[string]any{
		"id": id,
	}))
}

func (h *ChatEmail) GetChatEmailById(c *gin.Context) {
	res := api.AuthMiddleware(c)

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	chatEmail, err := h.chatEmail.GetChatEmailById(c, res.Data, id)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OK(chatEmail))
}

func (h *ChatEmail) UpdateChatEmailById(c *gin.Context) {
	res := api.AuthMiddleware(c)

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	var chatEmail model.ChatEmailRequest
	err := c.ShouldBindJSON(&chatEmail)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	if err := chatEmail.Validate(); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	err = h.chatEmail.UpdateChatEmailById(c, res.Data, id, chatEmail)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OK(map[string]any{
		"id": id,
	}))
}

func (h *ChatEmail) DeleteChatEmailById(c *gin.Context) {
	res := api.AuthMiddleware(c)

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	err := h.chatEmail.DeleteChatEmailById(c, res.Data, id)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}
