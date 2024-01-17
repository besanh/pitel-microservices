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
		InsertChatRouting(ctx context.Context, authUser *model.AuthUser, data *model.ChatRoutingRequest) error
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

func (s *ChatRouting) InsertChatRouting(ctx context.Context, authUser *model.AuthUser, data *model.ChatRoutingRequest) error {
	dbConn, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	if !slices.Contains[[]string](variables.CHAT_ROUTING, data.RoutingName) {
		err = errors.New("chat routing method is not supported")
		return err
	}

	total, _, err := repository.ChatRoutingRepo.GetChatRoutings(ctx, dbConn, model.ChatRoutingFilter{
		RoutingName: data.RoutingName,
		Status: sql.NullBool{
			Valid: true,
			Bool:  true,
		},
	}, 1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if total > 0 {
		err = errors.New("chat routing already exists")
		return err
	}

	chatRouting := model.ChatRouting{
		Base:        model.InitBase(),
		RoutingName: data.RoutingName,
		Status:      data.Status,
	}

	if err := repository.ChatRoutingRepo.Insert(ctx, dbConn, chatRouting); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (s *ChatRouting) GetChatRoutings(ctx context.Context, authUser *model.AuthUser, filter model.ChatRoutingFilter, limit, offset int) (int, *[]model.ChatRouting, error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

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
