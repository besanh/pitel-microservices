package model

import "database/sql"

type AuthSourceFilter struct {
	Source string
	Status sql.NullBool
}
