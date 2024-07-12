package v1

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type (
	IAPIChatVendor interface {
		HandlePostChatVendorLogoUpload(c *gin.Context)
		HandlePutChatVendorLogoUpload(c *gin.Context)
	}

	APIChatVendor struct{}
)

var APIChatVendorHandler IAPIChatVendor

func NewChatVendor() IAPIChatVendor {
	return &APIChatVendor{}
}

func (handler *APIChatVendor) HandlePostChatVendorLogoUpload(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		return
	}
	var payload model.ChatVendorRequest
	if err := c.ShouldBind(&payload); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	id, err := service.ChatVendorService.PostChatVendorUpload(c, res.Data, payload, payload.File)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err))
		return
	}
	c.JSON(response.OK(map[string]any{"id": id}))
}

func (handler *APIChatVendor) HandlePutChatVendorLogoUpload(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		return
	}
	id := strings.TrimPrefix(c.Request.RequestURI, "/bss-chat/v1/chat-vendor/upload/")
	if id == "" {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}
	var payload model.ChatVendorRequest
	if err := c.ShouldBind(&payload); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	if payload.File == nil {
		c.JSON(response.BadRequestMsg("Import without file"))
		return
	}
	if err := service.ChatVendorService.PutChatVendorUpload(c, res.Data, id, payload, payload.File); err != nil {
		c.JSON(response.ServiceUnavailableMsg(err))
		return
	}
	c.JSON(response.OK(map[string]any{"id": id}))
}
