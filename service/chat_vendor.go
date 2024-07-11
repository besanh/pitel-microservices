package service

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatVendor interface {
		GetChatVendors(ctx context.Context, authUser *model.AuthUser, filter model.ChatVendorFilter, limit, offset int) (total int, result *[]model.ChatVendor, err error)
		InsertChatVendor(ctx context.Context, authUser *model.AuthUser, data model.ChatVendorRequest) (id string, err error)

		// TODO: write api insert use form
	}
	ChatVendor struct{}
)

var ChatVendorService IChatVendor

func NewChatVendor() IChatVendor {
	return &ChatVendor{}
}

func (s *ChatVendor) GetChatVendors(ctx context.Context, authUser *model.AuthUser, filter model.ChatVendorFilter, limit, offset int) (total int, result *[]model.ChatVendor, err error) {
	total, result, err = repository.VendorRepo.GetVendors(ctx, repository.DBConn, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *ChatVendor) InsertChatVendor(ctx context.Context, authUser *model.AuthUser, data model.ChatVendorRequest) (id string, err error) {
	vendor := model.ChatVendor{
		Base:       model.InitBase(),
		VendorName: data.VendorName,
		VendorType: data.VendorType,
		Logo:       data.Logo,
		Status:     data.Status,
	}

	if err = repository.VendorRepo.Insert(ctx, repository.DBConn, vendor); err != nil {
		log.Error(err)
		return
	}

	id = vendor.Base.GetId()
	return
}
