package elasticsearch

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/tel4vn/fins-microservices/common/util"

	es "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	log "github.com/sirupsen/logrus"
)

type (
	IESClient interface {
		GetClient() *es.Client
		Ping() error
	}
	ESClient struct {
		config Config
		client *es.Client
	}
)

func NewES(config Config) IESClient {
	e := &ESClient{
		config: config,
	}
	cfg := es.Config{
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Minute * 10,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
		Username:      e.config.Username,
		Password:      e.config.Password,
		Addresses:     e.config.Host,
		MaxRetries:    e.config.MaxRetries,
		RetryOnStatus: e.config.RetryStatuses,
	}
	client, err := es.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	e.client = client
	return e
}

func (e *ESClient) GetClient() *es.Client {
	return e.client
}

func (e *ESClient) Ping() error {
	res, err := e.GetClient().Info()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Error("Error: " + res.String())
	}
	info := map[string]any{}
	if err := DecodeAny(res.Body, &info); err != nil {
		return err
	}
	log.Infof("Elasticsearch returned with code %d and version %s", res.StatusCode, info["version"].(map[string]any)["number"])
	return nil
}

func ParseAnyToAny(value any, dest any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, dest); err != nil {
		return err
	}
	return nil
}

func DecodeAny(value io.Reader, dest any) error {
	if err := json.NewDecoder(value).Decode(dest); err != nil {
		return err
	}
	return nil
}

func EncodeAny(value any) (bytes.Buffer, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(value); err != nil {
		return buf, err
	}
	return buf, nil
}

func RangeQuery(field string, from, to any) map[string]any {
	return map[string]any{
		"range": map[string]any{
			field: map[string]any{
				"from":          from,
				"to":            to,
				"include_lower": true,
				"include_upper": true,
			},
		},
	}
}
func MatchQuery(field string, value any) map[string]any {
	return map[string]any{
		"match": map[string]any{
			field: value,
		},
	}
}

func TermQuery(field string, value any) map[string]any {
	return map[string]any{
		"term": map[string]any{
			field: map[string]any{
				"value": value,
			},
		},
	}
}

func TermsQuery(field string, values ...any) map[string]any {
	arr := make([]any, 0)
	if len(values) > 0 {
		arr = append(arr, values...)
	}
	return map[string]any{
		"terms": map[string]any{
			field: arr,
		},
	}
}

func WildcardQuery(field string, value any) map[string]any {
	return map[string]any{
		"wildcard": map[string]any{
			field: map[string]any{
				"value": value,
			},
		},
	}
}

func ShouldQuery(queries ...map[string]any) map[string]any {
	return map[string]any{
		"should": queries,
	}
}

func BoolQuery(queries ...map[string]any) map[string]any {
	query := map[string]any{
		"bool": map[string]any{},
	}
	if len(queries) == 1 {
		query["bool"] = queries[0]
	} else if len(queries) > 1 {
		var clauses []any
		for _, subQuery := range queries {
			clauses = append(clauses, subQuery)
		}
		query["bool"] = clauses
	}
	return query
}

func MustQuery(queries ...map[string]any) map[string]any {
	query := map[string]any{
		"must": map[string]any{},
	}
	if len(queries) == 1 {
		query["must"] = queries[0]
	} else if len(queries) > 1 {
		var clauses []any
		for _, subQuery := range queries {
			clauses = append(clauses, subQuery)
		}
		query["must"] = clauses
	}
	return query
}

func Order(field string, isAsc bool) map[string]any {
	order := "asc"
	if !isAsc {
		order = "desc"
	}
	return map[string]any{
		field: map[string]any{
			"order": order,
		},
	}
}

func ParseSearchResponse(response *esapi.Response) (*SearchReponse, error) {
	if response.IsError() {
		var e map[string]any
		if err := json.NewDecoder(response.Body).Decode(&e); err != nil {
			return nil, err
		} else {
			typeErr := e["error"].(map[string]any)["type"]
			reason := e["error"].(map[string]any)["reason"]
			return nil, errors.New(util.ParseString(typeErr) + ":" + util.ParseString(reason))
		}
	}
	result := SearchReponse{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return &result, err
	}
	return &result, nil
}

type SearchReponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore interface{} `json:"max_score"`
		Hits     []struct {
			Index   string      `json:"_index"`
			Type    string      `json:"_type"`
			ID      string      `json:"_id"`
			Score   interface{} `json:"_score"`
			Routing string      `json:"_routing"`
			Source  any         `json:"_source"`
			Sort    []string    `json:"sort"`
		} `json:"hits"`
	} `json:"hits"`
}
