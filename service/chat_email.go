package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatEmail interface {
		GetChatEmails(ctx context.Context, authUser *model.AuthUser, filter model.ChatEmailFilter, limit, offset int) (int, *[]model.ChatEmailCustom, error)
		InsertChatEmail(ctx context.Context, authUser *model.AuthUser, request model.ChatEmailRequest) (string, error)
		GetChatEmailById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ChatEmail, error)
		UpdateChatEmailById(ctx context.Context, authUser *model.AuthUser, id string, request model.ChatEmailRequest) error
		DeleteChatEmailById(ctx context.Context, authUser *model.AuthUser, id string) error
		HandleJobExpireToken()
	}
	ChatEmail struct{}
)

func NewChatEmail() IChatEmail {
	repo := &ChatEmail{}
	return repo
}

func (s *ChatEmail) GetChatEmails(ctx context.Context, authUser *model.AuthUser, filter model.ChatEmailFilter, limit, offset int) (int, *[]model.ChatEmailCustom, error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}
	filter.TenantId = authUser.TenantId
	total, emails, err := repository.NewChatEmail().GetChatEmailsCustom(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	return total, emails, nil
}

func (s *ChatEmail) InsertChatEmail(ctx context.Context, authUser *model.AuthUser, request model.ChatEmailRequest) (string, error) {
	chatEmail := model.ChatEmail{
		Base: model.InitBase(),
	}
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return chatEmail.Id, err
	}

	// Check oa exist
	filterConnection := model.ChatConnectionAppFilter{
		TenantId: authUser.TenantId,
		OaId:     request.OaId,
	}
	_, connections, err := repository.NewConnectionApp().GetChatConnectionApp(ctx, dbCon, filterConnection, 1, 0)
	if err != nil {
		log.Error(err)
		return chatEmail.Id, err
	}
	if len(*connections) < 1 {
		log.Error("oa " + request.OaId + " not found")
		return chatEmail.Id, errors.New("oa " + request.OaId + " not found")
	}

	filter := model.ChatEmailFilter{
		TenantId: authUser.TenantId,
		OaId:     request.OaId,
	}
	_, chatEmails, err := repository.NewChatEmail().GetChatEmails(ctx, dbCon, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return chatEmail.Id, err
	}
	if len(*chatEmails) > 0 {
		log.Error("chat email oa_id " + request.OaId + " already exist")
		return chatEmail.Id, errors.New("chat email oa_id " + request.OaId + " already exist")
	}

	chatEmail.TenantId = authUser.TenantId
	chatEmail.OaId = request.OaId
	chatEmail.EmailSubject = request.EmailSubject
	chatEmail.EmailRecipient = request.EmailRecipient
	chatEmail.EmailContent = request.EmailContent

	if request.EmailRequestType == "manual" {
		chatEmail.EmailServer = request.EmailServer
		chatEmail.EmailUsername = request.EmailUsername
		chatEmail.EmailPassword = request.EmailPassword
		chatEmail.EmailPort = request.EmailPort
		chatEmail.EmailEncryptType = request.EmailEncryptType
	} else {
		chatEmail.EmailServer = SMTP_SERVER
		chatEmail.EmailUsername = SMTP_USERNAME
		chatEmail.EmailPassword = SMTP_PASSWORD
		chatEmail.EmailPort = fmt.Sprintf("%d", SMTP_MAILPORT)
		chatEmail.EmailEncryptType = "tls"
	}

	chatEmail.EmailStatus = request.EmailStatus
	chatEmail.CreatedAt = time.Now()

	if err := repository.NewChatEmail().Insert(ctx, dbCon, chatEmail); err != nil {
		log.Error(err)
		return chatEmail.Id, err
	}
	return chatEmail.Id, nil
}

func (s *ChatEmail) GetChatEmailById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ChatEmail, error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	chatEmail, err := repository.NewChatEmail().GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return chatEmail, nil
}

func (s *ChatEmail) UpdateChatEmailById(ctx context.Context, authUser *model.AuthUser, id string, request model.ChatEmailRequest) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	// Check oa exist
	filterConnection := model.ChatConnectionAppFilter{
		TenantId: authUser.TenantId,
		OaId:     request.OaId,
	}
	_, connections, err := repository.NewConnectionApp().GetChatConnectionApp(ctx, dbCon, filterConnection, 1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if len(*connections) < 1 {
		log.Error("oa " + request.OaId + " not found")
		return errors.New("oa " + request.OaId + " not found")
	}

	chatEmailExist, err := repository.NewChatEmail().GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	} else if chatEmailExist == nil {
		log.Error("chat email " + id + " not found")
		return errors.New("chat email " + id + " not found")
	}

	if request.EmailRequestType == "manual" {
		chatEmailExist.EmailServer = request.EmailServer
		chatEmailExist.EmailUsername = request.EmailUsername
		chatEmailExist.EmailPassword = request.EmailPassword
		chatEmailExist.EmailPort = request.EmailPort
		chatEmailExist.EmailEncryptType = request.EmailEncryptType
	} else {
		chatEmailExist.EmailServer = SMTP_SERVER
		chatEmailExist.EmailUsername = SMTP_USERNAME
		chatEmailExist.EmailPassword = SMTP_PASSWORD
		chatEmailExist.EmailPort = fmt.Sprintf("%d", SMTP_MAILPORT)
		chatEmailExist.EmailEncryptType = "tls"
	}
	chatEmailExist.EmailRecipient = request.EmailRecipient
	chatEmailExist.EmailServer = request.EmailServer
	chatEmailExist.EmailSubject = request.EmailSubject
	chatEmailExist.EmailUsername = request.EmailUsername
	chatEmailExist.UpdatedAt = time.Now()

	if err = repository.NewChatEmail().Update(ctx, dbCon, *chatEmailExist); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (s *ChatEmail) DeleteChatEmailById(ctx context.Context, authUser *model.AuthUser, id string) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}
	chatEmailExist, err := repository.NewChatEmail().GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	} else if chatEmailExist == nil {
		log.Error("chat email " + id + " not found")
		return errors.New("chat email " + id + " not found")
	}
	if err = repository.NewChatEmail().Delete(ctx, dbCon, id); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
