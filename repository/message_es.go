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
	IMessageES interface {
		GetMessages(ctx context.Context, tenantId, index string, filter model.MessageFilter, limit, offset int) (int, *[]model.Message, error)
		GetMessageById(ctx context.Context, tenantId, index, id string) (*model.Message, error)
	}
	MessageES struct {
	}
)

var MessageESRepo IMessageES

func NewMessageES() IMessageES {
	return &MessageES{}
}

func (m *MessageES) GetMessages(ctx context.Context, tenantId, index string, filter model.MessageFilter, limit, offset int) (int, *[]model.Message, error) {
	filters := []map[string]any{}
	musts := []map[string]any{}
	if len(tenantId) > 0 {
		filters = append(filters, elasticsearch.TermQuery("_routing", index+"_"+tenantId))
	}
	if len(filter.AppId) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("app_id", util.ParseToAnyArray([]string{filter.AppId})...))
	}
	if len(filter.ConversationId) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("conversation_id", util.ParseToAnyArray([]string{filter.ConversationId})...))
	}
	if len(filter.IsRead) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("is_read", util.ParseToAnyArray([]string{filter.IsRead})...))
	}

	boolQuery := map[string]any{
		"bool": map[string]any{
			"filter": filters,
			"must":   musts,
		},
	}
	searchSource := map[string]any{
		"from":    offset,
		"size":    limit,
		"_source": true,
		"query":   boolQuery,
		"sort": []any{
			elasticsearch.Order("send_time", false),
		},
	}
	buf, err := elasticsearch.EncodeAny(searchSource)

	if err != nil {
		return 0, nil, err
	}
	client := ESClient.GetClient()
	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(index),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
	)
	if err != nil {
		return 0, nil, err
	}

	// handle res error
	if res.IsError() {
		var e map[string]any
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return 0, nil, err
		} else {
			// Print the response status and error information.
			return 0, nil, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]any)["type"],
				e["error"].(map[string]any)["reason"],
			)
		}
	}

	defer res.Body.Close()

	body := model.ElasticsearchChatResponse{}

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return 0, nil, err
	}
	result := []model.Message{}
	total := body.Hits.Total.Value
	// mapping
	for _, bodyHits := range body.Hits.Hits {
		data := model.Message{}
		if err := util.ParseAnyToAny(bodyHits.Source, &data); err != nil {
			return 0, nil, err
		}
		result = append(result, data)
	}

	return total, &result, nil

}

func (repo *MessageES) GetMessageById(ctx context.Context, tenantId, index, id string) (*model.Message, error) {
	filters := []map[string]any{}
	musts := []map[string]any{}
	if len(tenantId) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("_routing", index+"_"+tenantId))
		musts = append(musts, elasticsearch.MatchQuery("tenant_id", tenantId))
	}
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
		"sort": []any{
			elasticsearch.Order("send_time", false),
		},
	}
	buf, err := elasticsearch.EncodeAny(searchSource)

	if err != nil {
		return nil, err
	}
	client := ESClient.GetClient()
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

	body := model.ElasticsearchChatResponse{}

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}
	result := model.Message{}
	// mapping
	for _, bodyHits := range body.Hits.Hits {
		data := model.Message{}
		if err := util.ParseAnyToAny(bodyHits.Source, &data); err != nil {
			return nil, err
		}
		result = data
	}
	return &result, nil
}
