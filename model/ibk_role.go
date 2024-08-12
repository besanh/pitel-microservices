package model

import (
	"github.com/uptrace/bun"
)

type (
	IBKRole struct {
		*Base
		bun.BaseModel `bun:"table:ibk_roles"`
		TenantId      string        `bun:"tenant_id,type:uuid,notnull"`
		RoleName      string        `bun:"role_name,type:text,notnull"`
		Description   string        `bun:"description,type:text"`
		Permission    IBKPermission `bun:"permission"`
	}
	IBKRoleInfo struct {
		*IBKRole
		bun.BaseModel `bun:"table:ibk_roles,alias:role"`
	}
)
type (
	IBKRoleQueryParam struct {
		RoleName_Eq string `query:"role_name"`
	}
	IBKRoleBody struct {
		RoleName    string        `json:"role_name" required:"true" pattern:"^[a-zA-Z0-9_ ]+$"`
		Permission  IBKPermission `json:"permission" required:"true"`
		Description string        `json:"description"`
	}
)

type (
	IBKPermission struct {
		PermissionMain           PermissionDetail `json:"permission_main"`
		PermissionMonitor        PermissionDetail `json:"permission_monitor"`
		PermissionRole           PermissionDetail `json:"permission_role"`
		PermissionUser           PermissionDetail `json:"permission_user"`
		PermissionBU             PermissionDetail `json:"permission_bu"`
		PermissionConnection     PermissionDetail `json:"permission_connection"`
		PermissionOaName         PermissionDetail `json:"permission_oa_name"`
		PermissionSampleTemplate PermissionDetail `json:"permission_sample_template"`
		PermissionLead           PermissionDetail `json:"permission_lead"`
		PermissionSendMessage    PermissionDetail `json:"permission_send_message"`
		PermissionCampaign       PermissionDetail `json:"permission_campaign"`
		PermissionSendingHistory PermissionDetail `json:"permission_sending_history"`
		PermissionReport         PermissionDetail `json:"permission_report"`
	}

	PermissionDetail struct {
		All struct {
			Enable bool `json:"enable"`
			Value  bool `json:"value"`
		} `json:"all"`
		View struct {
			Enable bool `json:"enable"`
			Value  bool `json:"value"`
		} `json:"view"`
		Search struct {
			Enable bool `json:"enable"`
			Value  bool `json:"value"`
		}
		Create struct {
			Enable bool `json:"enable"`
			Value  bool `json:"value"`
		}
		Edit struct {
			Enable bool `json:"enable"`
			Value  bool `json:"value"`
		}
		Delete struct {
			Enable bool `json:"enable"`
			Value  bool `json:"value"`
		}
		Import struct {
			Enable bool `json:"enable"`
			Value  bool `json:"value"`
		}
		Export struct {
			Enable bool `json:"enable"`
			Value  bool `json:"value"`
		}
	}
)
