package model

import (
	"errors"
	"regexp"
	"slices"

	"github.com/tel4vn/fins-microservices/common/regex"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/uptrace/bun"
)

type User struct {
	AuthUser               *AuthUser     `json:"auth_user"`
	ConnectionId           string        `json:"connection_id"`
	QueueId                string        `json:"queue_id"`
	ConnectionQueueId      string        `json:"connection_queue_id"`
	IsOk                   bool          `json:"is_ok"`
	IsReassignNew          bool          `json:"is_reassign_new"`
	IsReassignSame         bool          `json:"is_reassign_same"`
	UserIdReassignNew      string        `json:"user_id_reassign_new"`
	UserIdReassignSame     string        `json:"user_id_reassign_same"`
	UserIdRemove           string        `json:"user_id_remove"`
	PreviousAssign         *AllocateUser `json:"previous_assign"`
	ConversationId         string        `json:"conversation_id"`
	ExternalConversationId string        `json:"external_conversation_id"`
	NewAllocateUserId      string        `json:"new_allocate_user_id"`
}

type ChatUser struct {
	*Base
	bun.BaseModel `bun:"table:chat_user,alias:cu"`
	Username      string `json:"username" bun:"username,type:text,notnull"`
	Password      string `json:"password" bun:"password,type:text,notnull"`
	Email         string `json:"email" bun:"email,type:text,notnull"`
	Salt          string `json:"salt" bun:"salt,type:text,notnull"`
	Level         string `json:"level" bun:"level,type:text,notnull"`
	Avatar        string `json:"avatar" bun:"avatar,type:text,default:null"`
	Fullname      string `json:"fullname" bun:"fullname,type:text,notnull"`
	RoleId        string `json:"role_id" bun:"role_id,type:uuid,notnull"`
	Status        bool   `json:"status" bun:"status,type:boolean,nullzero,default:false"`
}

type ChatUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Level    string `json:"level" binding:"required"`
	Avatar   string `json:"avatar"`
	Fullname string `json:"fullname" binding:"required,min=3"`
	RoleId   string `json:"role_id" binding:"required"`
	Status   bool   `json:"status" binding:"required"`
}

func (m *ChatUserRequest) Validate() error {
	if len(m.Username) < 1 {
		return errors.New("username is required")
	}
	if len(m.Password) < 1 {
		return errors.New("password is required")
	}
	if len(m.Email) < 1 {
		return errors.New("email is required")
	}

	regex := regexp.MustCompile(regex.REGEX_EMAIL)
	if isOk := regex.MatchString(m.Email); !isOk {
		return errors.New("invalid email")
	}

	if len(m.Level) < 1 {
		return errors.New("level is required")
	}
	if !slices.Contains(variables.LEVELS, m.Level) {
		return errors.New("level " + m.Level + " is not supported")
	}
	if len(m.Fullname) < 1 {
		return errors.New("fullname is required")
	}
	if len(m.RoleId) < 1 {
		return errors.New("role id is required")
	}

	return nil
}
