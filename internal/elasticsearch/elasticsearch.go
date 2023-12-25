package elasticsearch

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	elastic "github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Username              string
	Password              string
	Host                  []string
	RetryStatuses         []int
	MaxRetries            int
	ResponseHeaderTimeout int
	Index                 string
}

type Response struct {
	StatusCode int
	Header     http.Header
	Body       map[string]interface{}
}

type IElasticsearchClient interface {
	GetClient() *elastic.Client
}

type ElasticsearchClient struct {
	config Config
	client *elastic.Client
}

func NewElasticsearchClient(config Config) IElasticsearchClient {
	es := &ElasticsearchClient{
		config: config,
	}
	if err := es.Connect(); err != nil {
		log.Fatal(err)
		return nil
	}
	if err := es.Ping(); err != nil {
		log.Fatal(err)
		return nil
	}
	return es
}

func (e *ElasticsearchClient) Connect() error {
	client, err := elastic.NewClient(
		elastic.SetBasicAuth(e.config.Username, e.config.Password),
		elastic.SetURL(strings.Join(e.config.Host[:], ",")),
		elastic.SetSniff(false),
		elastic.SetRetryStatusCodes(e.config.RetryStatuses...),
		elastic.SetMaxRetries(e.config.MaxRetries),
		elastic.SetHealthcheck(true),
		elastic.SetHealthcheckInterval(time.Duration(e.config.ResponseHeaderTimeout)*time.Second),
		elastic.SetHealthcheckTimeout(time.Duration(e.config.ResponseHeaderTimeout)*time.Second),
	)
	if err != nil {
		return err
	}
	e.client = client
	return nil
}

func (e *ElasticsearchClient) GetClient() *elastic.Client {
	return e.client
}

func (e *ElasticsearchClient) Ping() error {
	for _, hostUrl := range e.config.Host {
		info, code, err := e.client.Ping(hostUrl).Do(context.Background())
		if err != nil {
			return err
		} else if code != 200 {
			return errors.New("elasticsearch ping fail")
		}
		log.Infof("Elasticsearch returned with code %d and version %s", code, info.Version.Number)
	}
	return nil
}
