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

type ChatApp struct {
	chatAppService service.IChatApp
}

var CRM_AUTH_URL string

func NewChatApp(engine *gin.Engine, chatAppService service.IChatApp, crmAuthUrl string) {
	handler := &ChatApp{
		chatAppService: chatAppService,
	}
	CRM_AUTH_URL = crmAuthUrl
	Group := engine.Group("bss-message/v1/chat-app")
	{
		Group.POST("", handler.InsertChatApp)
		Group.GET("", handler.GetChatApp)
		Group.GET(":id", handler.GetChatAppById)
		Group.PUT(":id", handler.UpdateChatAppById)
		Group.DELETE(":id", handler.DeleteChatAppById)
	}
}

func (handler *ChatApp) InsertChatApp(c *gin.Context) {
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

	var data model.ChatAppRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("insert chat app payload -> ", &data)

	if err := data.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	id, err := handler.chatAppService.InsertChatApp(c, res.Data, data)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.OK(map[string]any{
		"id": id,
	}))
}

func (handler *ChatApp) GetChatApp(c *gin.Context) {
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

	filter := model.AppFilter{
		AppName: c.Query("app_name"),
		Status:  c.Query("status"),
	}
	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))

	total, chatApps, err := handler.chatAppService.GetChatApp(c, res.Data, filter, limit, offset)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
	}
	c.JSON(response.Pagination(chatApps, total, limit, offset))
}

func (handler *ChatApp) GetChatAppById(c *gin.Context) {
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

	chatApp, err := handler.chatAppService.GetChatAppById(c, res.Data, id)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.OK(map[string]any{
		"id": chatApp.Id,
	}))
}

func (handler *ChatApp) UpdateChatAppById(c *gin.Context) {
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

	var data model.ChatAppRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("update chat app payload -> ", &data)

	err := handler.chatAppService.UpdateChatAppById(c, res.Data, id, data)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}

func (handler *ChatApp) DeleteChatAppById(c *gin.Context) {
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

	err := handler.chatAppService.DeleteChatAppById(c, res.Data, id)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}
