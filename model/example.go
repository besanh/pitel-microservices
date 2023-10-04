package model

import "github.com/uptrace/bun"

type Example struct {
	*Base
	bun.BaseModel `bun:"table:example"`
	ExampleName   string `bun:"example_name,type:text,notnull"`
}
