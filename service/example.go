package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IExample interface {
		GetExamples(ctx context.Context, authUser *model.AuthUser, queryParams []model.Param, limit int, offset int) (examples *[]model.Example, total int, err error)
		GetExampleById(ctx context.Context, authUser *model.AuthUser, id string) (example model.Example, err error)
		PostExample(ctx context.Context, authUser *model.AuthUser, example model.Example) (err error)
		PutExample(ctx context.Context, authUser *model.AuthUser, id string, example model.Example) (err error)
		DeleteExampleById(ctx context.Context, authUser *model.AuthUser, id string) (err error)
	}
	Example struct{}
)

var ExampleService IExample

func NewExample() IExample {
	return &Example{}
}

func (c *Example) GetExamples(ctx context.Context, authUser *model.AuthUser, queryParams []model.Param, limit int, offset int) (examples *[]model.Example, total int, err error) {
	dbConn, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return nil, 0, errors.New(response.ERR_EMPTY_CONN)
	}
	queryParams = append(queryParams, model.Param{
		// Key:      "tenant_id",
		// Value:    authUser.TenantId,
		// Operator: "=",
	})
	examples, total, err = repository.ExampleRepo.SelectByQuery(ctx, dbConn, queryParams, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (c *Example) GetExampleById(ctx context.Context, authUser *model.AuthUser, id string) (result model.Example, err error) {
	dbConn, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		err = errors.New(response.ERR_EMPTY_CONN)
		return
	}
	example, err := repository.ExampleRepo.GetById(ctx, dbConn, id)
	if err != nil {
		log.Error(err)
		return
	} else if example == nil {
		err = errors.New(response.ERR_EXAMPLE_NOT_FOUND)
		return
	}
	result = *example
	return
}

func (c *Example) PostExample(ctx context.Context, authUser *model.AuthUser, example model.Example) (err error) {
	dbConn, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		err = errors.New(response.ERR_EMPTY_CONN)
		return
	}
	example.Base = &model.Base{
		Id: uuid.NewString(),
	}
	err = repository.ExampleRepo.Insert(ctx, dbConn, example)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (c *Example) PutExample(ctx context.Context, authUser *model.AuthUser, id string, example model.Example) (err error) {
	dbConn, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		err = errors.New(response.ERR_EMPTY_CONN)
		return
	}
	exampleExist, err := repository.ExampleRepo.GetById(ctx, dbConn, id)
	if err != nil {
		log.Error(err)
		return
	} else if exampleExist == nil {
		err = errors.New(response.ERR_EXAMPLE_NOT_FOUND)
		return
	}
	exampleExist.ExampleName = example.ExampleName
	err = repository.ExampleRepo.Update(ctx, dbConn, *exampleExist)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (c *Example) DeleteExampleById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbConn, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		err = errors.New(response.ERR_EMPTY_CONN)
		return
	}
	exampleExist, err := repository.ExampleRepo.GetById(ctx, dbConn, id)
	if err != nil {
		log.Error(err)
		return
	} else if exampleExist == nil {
		err = errors.New(response.ERR_EXAMPLE_NOT_FOUND)
		return
	}
	err = repository.ExampleRepo.Delete(ctx, dbConn, id)
	if err != nil {
		log.Error(err)
		return
	}
	return
}
