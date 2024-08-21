package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/internal/elasticsearch"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IConversationES interface {
		IESGenericRepo[model.Conversation]
		GetConversations(ctx context.Context, tenantId, index string, filter model.ConversationFilter, limit, offset int) (int, *[]model.ConversationView, error)
		GetConversationById(ctx context.Context, tenantId, index, appId, id string) (*model.Conversation, error)
		SearchWithScroll(ctx context.Context, tenantId, index string, filter model.ConversationFilter, size int, scrollId string, scrollDurations ...time.Duration) (total int, entries []*model.ConversationView, respScrollId string, err error)
		GetNotesList(ctx context.Context, tenantId, index string, filter model.ConversationNotesListFilter, limit, offset int) (total int, entries []*model.NotesList, err error)
	}
	ConversationES struct {
		ESGenericRepo[model.Conversation]
	}
)

var ConversationESRepo IConversationES

func NewConversationES() IConversationES {
	return &ConversationES{}
}

/*
 * script used to filter notes list in ES query
 * int limit = params.containsKey('limit') ? (int)params['limit'] : 10;
 * int offset = params.containsKey('offset') ? (int)params['offset'] : 0;
 * def notes = params['_source'].containsKey('notes_list') && params['_source']['notes_list'] != null ? params['_source']['notes_list'] : [];
 * int totalSize = notes.size();
 * def slicedNotes = [];
 * if (totalSize > 0) {
 *    notes.sort((a, b) -> b['created_at'].compareTo(a['created_at']));
 *    int end = (int) Math.min(totalSize, offset + limit);
 *    slicedNotes = totalSize > offset ? notes.subList(offset, end) : [];
 * }
 * return ['notes_list': slicedNotes, 'total_size': totalSize];
 */
const notesListScript = "int limit = params.containsKey('limit') ? (int)params['limit'] : 10; int offset = params.containsKey('offset') ? (int)params['offset'] : 0; def notes = params['_source'].containsKey('notes_list') && params['_source']['notes_list'] != null ? params['_source']['notes_list'] : []; int totalSize = notes.size(); def slicedNotes = []; if (totalSize > 0) { notes.sort((a, b) -> b['created_at'].compareTo(a['created_at'])); int end = (int) Math.min(totalSize, offset + limit); slicedNotes = totalSize > offset ? notes.subList(offset, end) : []; } return ['notes_list': slicedNotes, 'total_size': totalSize];"

func (repo *ConversationES) GetConversations(ctx context.Context, tenantId, index string, filter model.ConversationFilter, limit, offset int) (int, *[]model.ConversationView, error) {
	filters := []map[string]any{}
	musts := []map[string]any{}
	insensitiveBool, _ := strconv.ParseBool(filter.Insensitive)
	insensitive := sql.NullBool{
		Bool:  insensitiveBool,
		Valid: true,
	}

	// Remove because routing maybe having pitel_conversation_
	if len(tenantId) > 0 {
		musts = append(musts, elasticsearch.MatchQuery("_routing", index+"_"+tenantId))
	}
	if len(filter.TenantId) > 0 {
		filters = append(filters, elasticsearch.MatchQuery("tenant_id", filter.TenantId))
	}
	if len(filter.AppId) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("app_id", util.ParseToAnyArray(filter.AppId)...))
	}
	if len(filter.ConversationId) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("conversation_id", util.ParseToAnyArray(filter.ConversationId)...))
	}
	if len(filter.ExternalConversationId) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("external_conversation_id", util.ParseToAnyArray(filter.ExternalConversationId)...))
	}
	if len(filter.Username) > 0 {
		// Search like
		filters = append(filters, elasticsearch.WildcardQuery("username", "*"+filter.Username, insensitive))
	}
	if len(filter.PhoneNumber) > 0 {
		filters = append(filters, elasticsearch.WildcardQuery("phone_number", "*"+filter.PhoneNumber, sql.NullBool{}))
	}
	if len(filter.Email) > 0 {
		filters = append(filters, elasticsearch.WildcardQuery("email", "*"+filter.Email, insensitive))
	}
	if filter.IsDone.Valid {
		bq := map[string]any{
			"bool": map[string]any{
				"filter": []map[string]any{
					{
						"bool": map[string]any{
							"must": map[string]any{
								"wildcard": map[string]any{
									"is_done": strconv.FormatBool(filter.IsDone.Bool),
								},
							},
						},
					},
				},
			},
		}
		filters = append(filters, bq)
	}
	if filter.Major.Valid {
		bq := map[string]any{
			"bool": map[string]any{
				"filter": []map[string]any{
					{
						"bool": map[string]any{
							"must": map[string]any{
								"term": map[string]any{
									"major": strconv.FormatBool(filter.Major.Bool),
								},
							},
						},
					},
				},
			},
		}
		filters = append(filters, bq)
	}
	if filter.Following.Valid {
		bq := map[string]any{
			"bool": map[string]any{
				"filter": []map[string]any{
					{
						"bool": map[string]any{
							"must": map[string]any{
								"term": map[string]any{
									"following": strconv.FormatBool(filter.Following.Bool),
								},
							},
						},
					},
				},
			},
		}
		filters = append(filters, bq)
	}
	scriptFields := map[string]any{
		"notes_list": map[string]any{
			"script": map[string]any{
				"source": notesListScript,
				"params": map[string]any{
					"limit":  2,
					"offset": 0,
				},
			},
		},
	}

	boolQuery := map[string]any{
		"bool": map[string]any{
			"filter": filters,
			"must":   musts,
		},
	}
	searchSource := map[string]any{
		"from":          offset,
		"size":          limit,
		"_source":       true,
		"script_fields": scriptFields,
		"query":         boolQuery,
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
	result := []model.ConversationView{}
	total := body.Hits.Total.Value
	// mapping
	for _, bodyHits := range body.Hits.Hits {
		data := model.ConversationView{}
		if err := util.ParseAnyToAny(bodyHits.Source, &data); err != nil {
			return 0, nil, err
		}
		// replace notes list in _source with the one in script_fields to keep only 2 latest items
		if len(bodyHits.Fields.NotesList) > 0 {
			hitData, ok := bodyHits.Fields.NotesList[0].(map[string]any)
			if !ok {
				err = errors.New("failed to convert notes list")
				return 0, nil, err
			}
			notesList := &[]model.NotesList{}
			if err = util.ParseAnyToAny(hitData["notes_list"], &notesList); err != nil {
				log.Error(err)
				return 0, nil, err
			}
			data.NotesList = notesList
		}
		result = append(result, data)
	}

	return total, &result, nil
}

func (repo *ConversationES) GetConversationById(ctx context.Context, tenantId, index, appId, id string) (*model.Conversation, error) {
	filters := []map[string]any{}
	musts := []map[string]any{}
	if len(tenantId) > 0 {
		musts = append(musts, elasticsearch.MatchQuery("_routing", index+"_"+tenantId))
		filters = append(filters, elasticsearch.MatchQuery("tenant_id", tenantId))
	}
	if len(appId) > 0 {
		filters = append(filters, elasticsearch.MatchQuery("app_id", appId))
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

		// replace notes list in _source with the one in script_fields to keep only 2 latest items
		if len(bodyHits.Fields.NotesList) > 0 {
			hitData, ok := bodyHits.Fields.NotesList[0].(map[string]any)
			if !ok {
				err = errors.New("failed to convert notes list")
				return nil, err
			}
			notesList := &[]model.NotesList{}
			if err = util.ParseAnyToAny(hitData["notes_list"], &notesList); err != nil {
				log.Error(err)
				return nil, err
			}
			data.NotesList = notesList
		}

		result = data
	}
	return &result, nil
}

func (repo *ConversationES) SearchWithScroll(ctx context.Context, tenantId, index string, filter model.ConversationFilter, size int, scrollId string, scrollDurations ...time.Duration) (total int, entries []*model.ConversationView, respScrollId string, err error) {
	var body *model.SearchReponse
	if len(scrollId) < 1 {
		scrollDuration := 5 * time.Minute
		if len(scrollDurations) > 0 {
			scrollDuration = scrollDurations[0]
		}
		body, err = repo.searchWithScroll(ctx, tenantId, index, filter, size, scrollDuration)
	} else {
		body, err = repo.ScrollAPI(ctx, scrollId)
	}
	if err != nil || body == nil {
		return
	}
	total = body.Hits.Total.Value
	hits := body.Hits.Hits
	respScrollId = body.ScrollId
	entries = make([]*model.ConversationView, 0)
	for _, hit := range hits {
		entry := &model.ConversationView{}
		if err = util.ParseAnyToAny(hit.Source, entry); err != nil {
			return
		}

		// replace notes list in _source with the one in script_fields to keep only 2 latest items
		if len(hit.Fields.NotesList) > 0 {
			hitData, ok := hit.Fields.NotesList[0].(map[string]any)
			if !ok {
				err = errors.New("failed to convert notes list")
				return
			}
			notesList := &[]model.NotesList{}
			if err = util.ParseAnyToAny(hitData["notes_list"], &notesList); err != nil {
				log.Error(err)
				return
			}
			entry.NotesList = notesList
		}
		entries = append(entries, entry)
	}
	return total, entries, respScrollId, nil
}

func (repo *ConversationES) searchWithScroll(ctx context.Context, tenantId, index string, filter model.ConversationFilter, size int, scrollDuration time.Duration) (result *model.SearchReponse, err error) {
	filters := []map[string]any{}
	musts := []map[string]any{}
	insensitiveBool, _ := strconv.ParseBool(filter.Insensitive)
	insensitive := sql.NullBool{
		Bool:  insensitiveBool,
		Valid: true,
	}

	// Remove because routing maybe having pitel_conversation_
	// filters = append(filters, elasticsearch.TermQuery("_routing", index+"_"+tenantId))
	if len(tenantId) > 0 {
		musts = append(musts, elasticsearch.MatchQuery("_routing", index+"_"+tenantId))
		//filters = append(filters, elasticsearch.MatchQuery("tenant_id", tenantId))
	}
	if len(filter.AppId) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("app_id", util.ParseToAnyArray(filter.AppId)...))
	}
	if len(filter.ConversationId) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("conversation_id", util.ParseToAnyArray(filter.ConversationId)...))
	}
	if len(filter.Username) > 0 {
		// Search like
		filters = append(filters, elasticsearch.WildcardQuery("username", "*"+filter.Username, insensitive))
	}
	if len(filter.PhoneNumber) > 0 {
		filters = append(filters, elasticsearch.WildcardQuery("phone_number", "*"+filter.PhoneNumber, sql.NullBool{}))
	}
	if len(filter.Email) > 0 {
		filters = append(filters, elasticsearch.WildcardQuery("email", "*"+filter.Email, insensitive))
	}
	if filter.IsDone.Valid {
		bq := map[string]any{
			"bool": map[string]any{
				"filter": []map[string]any{
					{
						"bool": map[string]any{
							"must": map[string]any{
								"wildcard": map[string]any{
									"is_done": strconv.FormatBool(filter.IsDone.Bool),
								},
							},
						},
					},
				},
			},
		}
		filters = append(filters, bq)
	}
	if filter.Major.Valid {
		bq := map[string]any{
			"bool": map[string]any{
				"filter": []map[string]any{
					{
						"bool": map[string]any{
							"must": map[string]any{
								"term": map[string]any{
									"major": strconv.FormatBool(filter.Major.Bool),
								},
							},
						},
					},
				},
			},
		}
		filters = append(filters, bq)
	}
	if filter.Following.Valid {
		bq := map[string]any{
			"bool": map[string]any{
				"filter": []map[string]any{
					{
						"bool": map[string]any{
							"must": map[string]any{
								"term": map[string]any{
									"following": strconv.FormatBool(filter.Following.Bool),
								},
							},
						},
					},
				},
			},
		}
		filters = append(filters, bq)
	}

	boolQuery := map[string]any{
		"bool": map[string]any{
			"filter": filters,
			"must":   musts,
		},
	}
	scriptFields := map[string]any{
		"notes_list": map[string]any{
			"script": map[string]any{
				"source": notesListScript,
				"params": map[string]any{
					"limit":  2,
					"offset": 0,
				},
			},
		},
	}

	// search
	searchSource := map[string]any{
		"query":         boolQuery,
		"script_fields": scriptFields,
		"_source":       true,
		"size":          0,
		"sort": []any{
			elasticsearch.Order("updated_at", false),
			elasticsearch.Order("created_at", false),
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

func (repo *ConversationES) GetNotesList(ctx context.Context, tenantId, index string, filter model.ConversationNotesListFilter, limit, offset int) (total int, entries []*model.NotesList, err error) {
	filters := []map[string]any{}
	musts := []map[string]any{}

	// Remove because routing maybe having pitel_conversation_
	// filters = append(filters, elasticsearch.TermQuery("_routing", index+"_"+tenantId))
	if len(tenantId) > 0 {
		musts = append(musts, elasticsearch.MatchQuery("_routing", index+"_"+tenantId))
		//filters = append(filters, elasticsearch.MatchQuery("tenant_id", tenantId))
	}
	if len(filter.AppId) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("app_id", util.ParseToAnyArray([]string{filter.AppId})...))
	}
	if len(filter.OaId) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("oa_id", util.ParseToAnyArray([]string{filter.OaId})...))
	}
	if len(filter.ConversationId) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("conversation_id", util.ParseToAnyArray([]string{filter.ConversationId})...))
	}

	boolQuery := map[string]any{
		"bool": map[string]any{
			"filter": filters,
			"must":   musts,
		},
	}
	scriptFields := map[string]any{
		"notes_list": map[string]any{
			"script": map[string]any{
				"source": notesListScript,
				"params": map[string]any{
					"limit":  limit,
					"offset": offset,
				},
			},
		},
	}

	// search
	searchSource := map[string]any{
		"query":         boolQuery,
		"script_fields": scriptFields,
		"_source":       false,
		"size":          1,
		"sort":          []any{},
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
	)
	if err != nil {
		return
	}
	defer res.Body.Close()
	result, err := elasticsearch.ParseSearchResponse((*esapi.Response)(res))
	if err != nil {
		return
	}
	if result == nil || len(result.Hits.Hits) < 1 {
		return
	}
	hit := result.Hits.Hits[0]
	if len(hit.Fields.NotesList) < 1 {
		err = errors.New("not found any notes list")
		return
	}
	hitData, ok := hit.Fields.NotesList[0].(map[string]any)
	if !ok {
		err = errors.New("failed to convert notes list")
		return
	}
	total = int(hitData["total_size"].(float64))
	entries = make([]*model.NotesList, 0)
	if err = util.ParseAnyToAny(hitData["notes_list"], &entries); err != nil {
		return
	}
	return
}
