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
)

type ChatConnectionApp struct {
	chatConnectionAppService service.IChatConnectionApp
}

func NewChatConnectionApp(engin *gin.Engine, chatConnectionAppService service.IChatConnectionApp, crmAuthUrl string) {
	handler := &ChatConnectionApp{
		chatConnectionAppService: chatConnectionAppService,
	}
	CRM_AUTH_URL = crmAuthUrl
	Group := engin.Group("bss-message/v1/chat-connection-app")
	{
		Group.GET("", handler.GetChatConnectionApp)
		Group.POST("", handler.InsertChatConnectionApp)
		Group.PUT(":id", handler.UpdateChatConnectionAppById)
		Group.GET(":id", handler.GetChatConnectionAppById)
	}
}

func (handler *ChatConnectionApp) GetChatConnectionApp(c *gin.Context) {
	bssAuthRequest := model.BssAuthRequest{
		Token:   c.Query("token"),
		AuthUrl: c.Query("auth_url"),
		Source:  c.Query("source"),
	}

	res := api.AAAMiddleware(c, CRM_AUTH_URL, bssAuthRequest)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}
	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))

	var status sql.NullBool
	if len(c.Query("status")) > 0 {
		statusTmp, _ := strconv.ParseBool(c.Query("status"))
		status.Valid = true
		status.Bool = statusTmp
	}

	filter := model.ChatConnectionAppFilter{
		ConnectionName: c.Query("connection_name"),
		ConnectionType: c.Query("connection_type"),
		Status:         status,
	}

	total, result, err := handler.chatConnectionAppService.GetChatConnectionApp(c, res.Data, filter, limit, offset)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.Pagination(result, total, limit, offset))
}

func (handler *ChatConnectionApp) InsertChatConnectionApp(c *gin.Context) {
	bssAuthRequest := model.BssAuthRequest{
		Token:   c.Query("token"),
		AuthUrl: c.Query("auth_url"),
		Source:  c.Query("source"),
	}

	res := api.AAAMiddleware(c, CRM_AUTH_URL, bssAuthRequest)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	var data model.ChatConnectionAppRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("insert chat connection app payload -> ", &data)

	if err := data.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	id, err := handler.chatConnectionAppService.InsertChatConnectionApp(c, res.Data, data)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.OK(map[string]any{
		"id": id,
	}))
}

func (handler *ChatConnectionApp) GetChatConnectionAppById(c *gin.Context) {
	bssAuthRequest := model.BssAuthRequest{
		Token:   c.Query("token"),
		AuthUrl: c.Query("auth_url"),
		Source:  c.Query("source"),
	}

	res := api.AAAMiddleware(c, CRM_AUTH_URL, bssAuthRequest)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	chatConnectionApp, err := handler.chatConnectionAppService.GetChatConnectionAppById(c, res.Data, id)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.OK(map[string]any{
		"id": chatConnectionApp.Id,
	}))
}

func (handler *ChatConnectionApp) UpdateChatConnectionAppById(c *gin.Context) {
	bssAuthRequest := model.BssAuthRequest{
		Token:   c.Query("token"),
		AuthUrl: c.Query("auth_url"),
		Source:  c.Query("source"),
	}

	res := api.AAAMiddleware(c, CRM_AUTH_URL, bssAuthRequest)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	var data model.ChatConnectionAppRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("update chat connection app payload -> ", data)

	err := handler.chatConnectionAppService.UpdateChatConnectionAppById(c, res.Data, id, data)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}
