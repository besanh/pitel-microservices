package model

import (
	"errors"
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
	if len(r.ConnectionQueueId) > 0 {
		return nil
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
