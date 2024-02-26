package service

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/constant"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IShareInfo interface {
		PostConfigForm(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormRequest, files []*multipart.FileHeader) (int, any)
		PostRequestShareInfo(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormRequest) (int, any)
	}
	ShareInfo struct{}
)

func NewShareInfo() IShareInfo {
	return &ShareInfo{}
}

func (s *ShareInfo) PostConfigForm(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormRequest, files []*multipart.FileHeader) (int, any) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return response.OKResponse()
	}
	// Can upload to s3
	dir := constant.PUBLIC_DIR + data.OaId + "/share_info/" + data.Uid
	var filePath string
	for _, file := range files {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		}
		filePath = dir + "/" + util.TimeToStringLayout(time.Now(), "2006_01_02_15_04") + "_" + file.Filename
		var c *gin.Context
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			continue
		}
	}

	filter := model.ShareInfoFormFilter{
		TenantId:  authUser.TenantId,
		ShareType: data.ShareType,
	}

	total, _, err := repository.ShareInfoRepo.GetShareInfos(ctx, dbCon, filter, -1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if total > 0 {
		log.Error("share config " + data.OaId + " already exist")
	}

	shareForm := model.ShareForm{}
	if data.ShareType == "facebook" {
	} else if data.ShareType == "zalo" {
		shareForm.Zalo.ImageUrl = filePath
		shareForm.Zalo.Title = data.Title
		shareForm.Zalo.Subtitle = data.Subtitle
	}

	shareInfoForm := model.ShareInfoForm{
		Base:      model.InitBase(),
		TenantId:  authUser.TenantId,
		ShareType: data.ShareType,
		ShareForm: shareForm,
	}

	if err := repository.ShareInfoRepo.Insert(ctx, dbCon, shareInfoForm); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(
		map[string]any{
			"id": shareInfoForm.GetId(),
		},
	)
}

func (s *ShareInfo) PostRequestShareInfo(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormRequest) (int, any) {
	var result model.OttResponse
	var body any
	tmp := model.OttShareInfoRequest{
		Type:      data.ShareType,
		EventName: data.EventName,
		AppId:     data.AppId,
		OaId:      data.OaId,
		Uid:       data.Uid,
		ImageUrl:  data.ImageUrl,
		Title:     data.Title,
		Subtitle:  data.Subtitle,
	}
	if err := util.ParseAnyToAny(tmp, &body); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	url := OTT_URL + "/ott/v1/crm"
	client := resty.New()

	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		// SetHeader("Authorization", "Bearer "+token).
		SetBody(body).
		Post(url)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	if err := json.Unmarshal([]byte(res.Body()), &result); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if res.StatusCode() == 200 {
		return response.OKResponse()
	} else {
		return response.ServiceUnavailableMsg(result.Message)
	}
}
