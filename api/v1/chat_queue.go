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

type ChatQueue struct {
	chatQueueService service.IChatQueue
}

func NewChatQueue(engine *gin.Engine, chatQueueService service.IChatQueue, crmAuthUrl string) {
	handler := &ChatQueue{
		chatQueueService: chatQueueService,
	}
	CRM_AUTH_URL = crmAuthUrl
	Group := engine.Group("bss-message/v1/chat-queue")
	{
		Group.POST("", handler.InsertChatQueue)
		Group.GET("", handler.GetChatQueues)
		Group.GET(":id", handler.GetChatQueueById)
		Group.PUT(":id", handler.UpdateChatQueueById)
		Group.DELETE(":id", handler.DeleteChatQueueById)
	}
}

func (handler *ChatQueue) InsertChatQueue(c *gin.Context) {
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

	var data model.ChatQueueRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("insert chat queue payload -> ", data)

	if err := data.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	id, err := handler.chatQueueService.InsertChatQueue(c, res.Data, data)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.OK(map[string]any{
		"id": id,
	}))
}

func (handler *ChatQueue) GetChatQueues(c *gin.Context) {
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

	filter := model.QueueFilter{
		AppId:     c.Query("app_id"),
		QueueName: c.Query("queue_name"),
	}

	total, chatQueues, err := handler.chatQueueService.GetChatQueues(c, res.Data, filter, limit, offset)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.Pagination(chatQueues, total, limit, offset))
}

func (handler *ChatQueue) GetChatQueueById(c *gin.Context) {
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

	chatQueue, err := handler.chatQueueService.GetChatQueueById(c, res.Data, id)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.OK(chatQueue))
}

func (handler *ChatQueue) UpdateChatQueueById(c *gin.Context) {
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

	var data model.ChatQueueRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("update chat queue payload -> ", data)

	err := handler.chatQueueService.UpdateChatQueueById(c, res.Data, id, data)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}

func (handler *ChatQueue) DeleteChatQueueById(c *gin.Context) {
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

	err := handler.chatQueueService.DeleteChatQueueById(c, res.Data, id)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}
