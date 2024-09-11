package v1

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/tel4vn/pitel-microservices/api"
	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/common/response"
	"github.com/tel4vn/pitel-microservices/common/util"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/service"
	"strconv"
)

type ChatAutoScript struct {
	chatAutoScriptService service.IChatAutoScript
}

func NewChatAutoScript(engine *gin.Engine, chatAutoScriptService service.IChatAutoScript) {
	handler := ChatAutoScript{
		chatAutoScriptService: chatAutoScriptService,
	}

	group := engine.Group("bss-message/v1/chat-auto-script")
	{
		group.GET("", handler.GetChatAutoScripts)
		group.GET(":id", handler.GetChatAutoScriptById)
		group.POST("", handler.InsertChatAutoScript)
		group.PUT(":id", handler.UpdateChatAutoScript)
		group.PUT("status/:id", handler.UpdateChatAutoScriptStatus)
		group.DELETE(":id", handler.DeleteChatAutoScriptById)
	}
}

func (handler *ChatAutoScript) GetChatAutoScripts(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("invalid token"))
		return
	}

	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))

	statusTmp := c.Query("status")
	var status sql.NullBool
	if len(statusTmp) > 0 {
		statusTmp, _ := strconv.ParseBool(statusTmp)
		status.Valid = true
		status.Bool = statusTmp
	}

	filter := model.ChatAutoScriptFilter{
		ScriptName:   c.Query("script_name"),
		Channel:      c.Query("channel"),
		OaId:         c.Query("oa_id"),
		Status:       status,
		TriggerEvent: c.Query("trigger_event"),
	}

	total, result, err := handler.chatAutoScriptService.GetChatAutoScripts(c, res.Data, filter, limit, offset)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.Pagination(result, total, limit, offset))
}

func (handler *ChatAutoScript) GetChatAutoScriptById(c *gin.Context) {
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

	chatScript, err := handler.chatAutoScriptService.GetChatAutoScriptById(c, res.Data, id)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OK(chatScript))
}

func (handler *ChatAutoScript) InsertChatAutoScript(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("invalid token"))
		return
	}

	var chatScriptRequest model.ChatAutoScriptRequest
	err := c.ShouldBindJSON(&chatScriptRequest)
	if err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	if err := chatScriptRequest.Validate(); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	id, err := handler.chatAutoScriptService.InsertChatAutoScript(c, res.Data, chatScriptRequest)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.OK(map[string]any{
		"id": id,
	}))
}

func (handler *ChatAutoScript) UpdateChatAutoScript(c *gin.Context) {
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

	var chatScriptRequest model.ChatAutoScriptRequest
	err := c.ShouldBindJSON(&chatScriptRequest)
	if err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	if err := chatScriptRequest.Validate(); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	err = handler.chatAutoScriptService.UpdateChatAutoScriptById(c, res.Data, id, chatScriptRequest)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	c.JSON(response.OKResponse())
}

func (handler *ChatAutoScript) UpdateChatAutoScriptStatus(c *gin.Context) {
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

	var chatAutoScriptRequest model.ChatAutoScriptStatusRequest
	err := c.ShouldBind(&chatAutoScriptRequest)
	if err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	statusTmp := chatAutoScriptRequest.Status
	var status sql.NullBool
	if len(statusTmp) > 0 {
		statusTmp, _ := strconv.ParseBool(statusTmp)
		status.Valid = true
		status.Bool = statusTmp
	}

	err = handler.chatAutoScriptService.UpdateChatAutoScriptStatusById(c, res.Data, id, status)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	c.JSON(response.OKResponse())
}

func (handler *ChatAutoScript) DeleteChatAutoScriptById(c *gin.Context) {
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

	err := handler.chatAutoScriptService.DeleteChatAutoScriptById(c, res.Data, id)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}
