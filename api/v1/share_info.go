package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
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
	}
}

func (h *ShareInfo) PostConfigForm(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		log.Error("token is invalid")
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err))
		return
	}

	var shareType string
	shareTypeTmp := form.Value["share_type"]
	if len(shareTypeTmp) > 0 {
		shareType = shareTypeTmp[0]
	}
	var appId string
	appIdTmp := form.Value["app_id"]
	if len(appIdTmp) > 0 {
		appId = appIdTmp[0]
	}

	var oaId string
	oaIdTmp := form.Value["oa_id"]
	if len(oaIdTmp) > 0 {
		oaId = oaIdTmp[0]
	}

	var uid string
	uidTmp := form.Value["uid"]
	if len(uidTmp) > 0 {
		uid = uidTmp[0]
	}

	var title string
	titleTmp := form.Value["title"]
	if len(titleTmp) > 0 {
		title = titleTmp[0]
	}

	var subTitle string
	subTitleTmp := form.Value["subtitle"]
	if len(subTitleTmp) > 0 {
		subTitle = subTitleTmp[0]
	}

	data := model.ShareInfoFormRequest{
		ShareType: shareType,
		AppId:     appId,
		OaId:      oaId,
		Uid:       uid,
		Title:     title,
		Subtitle:  subTitle,
	}
	if err := data.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	files := form.File["file"]
	if len(files) < 1 {
		log.Error("file not found")
		c.JSON(response.BadRequestMsg("file not found"))
		return
	}
	code, result := h.shareInfo.PostConfigForm(c, res.Data, data, files)
	c.JSON(code, result)
}

func (s *ShareInfo) PostRequestShareInfo(c *gin.Context) {
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
	log.Info("post share info payload -> ", &data)
	code, result := s.shareInfo.PostRequestShareInfo(c, res.Data, data)
	c.JSON(code, result)
}
