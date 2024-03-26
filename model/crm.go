package model

type UserCrmFilter struct {
	UnitUuid string `json:"unit_uuid"`
	Level    string `json:"level"`
	Username string `json:"username"`
	UserId   string `json:"user_id"`
}
