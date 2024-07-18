package model

import (
	"errors"
	"regexp"
	"slices"

	"github.com/tel4vn/fins-microservices/common/regex"
	"github.com/tel4vn/fins-microservices/common/variables"
)

type AttachConnectionQueueToConnectionAppRequest struct {
	ConnectionAppRequest ChatConnectionAppRequest           `json:"connection_app_request"`
	ConnectionQueueId    string                             `json:"connection_queue_id"` // for selecting an existed queue
	ChatQueue            PipelineChatQueueRequest           `json:"chat_queue"`
	ChatQueueUser        PipelineChatQueueUserRequest       `json:"chat_queue_user"`
	ChatManageQueueUser  PipelineChatManageQueueUserRequest `json:"chat_manage_queue_user"`
}

type PipelineChatManageQueueUserRequest struct {
	UserId string `json:"user_id"`
	IsNew  bool   `json:"is_new"`
}

type PipelineChatQueueRequest struct {
	QueueName     string `json:"queue_name"`
	Description   string `json:"description"`
	ChatRoutingId string `json:"chat_routing_id"`
	Status        string `json:"status"`
}

type PipelineChatQueueUserRequest struct {
	UserId []string `json:"user_id"`
	Source string   `json:"source"`
}

func (r *AttachConnectionQueueToConnectionAppRequest) Validate() error {
	if len(r.ConnectionAppRequest.ConnectionName) < 1 {
		return errors.New("connection name is required")
	}
	if len(r.ConnectionAppRequest.ConnectionType) < 1 {
		return errors.New("connection type is required")
	}

	if !slices.Contains[[]string](variables.CONNECTION_TYPE, r.ConnectionAppRequest.ConnectionType) {
		return errors.New("connection type " + r.ConnectionAppRequest.ConnectionType + " is not supported")
	}

	if r.ConnectionAppRequest.ConnectionType == "zalo" {
		if len(r.ConnectionAppRequest.OaInfo.Zalo) < 1 {
			return errors.New("oa info zalo is required for zalo connection type")
		}

		re := regexp.MustCompile(regex.REGEX_URL)
		if !re.MatchString(r.ConnectionAppRequest.OaInfo.Zalo[0].UrlOa) {
			return errors.New("url " + r.ConnectionAppRequest.OaInfo.Zalo[0].UrlOa + " is not valid")
		}
	}

	if r.ConnectionAppRequest.ConnectionType == "facebook" {
		if len(r.ConnectionAppRequest.OaInfo.Facebook) < 1 {
			return errors.New("oa info facebook is required for facebook connection type")
		}
	}

	if len(r.ConnectionAppRequest.Status) < 1 {
		return errors.New("status is required")
	}

	if len(r.ConnectionAppRequest.ChatAppId) < 1 {
		return errors.New("chat app id is required")
	}
	if len(r.ChatQueue.QueueName) < 1 {
		return errors.New("chat queue name is required")
	}
	if len(r.ChatQueue.ChatRoutingId) < 1 {
		return errors.New("chat queue routing id is required")
	}
	if len(r.ChatQueueUser.UserId) < 1 {
		return errors.New("chat queue user id is required")
	}
	if len(r.ChatManageQueueUser.UserId) < 1 {
		return errors.New("chat queue user id is required")
	}

	return nil
}
