package model

import "database/sql"

type AuthSourceFilter struct {
	Source string
	Status sql.NullBool
}

type AppFilter struct {
	AppName string
	Status  sql.NullBool
}

type ConnectionAppFilter struct {
	ConnectionName string
	ConnectionType string
	Status         sql.NullBool
}

type QueueFilter struct {
	AppId     string
	QueueName string
}

type ChatQueueAgentFilter struct {
	QueueId string
	AgentId string
	Source  string
}

type ChatRoutingFilter struct {
	RoutingName string
	Status      sql.NullBool
}

type ConversationFilter struct {
	AppId       []string
	UserIdByApp []string
	Username    []string
	PhoneNumber []string
	Email       []string
}
