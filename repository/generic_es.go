package repository

import (
	"context"
	"time"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/tel4vn/fins-microservices/internal/elasticsearch"
	"github.com/tel4vn/fins-microservices/model"
)

type IESGenericRepo[T any] interface {
	ScrollAPI(ctx context.Context, scrollId string) (result *model.SearchReponse, err error)
}

type ESGenericRepo[T any] struct {
}

func NewESRepo[T any]() IESGenericRepo[T] {
	return &ESGenericRepo[T]{}
}

func (r *ESGenericRepo[T]) ScrollAPI(ctx context.Context, scrollId string) (result *model.SearchReponse, err error) {
	client := ESClient.GetClient()
	res, err := client.Scroll(
		client.Scroll.WithContext(ctx),
		client.Scroll.WithScrollID(scrollId),
		client.Scroll.WithScroll(time.Minute),
	)
	if err != nil {
		return
	}
	defer res.Body.Close()
	result, err = elasticsearch.ParseSearchResponse((*esapi.Response)(res))
	return
}
