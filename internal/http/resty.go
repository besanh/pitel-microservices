package http

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/model"
)

type IResty[T model.Resty] interface {
	Get(ctx context.Context, setting model.RestySetting, entry map[string]string) (result *T, err error)
}

type Resty[T model.Resty] struct {
}

func NewResty[T model.Resty]() IResty[T] {
	return &Resty[T]{}
}

func (r *Resty[T]) GetClient(client *resty.Client, setting model.RestySetting) {
	if setting.InsecureSkipVerify.Valid {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: setting.InsecureSkipVerify.Bool})
	}
	if setting.Timeout > 0 {
		client.SetTimeout(setting.Timeout)
	}
	if setting.AuthType == "Basic" {
		client.SetBasicAuth(setting.RestyAuth.Username, setting.RestyAuth.Password)
	} else if setting.AuthType == "Bearer" {
		client.SetAuthToken(setting.RestyAuth.Token)
	}
}

func (r *Resty[T]) Get(ctx context.Context, setting model.RestySetting, entry map[string]string) (result *T, err error) {
	client := resty.New()
	r.GetClient(client, setting)

	resp, err := client.R().
		SetQueryParams(entry).
		SetHeader("Accept", setting.Accept).
		Get(setting.Url)

	if err != nil {
		return
	} else if resp.IsError() {
		err = errors.New(resp.String())
		return
	}

	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return
	}
	return
}

func (r *Resty[T]) Post(ctx context.Context, setting model.RestySetting, entry map[string]string) (result *T, err error) {
	client := resty.New()
	r.GetClient(client, setting)

	resp, err := client.R().
		SetBody(entry).
		SetHeader("Accept", setting.Accept).
		Post(setting.Url)

	if err != nil {
		return
	} else if resp.IsError() {
		err = errors.New(resp.String())
		return
	}

	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return
	}
	return
}

func (r *Resty[T]) Put(ctx context.Context, setting model.RestySetting, entry map[string]string) (result *T, err error) {
	client := resty.New()
	r.GetClient(client, setting)

	resp, err := client.R().
		SetBody(entry).
		SetHeader("Accept", setting.Accept).
		Put(setting.Url)

	if err != nil {
		return
	} else if resp.IsError() {
		err = errors.New(resp.String())
		return
	}

	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return
	}
	return
}

func (r *Resty[T]) Patch(ctx context.Context, setting model.RestySetting, entry map[string]string) (result *T, err error) {
	client := resty.New()
	r.GetClient(client, setting)

	resp, err := client.R().
		SetBody(entry).
		SetHeader("Accept", setting.Accept).
		Patch(setting.Url)

	if err != nil {
		return
	} else if resp.IsError() {
		err = errors.New(resp.String())
		return
	}

	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return
	}
	return
}

func (r *Resty[T]) Delete(ctx context.Context, setting model.RestySetting, entry map[string]string) (result *T, err error) {
	client := resty.New()
	r.GetClient(client, setting)

	resp, err := client.R().
		SetQueryParams(entry).
		SetHeader("Accept", setting.Accept).
		Delete(setting.Url)

	if err != nil {
		return
	} else if resp.IsError() {
		err = errors.New(resp.String())
		return
	}

	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return
	}
	return
}
