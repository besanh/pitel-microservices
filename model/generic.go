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

	GetEntityId() any
}

type Base struct {
	Id        string    `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `bun:"created_at,notnull" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull" json:"updated_at"`
	DeletedAt time.Time `bun:"deleted_at,nullzero" json:"-"`
	CreatedBy string    `bun:"created_by,type:uuid,nullzero" json:"created_by"`
	UpdatedBy string    `bun:"updated_by,type:uuid,nullzero" json:"updated_by"`
	DeletedBy string    `bun:"deleted_by,type:uuid,nullzero" json:"-"`
	IsDeleted bool      `bun:"is_deleted,type:bool,notnull" json:"-"`
}

type (
	GenericFilter struct {
		Field    string          `json:"field,omitempty" required:"false" pattern:"^[a-zA-Z0-9 _-]{0,200}$"`
		Operator string          `json:"operator,omitempty" required:"false" pattern:"^[=<>]|IN$"`
		Value    any             `json:"value,omitempty" required:"false" pattern:"^[a-zA-Z0-9 _-]{0,200}$"`
		And      []GenericFilter `json:"and,omitempty"`
		Or       []GenericFilter `json:"or,omitempty"`
	}
)

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

func (b *Base) GetEntityId() any {
	return b.Id
}

func InitBaseModel(createdBy string, updatedBy string) *Base {
	return &Base{
		Id:        uuid.NewString(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Time{},
		CreatedBy: createdBy,
		UpdatedBy: updatedBy,
		IsDeleted: false,
	}
}
