package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/internal/storage"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

/**
* Khong co delete form vi zalo khong ho tro
 */
type (
	IShareInfo interface {
		PostConfigForm(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormRequest, file *multipart.FileHeader) error
		PostRequestShareInfo(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormSubmitRequest) error
		UpdateConfigForm(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormRequest, file *multipart.FileHeader) error
		GetShareInfos(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormFilter, limit, offset int) (int, *[]model.ShareInfoForm, error)
		GetShareInfoById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ShareInfoForm, error)
	}
	ShareInfo struct{}
)

func NewShareInfo() IShareInfo {
	return &ShareInfo{}
}

func (s *ShareInfo) PostConfigForm(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormRequest, file *multipart.FileHeader) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		err = errors.New(response.ERR_EMPTY_CONN)
		return err
	}

	filePath := file.Filename

	filter := model.ShareInfoFormFilter{
		TenantId:  authUser.TenantId,
		ShareType: data.ShareType,
		OaId:      data.OaId,
	}

	total, _, err := repository.ShareInfoRepo.GetShareInfos(ctx, dbCon, filter, -1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if total > 0 {
		log.Error("share config app_id " + data.AppId + " already exist")
		err = errors.New("share config app_id " + data.AppId + " already exist")
		return err
	}

	shareForm := model.ShareForm{}
	if data.ShareType == "facebook" {
	} else if data.ShareType == "zalo" {
		shareForm.Zalo.AppId = data.AppId
		shareForm.Zalo.ImageUrl = filePath
		shareForm.Zalo.Title = data.Title
		shareForm.Zalo.Subtitle = data.Subtitle
		shareForm.Zalo.OaId = data.OaId
	}

	shareInfoForm := model.ShareInfoForm{
		Base:      model.InitBase(),
		TenantId:  authUser.TenantId,
		ShareType: data.ShareType,
		ShareForm: shareForm,
	}

	if err := repository.ShareInfoRepo.Insert(ctx, dbCon, shareInfoForm); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

/**
* Send to ott service
 */
func (s *ShareInfo) PostRequestShareInfo(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormSubmitRequest) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		err = errors.New(response.ERR_EMPTY_CONN)
		return err
	}
	filter := model.ShareInfoFormFilter{
		TenantId:  authUser.TenantId,
		ShareType: data.ShareType,
		AppId:     data.AppId,
	}
	total, shareInfos, err := repository.ShareInfoRepo.GetShareInfos(ctx, dbCon, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if total < 1 {
		log.Error("share config app_id " + data.AppId + " not exist")
		err = errors.New("share config app_id " + data.AppId + " not exist")
		return err
	}
	var result model.OttResponse
	var body any
	tmp := model.OttShareInfoRequest{
		Type:      data.ShareType,
		EventName: data.EventName,
		AppId:     data.AppId,
		OaId:      (*shareInfos)[0].ShareForm.Zalo.OaId,
		Uid:       data.ExternalUserId,
		ImageUrl:  API_SHARE_INFO_HOST + "/" + (*shareInfos)[0].ShareForm.Zalo.ImageUrl,
		Title:     (*shareInfos)[0].ShareForm.Zalo.Title,
		Subtitle:  (*shareInfos)[0].ShareForm.Zalo.Subtitle,
	}

	log.Info("request share info: ", tmp)

	if err := util.ParseAnyToAny(tmp, &body); err != nil {
		log.Error(err)
		return err
	}

	url := OTT_URL + "/ott/" + OTT_VERSION + "/crm"
	client := resty.New()

	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		// SetHeader("Authorization", "Bearer "+token).
		SetBody(body).
		Post(url)
	if err != nil {
		log.Error(err)
		return err
	}

	if err := json.Unmarshal([]byte(res.Body()), &result); err != nil {
		log.Error(err)
		return err
	}
	if res.StatusCode() == 200 {
		return nil
	} else {
		err = errors.New(result.Message)
		return err
	}
}

func GetAvatarPageShareInfo(ctx context.Context, fileName string) (string, error) {
	input := storage.NewRetrieveInput(fileName)
	_, err := storage.Instance.Retrieve(ctx, *input)
	if err != nil {
		log.Error(err)
		return err.Error(), err
	}
	return "", nil
}

func (s *ShareInfo) UpdateConfigForm(ctx context.Context, authUser *model.AuthUser, data model.ShareInfoFormRequest, file *multipart.FileHeader) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		err = errors.New(response.ERR_EMPTY_CONN)
		return err
	}
	if len(data.Id) < 1 {
		log.Errorf("share config id %s not exist", data.Id)
		err = fmt.Errorf("share config id %s not exist", data.Id)
		return err
	}

	shareInfoExist, err := repository.ShareInfoRepo.GetById(ctx, dbCon, data.Id)
	if err != nil {
		log.Error(err)
		return err
	} else if shareInfoExist == nil {
		log.Errorf("share config id %s not exist", data.Id)
		err = fmt.Errorf("share config id %s not exist", data.Id)
		return err
	}

	filePath := file.Filename

	shareForm := model.ShareForm{}
	if data.ShareType == "facebook" {
	} else if data.ShareType == "zalo" {
		shareForm.Zalo.AppId = data.AppId
		shareForm.Zalo.ImageUrl = filePath
		shareForm.Zalo.Title = data.Title
		shareForm.Zalo.Subtitle = data.Subtitle
		shareForm.Zalo.OaId = data.OaId
	}

	if err := util.ParseAnyToAny(shareForm, &shareInfoExist.ShareForm); err != nil {
		log.Error(err)
		return err
	}

	if err := repository.ShareInfoRepo.Update(ctx, dbCon, *shareInfoExist); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *ShareInfo) GetShareInfos(ctx context.Context, authUser *model.AuthUser, filter model.ShareInfoFormFilter, limit, offset int) (int, *[]model.ShareInfoForm, error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		err = errors.New(response.ERR_EMPTY_CONN)
		return 0, nil, err
	}

	total, shareInfos, err := repository.ShareInfoRepo.GetShareInfos(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}
	return total, shareInfos, nil
}

func (s *ShareInfo) GetShareInfoById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ShareInfoForm, error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		err = errors.New(response.ERR_EMPTY_CONN)
		return nil, err
	}

	shareInfo, err := repository.ShareInfoRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return shareInfo, nil
}
