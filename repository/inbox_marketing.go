package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/internal/elasticsearch"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IInboxMarketingES interface {
		GetDocById(ctx context.Context, index, id string) (*model.InboxMarketingLogInfo, error)
		GetDocByRoutingExternalMessageId(ctx context.Context, index, externalMessageId string) (*model.InboxMarketingLogInfo, error)
	}
	InboxMarketingES struct{}
)

var InboxMarketingESRepo IInboxMarketingES

func NewInboxMarketingES() IInboxMarketingES {
	return &InboxMarketingES{}
}

func (repo *InboxMarketingES) GetDocByRoutingExternalMessageId(ctx context.Context, index, externalMessageId string) (*model.InboxMarketingLogInfo, error) {
	filters := []map[string]any{}
	musts := []map[string]any{}
	filters = append(filters, elasticsearch.MatchQuery("external_message_id", externalMessageId))

	boolQuery := map[string]any{
		"bool": map[string]any{
			"filter": filters,
			"must":   musts,
		},
	}
	searchSource := map[string]any{
		"from":    0,
		"size":    1,
		"_source": true,
		"query":   boolQuery,
	}
	// tmp, _ := util.ParseMapToString(searchSource)
	// log.Info(tmp)
	buf, err := elasticsearch.EncodeAny(searchSource)

	if err != nil {
		return nil, err
	}
	client := ES.GetClient()
	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(index),
		client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}

	// handle res error
	if res.IsError() {
		var e map[string]any
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		} else {
			// Print the response status and error information.
			return nil, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]any)["type"],
				e["error"].(map[string]any)["reason"],
			)
		}
	}

	defer res.Body.Close()

	body := model.ElasticsearchInboxMarketingResponse{}

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}
	result := model.InboxMarketingLogInfo{}
	// mapping
	for _, bodyHits := range body.Hits.Hits {
		data := model.InboxMarketingLogInfo{}
		if err := util.ParseAnyToAny(bodyHits.Source, &data); err != nil {
			return nil, err
		}
		result = data
	}
	return &result, nil
}

func (repo *InboxMarketingES) GetDocById(ctx context.Context, index, id string) (*model.InboxMarketingLogInfo, error) {
	filters := []map[string]any{}
	musts := []map[string]any{}
	filters = append(filters, elasticsearch.MatchQuery("_id", id))

	boolQuery := map[string]any{
		"bool": map[string]any{
			"filter": filters,
			"must":   musts,
		},
	}
	searchSource := map[string]any{
		"from":    0,
		"size":    1,
		"_source": true,
		"query":   boolQuery,
		"collapse": map[string]any{
			"field": "routing_config_uuid",
			"inner_hits": []map[string]any{
				map[string]interface{}{
					"name": "hit_key",
					"size": 1,
					"sort": []map[string]any{
						{"created_at": map[string]any{"order": "desc"}},
					},
				},
			},
		},
	}
	buf, err := elasticsearch.EncodeAny(searchSource)

	if err != nil {
		return nil, err
	}
	client := ES.GetClient()
	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(index),
		client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}

	// handle res error
	if res.IsError() {
		var e map[string]any
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		} else {
			// Print the response status and error information.
			return nil, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]any)["type"],
				e["error"].(map[string]any)["reason"],
			)
		}
	}

	defer res.Body.Close()

	body := model.ElasticsearchInboxMarketingResponse{}

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}
	result := model.InboxMarketingLogInfo{}
	// mapping
	for _, bodyHits := range body.Hits.Hits {
		data := model.InboxMarketingLogInfo{}
		if err := util.ParseAnyToAny(bodyHits.Source, &data); err != nil {
			return nil, err
		}
		result = data
	}
	return &result, nil
}
