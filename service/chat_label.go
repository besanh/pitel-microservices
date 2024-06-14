package service

import (
	"context"
	"errors"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatLabel interface {
		InsertChatLabel(ctx context.Context, authUser *model.AuthUser, request *model.ChatLabelRequest) (id string, err error)
		GetChatLabels(ctx context.Context, authUser *model.AuthUser, filter model.ChatLabelFilter, limit, offset int) (int, *[]model.ChatLabel, error)
		GetChatLabelById(ctx context.Context, authUser *model.AuthUser, id string) (chatLabel *model.ChatLabel, err error)
		UpdateChatLabelById(ctx context.Context, authUser *model.AuthUser, id string, request *model.ChatLabelRequest) (err error)
		DeleteChatLabelById(ctx context.Context, authUser *model.AuthUser, id string) (err error)
	}
	ChatLabel struct{}
)

func NewChatLabel() IChatLabel {
	return &ChatLabel{}
}

func (s *ChatLabel) InsertChatLabel(ctx context.Context, authUser *model.AuthUser, request *model.ChatLabelRequest) (string, error) {
	chatLabel := model.ChatLabel{
		Base:     model.InitBase(),
		TenantId: authUser.TenantId,
	}
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return chatLabel.GetId(), err
	}

	filter := model.ChatConnectionAppFilter{
		AppId: request.AppId,
		OaId:  request.OaId,
	}
	_, connectionExist, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return chatLabel.GetId(), err
	} else if len(*connectionExist) < 1 {
		log.Error("connection with app_id: " + request.AppId + ", oa_id: " + request.OaId + " not found")
		return chatLabel.GetId(), errors.New("connection with app_id: " + request.AppId + ", oa_id: " + request.OaId + " not found")
	}

	chatLabel.AppId = request.AppId
	chatLabel.OaId = request.OaId
	chatLabel.LabelName = request.LabelName
	chatLabel.LabelType = (*connectionExist)[0].ConnectionType
	chatLabel.LabelColor = request.LabelColor
	chatLabel.LabelStatus = true
	chatLabel.CreatedBy = authUser.UserId

	if err := repository.ChatLabelRepo.Insert(ctx, dbCon, chatLabel); err != nil {
		log.Error(err)
		return chatLabel.GetId(), err
	}

	return chatLabel.GetId(), nil
}

func (s *ChatLabel) GetChatLabels(ctx context.Context, authUser *model.AuthUser, filter model.ChatLabelFilter, limit, offset int) (total int, result *[]model.ChatLabel, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	total, result, err = repository.ChatLabelRepo.GetChatLabels(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}

	return total, result, nil
}

func (s *ChatLabel) GetChatLabelById(ctx context.Context, authUser *model.AuthUser, id string) (chatLabel *model.ChatLabel, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	chatLabel, err = repository.ChatLabelRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return chatLabel, nil
}

func (s *ChatLabel) UpdateChatLabelById(ctx context.Context, authUser *model.AuthUser, id string, request *model.ChatLabelRequest) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	chatLabelExist, err := repository.ChatLabelRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	} else if chatLabelExist == nil {
		log.Error("chat label not found")
		return errors.New("chat label not found")
	}

	requestExternal := model.ChatExtenalLabelRequest{
		AppId:   chatLabelExist.AppId,
		OaId:    chatLabelExist.OaId,
		TagName: chatLabelExist.LabelName,
	}

	// TODO: update zalo and facebook
	// because zalo and facebook does not support update label, so we can only remove that label
	// with zalo, we only can create label with external_id~uid
	// with facebook, we can remove and create new label
	if chatLabelExist.LabelType == "zalo" {
		if err = s.RequestZaloLabel(ctx, "remove-label", requestExternal); err != nil {
			log.Error(err)
			return
		}
	} else if chatLabelExist.LabelType == "facebook" {
		if err = s.RequestFacebookLabel(ctx, "", requestExternal); err != nil {
			log.Error(err)
			return
		}
		if err = s.RequestFacebookLabel(ctx, "create-label", requestExternal); err != nil {
			log.Error(err)
			return
		}
	}

	filter := model.ChatConnectionAppFilter{
		AppId: request.AppId,
		OaId:  request.OaId,
	}
	_, connectionExist, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return err
	} else if len(*connectionExist) < 1 {
		log.Error("connection with app_id: " + request.AppId + ", oa_id: " + request.OaId + " not found")
		return errors.New("connection with app_id: " + request.AppId + ", oa_id: " + request.OaId + " not found")
	}

	chatLabelExist.LabelType = (*connectionExist)[0].ConnectionType
	chatLabelExist.AppId = request.AppId
	chatLabelExist.OaId = request.OaId
	chatLabelExist.LabelName = request.LabelName
	chatLabelExist.LabelColor = request.LabelColor
	chatLabelExist.UpdatedBy = authUser.UserId

	if err := repository.ChatLabelRepo.Update(ctx, dbCon, *chatLabelExist); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *ChatLabel) DeleteChatLabelById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}
	chatLabelExist, err := repository.ChatLabelRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	request := model.ChatExtenalLabelRequest{
		AppId:   chatLabelExist.AppId,
		OaId:    chatLabelExist.OaId,
		TagName: chatLabelExist.LabelName,
		LabelId: chatLabelExist.ExternalLabelId,
	}

	// TODO: zalo
	if chatLabelExist.LabelType == "zalo" {
		if err = s.RequestZaloLabel(ctx, "zalo/remove-label", request); err != nil {
			log.Error(err)
			return
		}
	} else if chatLabelExist.LabelType == "facebook" {
		if err = s.RequestFacebookLabel(ctx, "face/remove-label", request); err != nil {
			log.Error(err)
			return
		}
	}

	err = repository.ChatLabelRepo.Delete(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
