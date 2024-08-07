package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IBssUser interface {
		IRepo[model.BSS_User]
		GetUserByUsername(ctx context.Context, userName string) (result *model.BSS_User, err error)
	}
	BssUser struct {
		Repo[model.BSS_User]
	}
)

var BssUserRepo IBssUser

func NewBssUser(conn sqlclient.ISqlClientConn) IBssUser {
	return &BssUser{
		Repo[model.BSS_User]{
			Conn: conn,
		},
	}
}

func (repo *BssUser) GetUserByUsername(ctx context.Context, userName string) (result *model.BSS_User, err error) {
	result = new(model.BSS_User)
	query := repo.Conn.GetDB().NewSelect().
		Model(result).
		Where("user_name = ?", userName)
	err = query.Scan(ctx)
	if err == sql.ErrNoRows {
		return
	} else if err != nil {
		result = nil
	}
	return
}
