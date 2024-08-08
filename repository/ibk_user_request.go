package repository

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IIBKUserRequest interface {
		IRepo[model.IBKUserRequest]
	}
	IBKUserRequest struct {
		Repo[model.IBKUserRequest]
	}
)

var IBKUserRequestRepo IIBKUserRequest

func NewIBKUserRequest(conn sqlclient.ISqlClientConn) IIBKUserRequest {
	return &IBKUserRequest{
		Repo[model.IBKUserRequest]{
			Conn: conn,
		},
	}
}

func (repo *IBKUserRequest) InitTable(ctx context.Context) {
	if err := CreateTable(ctx, repo.Conn, (*model.IBKUserRequest)(nil)); err != nil {
		log.Error(err)
	}
}
