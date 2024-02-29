package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"golang.org/x/exp/slices"
)

type (
	IChatRouting interface {
		InsertChatRouting(ctx context.Context, authUser *model.AuthUser, data *model.ChatRoutingRequest) (string, error)
		GetChatRoutings(ctx context.Context, authUser *model.AuthUser, filter model.ChatRoutingFilter, limit, offset int) (int, *[]model.ChatRouting, error)
		GetChatRoutingById(ctx context.Context, authUser *model.AuthUser, id string) (model.ChatRouting, error)
		UpdateChatRoutingById(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatRoutingRequest) error
		DeleteChatRoutingById(ctx context.Context, authUser *model.AuthUser, id string) error
	}
	ChatRouting struct{}
)

func NewChatRouting() IChatRouting {
	return &ChatRouting{}
}

func (s *ChatRouting) InsertChatRouting(ctx context.Context, authUser *model.AuthUser, data *model.ChatRoutingRequest) (string, error) {
	chatRouting := model.ChatRouting{
		Base: model.InitBase(),
	}
	dbConn, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return chatRouting.GetId(), err
	}

	if !slices.Contains[[]string](variables.CHAT_ROUTING, data.RoutingAlias) {
		err = errors.New("chat routing method is not supported")
		return chatRouting.GetId(), err
	}

	total, _, err := repository.ChatRoutingRepo.GetChatRoutings(ctx, dbConn, model.ChatRoutingFilter{
		TenantId:     authUser.TenantId,
		RoutingAlias: data.RoutingAlias,
		Status: sql.NullBool{
			Valid: true,
			Bool:  true,
		},
	}, 1, 0)
	if err != nil {
		log.Error(err)
		return chatRouting.GetId(), err
	}
	if total > 0 {
		log.Error("chat routing already " + data.RoutingAlias + " exists")
		err = errors.New("chat routing " + data.RoutingAlias + " already exists")
		return chatRouting.GetId(), err
	}

	chatRouting.RoutingName = data.RoutingName
	chatRouting.RoutingAlias = data.RoutingAlias
	chatRouting.Status = data.Status
	chatRouting.TenantId = authUser.TenantId

	if err := repository.ChatRoutingRepo.Insert(ctx, dbConn, chatRouting); err != nil {
		log.Error(err)
		return chatRouting.GetId(), err
	}
	return chatRouting.GetId(), nil
}

func (s *ChatRouting) GetChatRoutings(ctx context.Context, authUser *model.AuthUser, filter model.ChatRoutingFilter, limit, offset int) (int, *[]model.ChatRouting, error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}
	filter.TenantId = authUser.TenantId

	total, chatRoutings, err := repository.ChatRoutingRepo.GetChatRoutings(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	return total, chatRoutings, nil
}

func (s *ChatRouting) GetChatRoutingById(ctx context.Context, authUser *model.AuthUser, id string) (model.ChatRouting, error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return model.ChatRouting{}, err
	}

	chatRouting, err := repository.ChatRoutingRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return model.ChatRouting{}, err
	}

	return *chatRouting, nil
}

func (s *ChatRouting) UpdateChatRoutingById(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatRoutingRequest) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	chatRouting, err := repository.ChatRoutingRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	chatRouting.RoutingName = data.RoutingName
	chatRouting.RoutingAlias = data.RoutingAlias
	chatRouting.Status = data.Status
	err = repository.ChatRoutingRepo.Update(ctx, dbCon, *chatRouting)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (s *ChatRouting) DeleteChatRoutingById(ctx context.Context, authUser *model.AuthUser, id string) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = repository.ChatRoutingRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}
	err = repository.ChatRoutingRepo.Delete(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
