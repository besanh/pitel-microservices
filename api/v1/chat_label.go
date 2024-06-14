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

type ChatLabel struct {
	chatLabel service.IChatLabel
}

func NewChatLabel(engine *gin.Engine, chatLabel service.IChatLabel) {
	handler := &ChatLabel{
		chatLabel: chatLabel,
	}

	Group := engine.Group("bss-message/v1/chat-label")
	{
		Group.POST("", handler.PostChatLabel)
		Group.GET("", handler.GetChatLabels)
		Group.GET(":id", handler.GetChatLabelById)
		Group.PUT(":id", handler.PutChatLabelById)
		Group.DELETE(":id", handler.DeleteChatLabelById)
	}
}

func (handler *ChatLabel) PostChatLabel(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	data := model.ChatLabelRequest{}
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	if err := data.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("insert chat label payload -> ", &data)

	id, err := handler.chatLabel.InsertChatLabel(c, res.Data, &data)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.OK(map[string]any{
		"id": id,
	}))
}

func (handler *ChatLabel) GetChatLabels(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	limit, offset := util.ParseLimit(c.Query("limit")), util.ParseOffset(c.Query("offset"))

	filter := model.ChatLabelFilter{
		AppId:     c.Query("app_id"),
		OaId:      c.Query("oa_id"),
		LabelType: c.Query("label_type"),
		LabelName: c.Query("label_name"),
	}

	total, result, err := handler.chatLabel.GetChatLabels(c, res.Data, filter, limit, offset)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.Pagination(result, total, limit, offset))
}

func (handler *ChatLabel) GetChatLabelById(c *gin.Context) {
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

	result, err := handler.chatLabel.GetChatLabelById(c, res.Data, id)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.OK(result))
}

func (handler *ChatLabel) PutChatLabelById(c *gin.Context) {
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

	data := model.ChatLabelRequest{}
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	if err := data.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("update chat label payload -> ", &data)

	err := handler.chatLabel.UpdateChatLabelById(c, res.Data, id, &data)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}

func (handler *ChatLabel) DeleteChatLabelById(c *gin.Context) {
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

	err := handler.chatLabel.DeleteChatLabelById(c, res.Data, id)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}
