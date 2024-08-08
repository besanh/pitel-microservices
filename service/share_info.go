package service

import (
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

/**
* Khong co delete form vi zalo khong ho tro
 */
type (
	IShareInfo interface {
		// Insert db
		PostConfigForm(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormRequest, file *multipart.FileHeader) (err error)
		// Send to ott service
		PostRequestShareInfo(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormSubmitRequest) (err error)
		UpdateConfigForm(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormRequest, file *multipart.FileHeader) (err error)
		GetShareInfos(ctx context.Context, authUser *model.AuthUser, filter model.ShareInfoFormFilter, limit, offset int) (total int, result *[]model.ShareInfoForm, err error)
		GetShareInfoById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.ShareInfoForm, err error)
		DeleteShareInfoById(ctx context.Context, authUser *model.AuthUser, id string) (err error)
	}
	ShareInfo struct{}
)

var ShareInfoService IShareInfo

func NewShareInfo() IShareInfo {
	return &ShareInfo{}
}

func (s *ShareInfo) PostConfigForm(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormRequest, file *multipart.FileHeader) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		err = errors.New(response.ERR_EMPTY_CONN)
		return
	}

	filter := model.ShareInfoFormFilter{
		TenantId:  authUser.TenantId,
		ShareType: data.ShareType,
		OaId:      data.OaId,
	}

	total, _, err := repository.ShareInfoRepo.GetShareInfos(ctx, dbCon, filter, -1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if total > 0 {
		err = errors.New("share config app_id " + data.AppId + " already exist")
		log.Error(err)
		return
	}

	shareForm := model.ShareForm{}
	if data.ShareType == "facebook" {
	} else if data.ShareType == "zalo" {
		shareForm.Zalo.AppId = data.AppId
		shareForm.Zalo.ImageName = file.Filename
		shareForm.Zalo.ImageUrl = API_SHARE_INFO_HOST + "/" + file.Filename
		shareForm.Zalo.Title = data.Title
		shareForm.Zalo.Subtitle = data.Subtitle
		shareForm.Zalo.OaId = data.OaId
	}

	shareInfoForm := model.ShareInfoForm{
		Base:         model.InitBase(),
		TenantId:     authUser.TenantId,
		ConnectionId: data.ConnectionId,
		ShareType:    data.ShareType,
		ShareForm:    shareForm,
	}

	if err = repository.ShareInfoRepo.Insert(ctx, dbCon, shareInfoForm); err != nil {
		log.Error(err)
		return
	}

	return
}

/**
* Send to ott service
 */
func (s *ShareInfo) PostRequestShareInfo(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormSubmitRequest) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		err = errors.New(response.ERR_EMPTY_CONN)
		return
	}
	filter := model.ShareInfoFormFilter{
		TenantId:  authUser.TenantId,
		ShareType: data.ShareType,
		AppId:     data.AppId,
		OaId:      data.OaId,
	}
	_, shareInfos, err := repository.ShareInfoRepo.GetShareInfos(ctx, dbCon, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*shareInfos) < 1 {
		err = errors.New("share config app_id " + data.AppId + " not exist")
		log.Error(err)
		return
	}
	var result model.OttResponse
	tmp := model.OttShareInfoRequest{
		Type:      data.ShareType,
		EventName: data.EventName,
		AppId:     data.AppId,
		OaId:      (*shareInfos)[0].ShareForm.Zalo.OaId,
		Uid:       data.ExternalUserId,
		ImageUrl:  (*shareInfos)[0].ShareForm.Zalo.ImageUrl,
		Title:     (*shareInfos)[0].ShareForm.Zalo.Title,
		Subtitle:  (*shareInfos)[0].ShareForm.Zalo.Subtitle,
	}

	log.Info("request share info: ", tmp)
	var body any
	if err = util.ParseAnyToAny(tmp, &body); err != nil {
		log.Error(err)
		return
	}

	url := OTT_URL + "/ott/" + OTT_VERSION + "/crm"
	client := resty.New().
		SetTimeout(1 * time.Minute)

	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
	if err != nil {
		log.Error(err)
		return
	}

	if err = json.Unmarshal([]byte(res.Body()), &result); err != nil {
		log.Error(err)
		return
	}
	if res.StatusCode() != 200 {
		err = errors.New(result.Message)
	}
	return
}

/**
* Share info use image from api
 */
func (s *ShareInfo) UpdateConfigForm(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormRequest, file *multipart.FileHeader) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		err = errors.New(response.ERR_EMPTY_CONN)
		return
	}
	if len(data.Id) < 1 {
		err = errors.New("share config id %s not exist " + data.Id)
		log.Errorf(err.Error())
		return
	}

	shareInfoExist, err := repository.ShareInfoRepo.GetById(ctx, dbCon, data.Id)
	if err != nil {
		log.Error(err)
		return
	} else if shareInfoExist == nil {
		err = errors.New("share config id %s not exist " + data.Id)
		log.Errorf(err.Error())
		return
	}

	if data.ShareType == "facebook" {
	} else if data.ShareType == "zalo" {
		// TODO: upload image
		var url string
		if file != nil {
			imageUrl, errTmp := uploadImageToStorageShareInfo(ctx, file)
			if errTmp != nil {
				log.Error(errTmp)
				return
			}
			url = imageUrl
		} else if len(data.ImageUrl) > 0 {
			url = data.ImageUrl
		}

		if len(data.AppId) > 0 {
			shareInfoExist.ShareForm.Zalo.AppId = data.AppId
		}
		if len(url) > 0 {
			if file != nil {
				shareInfoExist.ShareForm.Zalo.ImageName = file.Filename
			}
			shareInfoExist.ShareForm.Zalo.ImageUrl = url
		}
		if len(data.Title) > 0 {
			shareInfoExist.ShareForm.Zalo.Title = data.Title
		}
		if len(data.Subtitle) > 0 {
			shareInfoExist.ShareForm.Zalo.Subtitle = data.Subtitle
		}
		if len(data.OaId) > 0 {
			shareInfoExist.ShareForm.Zalo.OaId = data.OaId
		}
	}

	if err = repository.ShareInfoRepo.Update(ctx, dbCon, *shareInfoExist); err != nil {
		log.Error(err)
		return
	}

	return
}

func (s *ShareInfo) GetShareInfos(ctx context.Context, authUser *model.AuthUser, filter model.ShareInfoFormFilter, limit, offset int) (total int, result *[]model.ShareInfoForm, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		err = errors.New(response.ERR_EMPTY_CONN)
		return
	}

	total, result, err = repository.ShareInfoRepo.GetShareInfos(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *ShareInfo) GetShareInfoById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.ShareInfoForm, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		err = errors.New(response.ERR_EMPTY_CONN)
		return
	}

	result, err = repository.ShareInfoRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *ShareInfo) DeleteShareInfoById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		err = errors.New(response.ERR_EMPTY_CONN)
		return
	}

	shareInfo, err := repository.ShareInfoRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	} else if shareInfo == nil {
		err = errors.New("share config id %s not exist " + id)
		log.Errorf(err.Error())
		return
	}

	err = repository.ShareInfoRepo.Delete(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	}
	return
}
