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
	IConversationES interface {
		GetConversations(ctx context.Context, tenantId, index string, filter model.ConversationFilter, limit, offset int) (int, *[]model.Conversation, error)
		GetConversationById(ctx context.Context, appId, index, id string) (*model.Conversation, error)
	}
	ConversationES struct {
	}
)

var ConversationESRepo IConversationES

func NewConversationES() IConversationES {
	return &ConversationES{}
}

func (repo *ConversationES) GetConversations(ctx context.Context, tenantId, index string, filter model.ConversationFilter, limit, offset int) (int, *[]model.Conversation, error) {
	filters := []map[string]any{}
	musts := []map[string]any{}
	if len(tenantId) > 0 {
		filters = append(filters, elasticsearch.TermQuery("_routing", index+"_"+tenantId))
		musts = append(musts, elasticsearch.MatchQuery("app_id", tenantId))
	}
	if len(filter.AppId) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("app_id", util.ParseToAnyArray(filter.AppId)...))
	}
	if len(filter.ConversationId) > 0 {
		// filters = append(filters, elasticsearch.TermsQuery("external_user_id", util.ParseToAnyArray(filter.ConversationId)...))
	}
	if len(filter.Username) > 0 {
		// Search like
		filters = append(filters, elasticsearch.MustQuery(map[string]any{
			"multi_match": map[string]any{
				"query":  "%" + filter.Username + "%s",
				"fields": []string{"username"},
			},
		}))
	}
	if len(filter.PhoneNumber) > 0 {
		// filters = append(filters, elasticsearch.TermsQuery("phone_number", util.ParseToAnyArray(filter.PhoneNumber)...))
	}
	if len(filter.Email) > 0 {
		// filters = append(filters, elasticsearch.TermsQuery("email", util.ParseToAnyArray(filter.Email)...))
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
			elasticsearch.Order("updated_at", false),
			elasticsearch.Order("created_at", false),
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
	result := []model.Conversation{}
	total := body.Hits.Total.Value
	// mapping
	for _, bodyHits := range body.Hits.Hits {
		data := model.Conversation{}
		if err := util.ParseAnyToAny(bodyHits.Source, &data); err != nil {
			return 0, nil, err
		}
		result = append(result, data)
	}

	return total, &result, nil
}

func (repo *ConversationES) GetConversationById(ctx context.Context, appId, index, id string) (*model.Conversation, error) {
	filters := []map[string]any{}
	musts := []map[string]any{}
	if len(appId) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("_routing", index+"_"+appId))
		musts = append(musts, elasticsearch.MatchQuery("app_id", appId))
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
			elasticsearch.Order("updated_at", false),
			elasticsearch.Order("created_at", false),
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
	result := model.Conversation{}
	// mapping
	for _, bodyHits := range body.Hits.Hits {
		data := model.Conversation{}
		if err := util.ParseAnyToAny(bodyHits.Source, &data); err != nil {
			return nil, err
		}
		result = data
	}
	return &result, nil
}
