package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/internal/elasticsearch"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IInboxMarketingES interface {
		GetDocById(ctx context.Context, tenantId, index, id string) (*model.InboxMarketingLogInfo, error)
		GetDocByRoutingExternalMessageId(ctx context.Context, tenantId, index, externalMessageId string) (*model.InboxMarketingLogInfo, error)
		GetLogCurrentDay(ctx context.Context, index, plugin, startTime, endTime string) (int, []model.InboxMarketingLogInfo, error)
		GetReport(ctx context.Context, tenantId, index string, limit, offset int, filter model.InboxMarketingFilter) (int, []model.InboxMarketingLogReport, error)
	}
	InboxMarketingES struct{}
)

var InboxMarketingESRepo IInboxMarketingES

func NewInboxMarketingES() IInboxMarketingES {
	return &InboxMarketingES{}
}

func (repo *InboxMarketingES) GetDocById(ctx context.Context, tenantId, index, id string) (*model.InboxMarketingLogInfo, error) {
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

func (repo *InboxMarketingES) GetDocByRoutingExternalMessageId(ctx context.Context, tenantId, index, externalMessageId string) (*model.InboxMarketingLogInfo, error) {
	filters := []map[string]any{}
	musts := []map[string]any{}
	if len(tenantId) > 0 {
		filters = append(filters, elasticsearch.TermQuery("_routing", index+"_"+tenantId))
	}
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

func (repo *InboxMarketingES) GetLogCurrentDay(ctx context.Context, index, plugin, startTime, endTime string) (int, []model.InboxMarketingLogInfo, error) {
	filters := []map[string]any{}
	musts := []map[string]any{}

	filters = append(filters, elasticsearch.MatchQuery("plugin", plugin))
	filters = append(filters, elasticsearch.RangeQuery("count_action", 0, 5))
	filters = append(filters, elasticsearch.RangeQuery("created_at", startTime, nil))
	filters = append(filters, elasticsearch.RangeQuery("created_at", nil, endTime))

	musts = append(musts, map[string]any{
		"bool": map[string]any{
			"must": map[string]any{
				"term": map[string]any{
					"is_check": false,
				},
			},
		},
	})

	boolQuery := map[string]any{
		"bool": map[string]any{
			"filter": filters,
			"must":   musts,
		},
	}
	searchSource := map[string]any{
		"size":    10000,
		"_source": true,
		"query":   boolQuery,
		"sort": []any{
			elasticsearch.Order("created_at", false),
		},
	}
	buf, err := elasticsearch.EncodeAny(searchSource)

	if err != nil {
		return 0, nil, err
	}
	client := ES.GetClient()
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

	body := model.ElasticsearchInboxMarketingResponse{}

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return 0, nil, err
	}
	result := []model.InboxMarketingLogInfo{}
	total := body.Hits.Total.Value
	// mapping
	for _, bodyHits := range body.Hits.Hits {
		data := model.InboxMarketingLogInfo{}
		if err := util.ParseAnyToAny(bodyHits.Source, &data); err != nil {
			return 0, nil, err
		}
		result = append(result, data)
	}
	return total, result, nil
}

func (repo *InboxMarketingES) GetReport(ctx context.Context, tenantId, index string, limit, offset int, filter model.InboxMarketingFilter) (int, []model.InboxMarketingLogReport, error) {
	filters := []map[string]any{}
	musts := []map[string]any{}
	if len(tenantId) > 0 {
		filters = append(filters, elasticsearch.TermQuery("_routing", index+"_"+tenantId))
		musts = append(musts, elasticsearch.MatchQuery("tenant_id", tenantId))
	}
	if len(filter.Channel) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("channel", util.ParseToAnyArray(filter.Channel)...))
	}
	if len(filter.ErrorCode) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("error_code", util.ParseToAnyArray(filter.ErrorCode)...))
	}
	if len(filter.Ext) > 0 {
		filters = append(filters, elasticsearch.MatchQuery("ext", filter.Ext))
	}
	if len(filter.RoutingConfigUuid) > 0 {
		filters = append(filters, elasticsearch.MatchQuery("routing_config_uuid", filter.RoutingConfigUuid))
	}
	if len(filter.Id) > 0 {
		filters = append(filters, elasticsearch.MatchQuery("id", filter.Id))
	}
	if len(filter.Message) > 0 {
		filters = append(filters, elasticsearch.MatchQuery("message", filter.Message))
	}
	if len(filter.PhoneNumber) > 0 {
		filters = append(filters, elasticsearch.MatchQuery("phone_number", filter.PhoneNumber))
	}
	if len(filter.RouteRule) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("route_rule", util.ParseToAnyArray(filter.RouteRule)...))
	}
	if len(filter.ServiceTypeId) > 0 {
		serviceTypes := []string{}
		for _, serviceId := range filter.ServiceTypeId {
			telco := strconv.Itoa(serviceId)
			serviceTypes = append(serviceTypes, telco)
		}
		filters = append(filters, elasticsearch.TermsQuery("service_type_id", util.ParseToAnyArray(serviceTypes)...))
	}
	if len(filter.Status) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("status", util.ParseToAnyArray(filter.Status)...))
	}
	if len(filter.TelcoId) > 0 {
		telcos := []string{}
		for _, telcoId := range filter.TelcoId {
			telco := strconv.Itoa(telcoId)
			telcos = append(telcos, telco)
		}
		filters = append(filters, elasticsearch.TermsQuery("telco_id.keyword", util.ParseToAnyArray(telcos)...))
	}
	if len(filter.TemplateCode) > 0 {
		filters = append(filters, elasticsearch.MatchQuery("template_code", filter.TemplateCode))
	}
	if len(filter.CampaignUuid) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("campaign_uuid", util.ParseToAnyArray(filter.CampaignUuid)...))
	}
	if len(filter.Plugin) > 0 {
		musts = append(musts, elasticsearch.TermsQuery("plugin", util.ParseToAnyArray(filter.Plugin)...))
	}
	if filter.IsChargedZns.Valid {
		isChargedZns := map[string]any{
			"bool": map[string]any{
				"filter": map[string]any{
					"term": map[string]any{
						"is_charged_zns": filter.IsChargedZns.Bool,
					},
				},
			},
		}
		filters = append(filters, isChargedZns)
	}
	if len(filter.Quantity) > 0 {
		filters = append(filters, elasticsearch.MatchQuery("quantity", filter.Quantity))
	}
	if len(filter.StartTime) > 0 {
		filters = append(filters, elasticsearch.RangeQuery("created_at", filter.StartTime, nil))
	}
	if len(filter.EndTime) > 0 {
		filters = append(filters, elasticsearch.RangeQuery("created_at", nil, filter.EndTime))
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
			elasticsearch.Order("created_at", false),
		},
	}
	buf, err := elasticsearch.EncodeAny(searchSource)

	if err != nil {
		return 0, nil, err
	}
	client := ES.GetClient()
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

	body := model.ElasticsearchInboxMarketingResponse{}

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return 0, nil, err
	}
	result := []model.InboxMarketingLogReport{}
	total := body.Hits.Total.Value
	// mapping
	for _, bodyHits := range body.Hits.Hits {
		data := model.InboxMarketingLogReport{}
		if err := util.ParseAnyToAny(bodyHits.Source, &data); err != nil {
			return 0, nil, err
		}
		result = append(result, data)
	}

	return total, result, nil
}
