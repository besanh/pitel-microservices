package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IAuthSource interface {
		InsertAuthSource(ctx context.Context, authUser *model.AuthUser, data model.AuthSource) (err error)
	}
	AuthSource struct{}
)

func NewAuthSource() IAuthSource {
	return &AuthSource{}
}

func (s *AuthSource) InsertAuthSource(ctx context.Context, authUser *model.AuthUser, data model.AuthSource) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	filter := model.AuthSourceFilter{
		Source: data.Source,
		Status: sql.NullBool{
			Valid: true,
			Bool:  true,
		},
	}

	total, _, err := repository.AuthSourceRepo.GetAuthSource(ctx, dbCon, filter, -1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if total > 0 {
		return errors.New("source already exists")
	}

	authSource := model.AuthSource{
		Base:     &model.Base{},
		TenantId: authUser.TenantId,
		Source:   data.Source,
		AuthUrl:  data.AuthUrl,
		Info:     data.Info,
		Status:   data.Status,
	}

	if err := repository.AuthSourceRepo.Insert(ctx, dbCon, authSource); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
