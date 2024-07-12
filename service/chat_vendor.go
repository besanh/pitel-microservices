package service

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatVendor interface {
		GetChatVendors(ctx context.Context, authUser *model.AuthUser, filter model.ChatVendorFilter, limit, offset int) (total int, result *[]model.ChatVendor, err error)
		InsertChatVendor(ctx context.Context, authUser *model.AuthUser, data model.ChatVendorRequest) (id string, err error)

		// TODO: write api insert use form
		PostChatVendorUpload(ctx context.Context, authUser *model.AuthUser, data model.ChatVendorRequest, file *multipart.FileHeader) (id string, err error)
		PutChatVendorUpload(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatVendorRequest, file *multipart.FileHeader) (err error)
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

func (s *ChatVendor) PostChatVendorUpload(ctx context.Context, authUser *model.AuthUser, data model.ChatVendorRequest, file *multipart.FileHeader) (id string, err error) {
	vendor := model.ChatVendor{
		Base:       model.InitBase(),
		VendorName: data.VendorName,
		VendorType: data.VendorType,
		Status:     data.Status,
	}

	fileUrl, err := UploadDoc(ctx, "", "", file)
	if err != nil {
		log.Error(err)
		return
	} else if fileUrl == "" {
		log.Error("file url is empty")
		return
	}

	vendor.Logo = fileUrl
	if err = repository.VendorRepo.Insert(ctx, repository.DBConn, vendor); err != nil {
		log.Error(err)
		return
	}

	id = vendor.Base.GetId()
	return
}

func (s *ChatVendor) PutChatVendorUpload(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatVendorRequest, file *multipart.FileHeader) (err error) {
	vendorExist, err := repository.VendorRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return err
	} else if vendorExist == nil {
		log.Error("vendor does not exist")
		return errors.New("vendor does not exist")
	}

	fileUrl, err := UploadDoc(ctx, "", "", file)
	if err != nil {
		log.Error(err)
		return
	} else if fileUrl == "" {
		log.Error("file url is empty")
		return
	}

	vendorExist.VendorName = data.VendorName
	vendorExist.VendorType = data.VendorType
	vendorExist.Status = data.Status
	vendorExist.Logo = fileUrl
	if err = repository.VendorRepo.Update(ctx, repository.DBConn, *vendorExist); err != nil {
		log.Error(err)
		return
	}

	return
}
