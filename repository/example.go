package repository

import (
	"github.com/tel4vn/pitel-microservices/model"
)

type (
	IExample interface {
		IRepo[model.Example]
	}
	Example struct {
		Repo[model.Example]
	}
)

var ExampleRepo IExample

func NewExample() IExample {
	repo := &Example{}
	return repo
}
