package elasticsearchsearch

import (
	"net/http"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Username              string
	Password              string
	Host                  []string
	RetryStatuses         []int
	MaxRetries            int
	ResponseHeaderTimeout int
}

type Response struct {
	StatusCode int
	Header     http.Header
	Body       map[string]interface{}
}

type IElasticsearchClient interface {
	GetClient() *elasticsearch.Client
}

type elasticsearchClient struct {
	config Config
	client *elasticsearch.Client
}

func NewElasticsearchClient(config Config) IElasticsearchClient {
	es := &elasticsearchClient{
		config: config,
	}
	if err := es.Connect(); err != nil {
		log.Fatal(err)
		return nil
	}
	_, err := es.GetClient().Ping()
	if err != nil {
		log.Fatal("Elasticsearch connection failed")
		return nil
	} else {
		log.Info("Elasticsearch connection successful")
	}
	return es
}

func (e *elasticsearchClient) Connect() error {
	client, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses:     e.config.Host,
			Username:      e.config.Username,
			Password:      e.config.Password,
			RetryOnStatus: e.config.RetryStatuses,
			MaxRetries:    e.config.MaxRetries,
		},
	)
	if err != nil {
		return err
	}
	e.client = client
	return nil
}

// func (e *elasticsearchClient) Ping() error {
// 	e.client.Ping()
// }

func (e *elasticsearchClient) GetClient() *elasticsearch.Client {
	return e.client
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
