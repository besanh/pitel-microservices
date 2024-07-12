package model

import (
	"errors"
	"reflect"

	"github.com/uptrace/bun"
)

type ChatRole struct {
	*Base
	bun.BaseModel `bun:"table:chat_role,alias:cr"`
	RoleName      string           `json:"role_name" bun:"role_name,type:text,notnull"`
	Status        bool             `json:"status" bun:"status,type:boolean,nullzero,default:false"`
	Setting       *ChatRoleSetting `json:"setting" bun:"setting,type:jsonb,notnull"`
}

type ChatRoleSetting struct {
	AssignConversation   bool `json:"assign_conversation"`
	ReassignConversation bool `json:"reassign_conversation"`
	MakeDone             bool `json:"make_done"`
	AddLabel             bool `json:"add_label"`
	RemoveLabel          bool `json:"remove_label"`
	Major                bool `json:"major"`
	Following            bool `json:"following"`
	SubmitForm           bool `json:"submit_form"`
}

type ChatRoleRequest struct {
	RoleName string          `json:"role_name"`
	Status   bool            `json:"status"`
	Setting  ChatRoleSetting `json:"setting"`
}

func (m *ChatRoleRequest) Validate() error {
	if len(m.RoleName) < 1 {
		return errors.New("role name is required")
	}
	if reflect.DeepEqual(m.Setting, "{}") {
		return errors.New("setting is required")
	}
	return nil
}
