package service

import (
	"context"
	"errors"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	ITemplateBss interface {
		InsertTemplateBss(ctx context.Context, authUser *model.AuthUser, data model.TemplateBssBodyRequest) error
		GetTemplateBssById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.TemplateBss, err error)
		GetTemplateBsses(ctx context.Context, authUser *model.AuthUser, filter model.TemplateBssFilter, limit, offset int) (total int, result *[]model.TemplateBssView, err error)
		DeleteTemplateBssById(ctx context.Context, authUser *model.AuthUser, id string) (err error)
		PutTemplateBssById(ctx context.Context, authUser *model.AuthUser, id string, data model.TemplateBssBodyRequest) (err error)
	}
	TemplateBss struct{}
)

func NewTemplateBss() ITemplateBss {
	return &TemplateBss{}
}

func (s *TemplateBss) InsertTemplateBss(ctx context.Context, authUser *model.AuthUser, data model.TemplateBssBodyRequest) error {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err != nil {
		return err
	}
	filter := model.TemplateBssFilter{
		TemplateCode: []string{data.TemplateCode},
		TemplateType: []string{data.TemplateType},
	}
	total, _, err := repository.TemplateBssRepo.GetTemplateBsses(ctx, dbCon, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if total > 0 {
		return errors.New("template code is existed")
	}

	partitionContents, checkWrongFormat := util.GetPartitionContentTemplate(data.Content)
	if checkWrongFormat {
		return errors.New("content wrong format")
	}
	joinPart := util.GetJoinPartTemplate(partitionContents)
	templateBss := model.TemplateBss{
		Base:         model.InitBase(),
		TemplateName: data.TemplateName,
		TemplateCode: data.TemplateCode,
		TemplateType: data.TemplateType,
		Content:      data.Content,
		Partition:    joinPart,
		Status:       data.Status,
	}

	if err := repository.TemplateBssRepo.Insert(ctx, dbCon, templateBss); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *TemplateBss) GetTemplateBsses(ctx context.Context, authUser *model.AuthUser, filter model.TemplateBssFilter, limit, offset int) (total int, result *[]model.TemplateBssView, err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err != nil {
		return 0, nil, err
	}

	total, result, err = repository.TemplateBssRepo.GetTemplateBsses(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	return total, result, nil
}

func (s *TemplateBss) GetTemplateBssById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.TemplateBss, err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err != nil {
		return nil, err
	}

	templateBssExist, err := repository.TemplateBssRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return templateBssExist, nil
}

func (s *TemplateBss) PutTemplateBssById(ctx context.Context, authUser *model.AuthUser, id string, data model.TemplateBssBodyRequest) (err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err != nil {
		return err
	}

	filter := model.TemplateBssFilter{
		TemplateCode: []string{data.TemplateCode},
		TemplateType: []string{data.TemplateType},
	}
	total, _, err := repository.TemplateBssRepo.GetTemplateBsses(ctx, dbCon, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if total > 1 {
		return errors.New("template code is existed")
	}

	templateBssExist, err := repository.TemplateBssRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	partitionContents, checkWrongFormat := util.GetPartitionContentTemplate(data.Content)
	if checkWrongFormat {
		return errors.New("content wrong format")
	}
	joinPart := util.GetJoinPartTemplate(partitionContents)

	templateBssExist.TemplateName = data.TemplateName
	templateBssExist.TemplateCode = data.TemplateCode
	templateBssExist.TemplateType = data.TemplateType
	templateBssExist.Content = data.Content
	templateBssExist.Partition = joinPart
	templateBssExist.Status = data.Status

	if err = repository.TemplateBssRepo.Update(ctx, dbCon, *templateBssExist); err != nil {
		log.Error(err)
		return err
	}

	return
}

func (s *TemplateBss) DeleteTemplateBssById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err != nil {
		return err
	}

	_, err = repository.TemplateBssRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	if err = repository.TemplateBssRepo.Delete(ctx, dbCon, id); err != nil {
		log.Error(err)
		return err
	}
	return
}
