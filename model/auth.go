package model

type AuthUser struct {
	TenantId           string   `json:"tenant_id"`
	BusinessUnitId     string   `json:"business_unit_id"`
	UserId             string   `json:"user_id"`
	Username           string   `json:"username"`
	Services           []string `json:"services"`
	DatabaseName       string   `json:"database_name"`
	DatabasePort       int      `json:"database_port"`
	DatabaseHost       string   `json:"database_host"`
	DatabaseUser       string   `json:"database_user"`
	DatabasePassword   string   `json:"database_password"`
	DatabaseEsHost     string   `json:"database_es_host"`
	DatabaseEsUser     string   `js:"database_es_user"`
	DatabaseEsPassword string   `json:"database_es_password"`
	DatabaseEsIndex    string   `json:"database_es_index"`
	RoutingConfigUuid  string   `json:"routing_config_uuid"`
}
