package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatPolicySetting interface {
		GetChatPolicySettings(ctx context.Context, authUser *model.AuthUser, filter model.ChatPolicyFilter, limit int, offset int) (int, *[]model.ChatPolicySetting, error)
		GetChatPolicySettingById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ChatPolicySetting, error)
		InsertChatPolicySetting(ctx context.Context, authUser *model.AuthUser, request model.ChatPolicyConfigRequest) (string, error)
		UpdateChatPolicySettingById(ctx context.Context, authUser *model.AuthUser, id string, request model.ChatPolicyConfigRequest) error
		DeleteChatPolicySettingById(ctx context.Context, authUser *model.AuthUser, id string) error
	}

	ChatPolicySetting struct{}
)

func NewChatPolicySetting() IChatPolicySetting {
	return &ChatPolicySetting{}
}

func (s *ChatPolicySetting) GetChatPolicySettings(ctx context.Context, authUser *model.AuthUser, filter model.ChatPolicyFilter, limit int, offset int) (total int, policySettings *[]model.ChatPolicySetting, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	filter.TenantId = authUser.TenantId
	total, policySettings, err = repository.ChatPolicySettingRepo.GetChatPolicySettings(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func (s *ChatPolicySetting) GetChatPolicySettingById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.ChatPolicySetting, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	result, err = repository.ChatPolicySettingRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	}
	if result == nil {
		log.Error(errors.New("not found policy setting"))
		return
	}

	return
}

func (s *ChatPolicySetting) InsertChatPolicySetting(ctx context.Context, authUser *model.AuthUser, request model.ChatPolicyConfigRequest) (string, error) {
	policySetting := model.ChatPolicySetting{
		Base:     model.InitBase(),
		TenantId: authUser.TenantId,
	}
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return policySetting.Id, err
	}

	// check already exist setting
	filter := model.ChatPolicyFilter{
		TenantId:       authUser.TenantId,
		ConnectionType: request.ConnectionType,
	}
	total, _, err := repository.ChatPolicySettingRepo.GetChatPolicySettings(ctx, dbCon, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return policySetting.Id, err
	}
	if total > 0 {
		return policySetting.Id, errors.New("policy setting already exists")
	}

	policySetting.CreatedBy = authUser.UserId
	policySetting.ConnectionType = request.ConnectionType
	policySetting.ChatWindowTime = request.ChatWindowTime

	if err = repository.ChatPolicySettingRepo.Insert(ctx, dbCon, policySetting); err != nil {
		log.Error(err)
		return policySetting.Id, err
	}

	key := GeneratePolicySettingKeyId(policySetting.TenantId, policySetting.ConnectionType)
	if err = cache.RCache.Set(key, policySetting, CHAT_POLICY_SETTING_EXPIRE); err != nil {
		log.Error(err)
		return policySetting.Id, err
	}

	return policySetting.Id, nil
}

func (s *ChatPolicySetting) UpdateChatPolicySettingById(ctx context.Context, authUser *model.AuthUser, id string, request model.ChatPolicyConfigRequest) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	// check already exist setting
	filter := model.ChatPolicyFilter{
		TenantId:       authUser.TenantId,
		ConnectionType: request.ConnectionType,
	}
	total, _, err := repository.ChatPolicySettingRepo.GetChatPolicySettings(ctx, dbCon, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if total > 0 {
		return errors.New("policy setting already exists")
	}

	policySetting, err := repository.ChatPolicySettingRepo.GetById(ctx, dbCon, id)
	if err != nil || policySetting == nil {
		err = fmt.Errorf("not found policy setting, err=%v", err)
		log.Error(err)
		return err
	}
	policySetting.ConnectionType = request.ConnectionType
	policySetting.UpdatedBy = authUser.UserId
	if err = repository.ChatPolicySettingRepo.Update(ctx, dbCon, *policySetting); err != nil {
		log.Error(err)
		return err
	}

	// clear cache
	key := GeneratePolicySettingKeyId(policySetting.TenantId, policySetting.ConnectionType)
	if err = cache.RCache.Del([]string{key}); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *ChatPolicySetting) DeleteChatPolicySettingById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	policySetting, err := repository.ChatPolicySettingRepo.GetById(ctx, dbCon, id)
	if err != nil || policySetting == nil {
		err = fmt.Errorf("not found policy setting, err=%v", err)
		log.Error(err)
		return err
	}

	if err = repository.ChatPolicySettingRepo.Delete(ctx, dbCon, id); err != nil {
		log.Error(err)
		return
	}

	// clear cache
	key := GeneratePolicySettingKeyId(policySetting.TenantId, policySetting.ConnectionType)
	if err = cache.RCache.Del([]string{key}); err != nil {
		log.Error(err)
		return err
	}
	return
}
