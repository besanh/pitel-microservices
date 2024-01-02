package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/log"
)

type Model interface {
	GetId() string
	SetId(id string)
	SetUpdatedAt(t time.Time)
	SetCreatedAt(t time.Time)
}

type Base struct {
	Id        string    `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	CreatedAt time.Time `bun:"created_at,notnull"`
	UpdatedAt time.Time `bun:"updated_at,notnull"`
}

func (b *Base) GetId() string {
	return b.Id
}

func (b *Base) SetId(id string) {
	b.Id = id
}

func (b *Base) SetUpdatedAt(t time.Time) {
	b.UpdatedAt = t
}

func (b *Base) SetCreatedAt(t time.Time) {
	log.Infof("SetCreatedAt: %v", b)
	b.CreatedAt = t
}

func InitBase() *Base {
	return &Base{
		Id:        uuid.NewString(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
