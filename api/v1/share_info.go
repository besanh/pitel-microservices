package v1

import (
	"fmt"
	"io"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/internal/storage"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type ShareInfo struct {
	shareInfo service.IShareInfo
}

func NewShareInfo(r *gin.Engine, shareInfo service.IShareInfo) {
	handler := &ShareInfo{
		shareInfo: shareInfo,
	}
	r.MaxMultipartMemory = 10 << 20
	Group := r.Group("bss-message/v1/share-info")
	{
		Group.POST("config", handler.PostConfigForm)
		Group.POST("", handler.PostRequestShareInfo)
		Group.GET("image/:filename", handler.GetImageShareInfo)
	}
}

func (h *ShareInfo) PostConfigForm(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		log.Error("token is invalid")
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	var data model.ShareInfoFormRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	if err := data.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	err := uploadShareInfo(c, data.Files, true)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	err = h.shareInfo.PostConfigForm(c, res.Data, data, data.Files)
	if err != nil {
		err := uploadShareInfo(c, data.Files, false)
		if err != nil {
			c.JSON(response.ServiceUnavailableMsg(err.Error()))
			return
		}
	}
	c.JSON(response.OKResponse())
}

func (s *ShareInfo) PostRequestShareInfo(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		log.Error("token is invalid")
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	var data model.ShareInfoFormSubmitRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	log.Info("post share info payload -> ", &data)
	if err := data.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	err := s.shareInfo.PostRequestShareInfo(c, res.Data, data)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.OKResponse())
}

func uploadShareInfo(c *gin.Context, file *multipart.FileHeader, isOk bool) error {
	f, err := file.Open()
	if err != nil {
		log.Error(err)
		return err
	}
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		log.Error(err)
		return err
	}
	if isOk {
		metaData := storage.NewStoreInput(fileBytes, file.Filename)
		isSuccess, err := storage.Instance.Store(c, *metaData)
		if err != nil || !isSuccess {
			log.Error(err)
			return err
		}
	} else {
		err := storage.Instance.RemoveFile(c, storage.RetrieveInput{
			Path: file.Filename,
		})
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func (h *ShareInfo) GetImageShareInfo(c *gin.Context) {
	fileName := c.Param("filename")
	if len(fileName) < 1 {
		c.JSON(response.BadRequestMsg("filename is required"))
		return
	}
	input := storage.NewRetrieveInput(fileName)
	result, err := storage.Instance.Retrieve(c, *input)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.Writer.Header().Add("Content-Disposition",
		fmt.Sprintf("attachment; filename=%s", fileName))
	c.Writer.Header().Add("Content-Type", c.GetHeader("Content-Type"))
	_, err = c.Writer.Write(result)
	if err != nil {
		log.Error(err)
		c.JSON(response.NotFoundMsg(err))
	}
}
