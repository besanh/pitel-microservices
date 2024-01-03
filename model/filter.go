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
