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
	}
	ChatRouting struct{}
)

func NewChatRouting() IChatRouting {
	return &ChatRouting{}
}

func (s *ChatRouting) InsertChatRouting(ctx context.Context, authUser *model.AuthUser, data *model.ChatRoutingRequest) error {
	dbConn, err := GetDBConnOfUser(*authUser)
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
