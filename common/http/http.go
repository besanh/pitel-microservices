package http

import (
	"crypto/tls"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/log"
)

func PostHeader(url string, header string, body any) (*resty.Response, error) {
	log.Infof("PostHeader | url : %v | headers : %v | body : %v", url, header, body)
	client := resty.New()
	client.SetTimeout(time.Second * 3)
	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	return client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", header).
		SetBody(body).
		Post(url)
}

func Post(url string, body any) (*resty.Response, error) {
	client := resty.New()
	client.SetTimeout(time.Second * 10)
	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	return client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
}
