package service

import (
	"context"
	"errors"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IFacebook interface {
		InsertFacebookPage(ctx context.Context, authUser *model.AuthUser, data model.FacebookPageInfo) (string, error)
		BulkInsertFacebookPage(ctx context.Context, authUser *model.AuthUser, data []model.FacebookPageInfo) error
		GetFacebookPages(ctx context.Context, authUser *model.AuthUser, filter model.FacebookPageFilter, limit, offset int) (int, *[]model.FacebookPage, error)
	}
	Facebook struct{}
)

func NewFacebook() IFacebook {
	return &Facebook{}
}

func (s *Facebook) InsertFacebookPage(ctx context.Context, authUser *model.AuthUser, data model.FacebookPageInfo) (string, error) {
	facebook := model.FacebookPage{
		Base: model.InitBase(),
	}
	dbConn, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return facebook.GetId(), err
	}

	filter := model.AppFilter{
		AppType:    "facebook",
		DefaultApp: "active",
	}

	total, appInfos, err := repository.ChatAppRepo.GetChatApp(ctx, dbConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return facebook.GetId(), err
	}
	if total < 1 {
		return facebook.GetId(), errors.New("facebook app not found")
	}

	facebook.TenantId = authUser.TenantId
	facebook.AppId = (*appInfos)[0].Id
	facebook.OaId = data.OaId
	facebook.OaName = data.OaName
	facebook.TokenType = data.TokenType
	facebook.AccessToken = data.AccessToken
	facebook.Avatar = data.Avatar
	facebook.Status = data.Status

	if err := repository.NewFacebook().Insert(ctx, dbConn, facebook); err != nil {
		log.Error(err)
		return facebook.GetId(), err
	}
	return facebook.GetId(), nil
}

func (s *Facebook) GetFacebookPages(ctx context.Context, authUser *model.AuthUser, filter model.FacebookPageFilter, limit, offset int) (int, *[]model.FacebookPage, error) {
	dbConn, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}
	return repository.NewFacebook().GetFacebookPages(ctx, dbConn, filter, limit, offset)
}

func (s *Facebook) BulkInsertFacebookPage(ctx context.Context, authUser *model.AuthUser, data []model.FacebookPageInfo) error {
	dbConn, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	filter := model.AppFilter{
		AppType:    "facebook",
		DefaultApp: "active",
	}

	total, appInfos, err := repository.ChatAppRepo.GetChatApp(ctx, dbConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if total < 1 {
		return errors.New("facebook app not found")
	}

	facebooks := []model.FacebookPage{}
	if len(data) > 0 {
		for _, item := range data {
			facebook := model.FacebookPage{
				Base:        model.InitBase(),
				TenantId:    authUser.TenantId,
				AppId:       (*appInfos)[0].Id,
				OaId:        item.OaId,
				OaName:      item.OaName,
				TokenType:   item.TokenType,
				AccessToken: item.AccessToken,
				Avatar:      item.Avatar,
				Status:      item.Status,
			}

			facebooks = append(facebooks, facebook)
		}
	}

	return repository.NewFacebook().BulkInsert(ctx, dbConn, facebooks)
}
