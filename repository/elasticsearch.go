package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/esapi"
)

type (
	IElasticsearch interface {
		CheckAliasExist(ctx context.Context, index, alias string) (bool, error)
		CreateAlias(ctx context.Context, index, alias string) error
		InsertLog(ctx context.Context, tenantId, index, appId, docId string, esDoc map[string]any) error
		UpdateDocById(ctx context.Context, index, appId, docId string, esDoc map[string]any) error
		BulkUpdateDoc(ctx context.Context, index string, esDoc map[string]any) error
		DeleteById(ctx context.Context, index, docId string) error
	}
	Elasticsearch struct{}
)

var ESRepo IElasticsearch

func NewES() IElasticsearch {
	return &Elasticsearch{}
}

func (repo *Elasticsearch) CheckAliasExist(ctx context.Context, index, alias string) (bool, error) {
	idx := index
	if len(alias) > 0 {
		idx += "_" + alias
	}
	res, err := esapi.IndicesExistsAliasRequest.Do(esapi.IndicesExistsAliasRequest{
		Index: []string{idx},
	}, ctx,
		ESClient.GetClient().Transport)
	if err != nil {
		return false, err
	} else if res.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}

func (repo *Elasticsearch) CreateAlias(ctx context.Context, index, alias string) error {
	_, err := esapi.IndicesPutAliasRequest.Do(esapi.IndicesPutAliasRequest{
		Index: []string{index},
		Name:  index + "_" + alias,
	}, ctx, ESClient.GetClient().Transport)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Elasticsearch) InsertLog(ctx context.Context, tenantId, index, appId, docId string, esDoc map[string]any) error {
	body, err := json.Marshal(esDoc)
	if err != nil {
		return err
	}
	req := esapi.CreateRequest{
		Index:      index,
		DocumentID: docId,
		Routing:    index + "_" + tenantId,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, ESClient.GetClient())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("insert: response: %s", res.String())
	}

	return nil
}

func (repo *Elasticsearch) UpdateDocById(ctx context.Context, index, appId, docId string, esDoc map[string]any) error {
	body, err := json.Marshal(esDoc)
	if err != nil {
		return err
	}
	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: docId,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, body))),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, ESClient.GetClient())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("update: response: %s", res.String())
	}

	return nil
}

func (repo *Elasticsearch) BulkUpdateDoc(ctx context.Context, index string, esDoc map[string]any) error {
	body, err := json.Marshal(esDoc)
	if err != nil {
		return err
	}
	req := esapi.BulkRequest{
		Index: index,
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(ctx, ESClient.GetClient())
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk update: response: %s", res.String())
	}

	return nil
}

func (repo *Elasticsearch) DeleteById(ctx context.Context, index, docId string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: docId,
	}

	res, err := req.Do(ctx, ESClient.GetClient())
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("delete: response: %s", res.String())
	}

	return nil
}
