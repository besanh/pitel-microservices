package model

type AuthUser struct {
	TenantId         string   `json:"tenant_id"`
	BusinessUnitId   string   `json:"business_unit_id"`
	UserId           string   `json:"user_id"`
	Username         string   `json:"username"`
	Services         []string `json:"services"`
	Level            string   `json:"level"`
	DatabaseName     string   `json:"database_name"`
	DatabasePort     int      `json:"database_port"`
	DatabaseHost     string   `json:"database_host"`
	DatabaseUser     string   `json:"database_user"`
	DatabasePassword string   `json:"database_password"`
	Source           string   `json:"source"`
	Token            string   `json:"token"`
	UnitUuid         string   `json:"unit_uuid"` //only for user crm
}

type AAAResponse struct {
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Data    *AuthUser `json:"data"`
}

type BssAuthRequest struct {
	Token   string `json:"token"`
	AuthUrl string `json:"auth_url"`
	Source  string `json:"source"`
	UserId  string `json:"user_id"`
}
