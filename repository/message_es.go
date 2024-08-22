package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/internal/elasticsearch"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IMessageES interface {
		IESGenericRepo[model.Message]
		GetMessages(ctx context.Context, tenantId, index string, filter model.MessageFilter, limit, offset int) (int, *[]model.Message, error)
		GetMessageById(ctx context.Context, tenantId, index, id string) (*model.Message, error)
		SearchWithScroll(ctx context.Context, tenantId, index string, filter model.MessageFilter, limit int, scrollId string, scrollDurations ...time.Duration) (total int, entries []*model.Message, respScrollId string, err error)
		GetMessageMediasWithScroll(ctx context.Context, tenantId, index string, filter model.MessageFilter, limit int, scrollId string, scrollDurations ...time.Duration) (total int, entries []*model.MessageAttachmentsDetails, respScrollId string, err error)
	}
	MessageES struct {
		ESGenericRepo[model.Message]
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
		musts = append(musts, elasticsearch.MatchQuery("_routing", index+"_"+tenantId))
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
	if len(filter.ExternalMessageId) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("external_message_id", util.ParseToAnyArray([]string{filter.ExternalMessageId})...))
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
		musts = append(musts, elasticsearch.MatchQuery("_routing", index+"_"+tenantId))
		filters = append(filters, elasticsearch.MatchQuery("tenant_id", tenantId))
	}
	musts = append(musts, elasticsearch.MatchQuery("_id", id))

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

func (repo *MessageES) SearchWithScroll(ctx context.Context, tenantId, index string, filter model.MessageFilter, limit int, scrollId string, scrollDurations ...time.Duration) (total int, entries []*model.Message, respScrollId string, err error) {
	var body *model.SearchReponse
	if len(scrollId) < 1 {
		scrollDuration := 5 * time.Minute
		if len(scrollDurations) > 0 {
			scrollDuration = scrollDurations[0]
		}
		body, err = repo.searchWithScroll(ctx, tenantId, index, filter, limit, scrollDuration)
	} else {
		body, err = repo.ScrollAPI(ctx, scrollId)
	}
	if err != nil || body == nil {
		return
	}
	total = body.Hits.Total.Value
	hits := body.Hits.Hits
	respScrollId = body.ScrollId
	entries = make([]*model.Message, 0)
	for _, hit := range hits {
		entry := &model.Message{}
		if err = util.ParseAnyToAny(hit.Source, entry); err != nil {
			return
		}
		entries = append(entries, entry)
	}
	return total, entries, respScrollId, nil
}

func (repo *MessageES) searchWithScroll(ctx context.Context, tenantId, index string, filter model.MessageFilter, size int, scrollDuration time.Duration) (result *model.SearchReponse, err error) {
	filters := []map[string]any{}
	musts := []map[string]any{}
	if len(tenantId) > 0 {
		musts = append(musts, elasticsearch.MatchQuery("_routing", index+"_"+tenantId))
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
	if len(filter.ExternalMessageId) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("external_message_id", util.ParseToAnyArray([]string{filter.ExternalMessageId})...))
	}

	boolQuery := map[string]any{
		"bool": map[string]any{
			"filter": filters,
			"must":   musts,
		},
	}

	// search
	searchSource := map[string]any{
		"query":   boolQuery,
		"_source": true,
		"size":    0,
		"sort": []any{
			elasticsearch.Order("send_time", false),
		},
	}
	if size > 0 {
		searchSource["size"] = size
	}

	buf, err := elasticsearch.EncodeAny(searchSource)
	if err != nil {
		return
	}
	client := ESClient.GetClient()
	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(index),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
		client.Search.WithScroll(scrollDuration),
	)
	if err != nil {
		return
	}
	defer res.Body.Close()
	result, err = elasticsearch.ParseSearchResponse((*esapi.Response)(res))
	return
}

func (repo *MessageES) GetMessageMediasWithScroll(ctx context.Context, tenantId, index string, filter model.MessageFilter, limit int, scrollId string, scrollDurations ...time.Duration) (total int, entries []*model.MessageAttachmentsDetails, respScrollId string, err error) {
	var body *model.SearchReponse
	if len(scrollId) < 1 {
		scrollDuration := 5 * time.Minute
		if len(scrollDurations) > 0 {
			scrollDuration = scrollDurations[0]
		}
		body, err = repo.searchMediasWithScroll(ctx, tenantId, index, filter, limit, scrollDuration)
	} else {
		body, err = repo.ScrollAPI(ctx, scrollId)
	}
	if err != nil || body == nil {
		return
	}
	total = body.Hits.Total.Value
	hits := body.Hits.Hits
	respScrollId = body.ScrollId
	entries = make([]*model.MessageAttachmentsDetails, 0)

	// attachment hits
	for _, messageHit := range hits {
		messageEntry := &model.Message{}
		if err = util.ParseAnyToAny(messageHit.Source, messageEntry); err != nil {
			log.Error(err)
			return
		}

		for _, attachmentHit := range messageEntry.Attachments {
			if attachmentHit == nil {
				log.Error("not found attachment")
				continue
			}
			entry := &model.MessageAttachmentsDetails{}
			if err = util.ParseAnyToAny(attachmentHit, entry); err != nil {
				log.Error(err)
				return
			}
			entry.MessageId = messageEntry.MessageId
			entry.SendTime = time.UnixMilli(messageEntry.SendTimestamp)
			entries = append(entries, entry)
		}
	}

	if filter.AttachmentType == "link" || filter.AttachmentType == "" {
		// Find all URLs in the content
		urlRegex := regexp.MustCompile(`https?://[^\s]+`)
		for _, messageHit := range hits {
			messageEntry := &model.Message{}
			if err = util.ParseAnyToAny(messageHit.Source, messageEntry); err != nil {
				log.Error(err)
				return
			}

			urls := urlRegex.FindAllString(messageEntry.Content, -1)
			for _, url := range urls {
				entry := &model.MessageAttachmentsDetails{
					AttachmentType: "link",
					Payload:        model.OttPayloadMedia{Url: url},
					MessageId:      messageEntry.MessageId,
					MessageContent: messageEntry.Content,
					SendTime:       time.UnixMilli(messageEntry.SendTimestamp),
				}
				entries = append(entries, entry)
			}
		}
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].SendTime.After(entries[j].SendTime)
		})
	}
	return total, entries, respScrollId, nil
}

func (repo *MessageES) searchMediasWithScroll(ctx context.Context, tenantId, index string, filter model.MessageFilter, size int, scrollDuration time.Duration) (result *model.SearchReponse, err error) {
	filters := []map[string]any{}
	musts := []map[string]any{}
	shoulds := []map[string]any{}
	if len(tenantId) > 0 {
		musts = append(musts, elasticsearch.MatchQuery("_routing", index+"_"+tenantId))
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
	if len(filter.ExternalMessageId) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("external_message_id", util.ParseToAnyArray([]string{filter.ExternalMessageId})...))
	}
	if filter.AttachmentType == "link" || filter.AttachmentType == "" {
		regexpString := ".*(http|https)://.*"
		if len(filter.SearchKeyword) > 0 {
			regexpString = ".*(http|https)://.*" + filter.SearchKeyword + ".*"
		}
		regexpQuery := map[string]any{
			"regexp": map[string]any{
				"content": regexpString,
			},
		}
		shoulds = append(shoulds, regexpQuery)
	}
	nestedMustQuery := filterMediaTypes(filter.AttachmentType)
	if len(filter.SearchKeyword) > 0 {
		wildcardQuery := map[string]any{
			"wildcard": map[string]any{
				"attachments.payload.url": "*" + filter.SearchKeyword + "*",
			},
		}
		nestedMustQuery = append(nestedMustQuery, wildcardQuery)
	}
	nested := map[string]any{
		"nested": map[string]any{
			"path": "attachments",
			"query": map[string]any{
				"bool": map[string]any{
					"must": nestedMustQuery,
				},
			},
		},
	}
	shoulds = append(shoulds, nested)

	boolQuery := map[string]any{
		"bool": map[string]any{
			"filter":               filters,
			"must":                 musts,
			"should":               shoulds,
			"minimum_should_match": 1,
		},
	}

	// search
	searchSource := map[string]any{
		"query":   boolQuery,
		"_source": true,
		"size":    0,
		"sort": []any{
			elasticsearch.Order("send_time", false),
		},
	}
	if size > 0 {
		searchSource["size"] = size
	}

	buf, err := elasticsearch.EncodeAny(searchSource)
	if err != nil {
		return
	}
	client := ESClient.GetClient()
	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(index),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
		client.Search.WithScroll(scrollDuration),
	)
	if err != nil {
		return
	}
	defer res.Body.Close()
	result, err = elasticsearch.ParseSearchResponse((*esapi.Response)(res))
	return
}

func filterMediaTypes(attachmentType string) []any {
	if len(attachmentType) < 1 {
		return []any{
			map[string]any{
				"match_all": map[string]any{},
			},
		}
	}
	args := []string{attachmentType}
	if attachmentType == "media" {
		args = []string{"image", "audio", "video", "sticker", "gif", "reacted", "unreacted"}
	}
	return []any{elasticsearch.TermsQuery("attachments.att_type", util.ParseToAnyArray(args)...)}
}
