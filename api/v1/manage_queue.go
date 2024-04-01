package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type ManageQueue struct {
	manageQueue service.IManageQueue
}

func NewManageQueue(engine *gin.Engine, manageQueue service.IManageQueue) {
	handler := &ManageQueue{
		manageQueue: manageQueue,
	}

	Group := engine.Group("bss-message/v1/manage-queue")
	{
		Group.POST("", handler.PostManageQueue)
		Group.PUT(":id", handler.UpdateManageQueueById)
		Group.DELETE(":id", handler.DeleteManageQueueById)
	}
}

func (handler *ManageQueue) PostManageQueue(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}
	data := model.ChatManageQueueUserRequest{}
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

	log.Info("update manage queue payload -> ", &data)

	id, err := handler.manageQueue.PostManageQueue(c, res.Data, data)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(response.OK(map[string]any{
		"id": id,
	}))
}

func (handler *ManageQueue) UpdateManageQueueById(c *gin.Context) {
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

	data := model.ChatManageQueueUserRequest{}
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

	log.Info("update manage queue payload -> ", &data)

	err := handler.manageQueue.UpdateManageQueueById(c, res.Data, id, data)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(response.OKResponse())
}

func (handler *ManageQueue) DeleteManageQueueById(c *gin.Context) {
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

	err := handler.manageQueue.DeleteManageQueueById(c, res.Data, id)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}
