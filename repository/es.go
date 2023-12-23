package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/tel4vn/fins-microservices/common/log"
	rabbitmq "github.com/tel4vn/fins-microservices/internal/rabbitmq/driver"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IElasticsearch interface {
		CheckAliasExist(ctx context.Context, index, alias string) (bool, error)
		CreateAlias(ctx context.Context, index, alias string) error
		CreateDocRabbitMQ(ctx context.Context, index, tenant, routing, uuid string, esDoc map[string]any) (bool, error)
		CreateAliasRabbitMQ(ctx context.Context, index, alias string) (bool, error)
		InsertLog(ctx context.Context, tenantId, index, docId string, esDoc map[string]any) error
		UpdateDocById(ctx context.Context, index, docId string, esDoc map[string]any) error
	}
	Elasticsearch struct{}
)

var ESRepo IElasticsearch

func NewES() IElasticsearch {
	return &Elasticsearch{}
}

func (repo *Elasticsearch) CheckAliasExist(ctx context.Context, index, alias string) (bool, error) {
	res, err := ESClient.GetClient().Aliases().
		Index(index).
		Do(ctx)
	if err != nil {
		return false, err
	}
	if len(res.Indices) > 0 {
		indices := res.Indices[index]
		if indices.HasAlias(index + "_" + alias) {
			return true, nil
		} else {
			return false, nil
		}
	}
	return false, nil
}

func (repo *Elasticsearch) CreateAlias(ctx context.Context, index, alias string) error {
	_, err := ESClient.GetClient().Alias().Action().
		Add(index, index+"_"+alias).
		Do(ctx)
	return err
}

func (repo *Elasticsearch) CreateAliasRabbitMQ(ctx context.Context, index, alias string) (bool, error) {
	log.Infof("create alias: %s", alias)
	data := model.AliasCreate{
		Index: index,
		Name:  index + "_" + alias,
	}
	var actions []any
	addAction := make(map[string]any)
	addAction["add"] = data
	actions = append(actions, addAction)
	bodyData := make(map[string]any)
	bodyData["actions"] = actions
	var payload model.RabbitMQPayload
	payload.HttpMethod = "POST"
	payload.Uri = "/_aliases"
	payload.Body = bodyData
	err := rabbitmq.RabbitConnector.Publish(payload)
	if err != nil {
		log.Error(err)
		return false, err
	} else {
		log.Infof("alias %s is created", alias)
		return true, nil
	}
}
func (repo *Elasticsearch) CreateDocRabbitMQ(ctx context.Context, index, tenant, routing, uuid string, esDoc map[string]any) (bool, error) {
	log.Infof("push log inbox marketing %s to rabbitmq", uuid)
	payload := model.RabbitMQPayload{
		HttpMethod: "POST",
		Uri:        "/" + index + "_" + tenant + "/_doc/" + uuid + "/_create?routing=" + index + "_" + routing,
		Body:       esDoc,
	}
	err := rabbitmq.RabbitConnector.Publish(payload)
	if err != nil {
		log.Error(err)
		return false, err
	} else {
		return true, nil
	}
}

func (repo *Elasticsearch) InsertLog(ctx context.Context, tenantId, index, docId string, esDoc map[string]any) error {
	bdy, err := json.Marshal(esDoc)
	if err != nil {
		return err
	}
	req := esapi.CreateRequest{
		Index:      index,
		DocumentID: docId,
		Routing:    index + "_" + tenantId,
		Body:       bytes.NewReader(bdy),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, ES.GetClient())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("insert: response: %s", res.String())
	}

	return nil
}

func (repo *Elasticsearch) UpdateDocById(ctx context.Context, index, docId string, esDoc map[string]any) error {
	bdy, err := json.Marshal(esDoc)
	if err != nil {
		return err
	}
	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: docId,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, bdy))),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, ES.GetClient())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("update: response: %s", res.String())
	}

	return nil
}
