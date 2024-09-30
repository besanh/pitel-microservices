package model

type NotifyPayload struct {
	Detail   map[string]string `json:"detail"`
	DeviceId string            `json:"device_id"` // user@tenant
	Message  string            `json:"message"`
	Title    string            `json:"title"`
	Type     string            `json:"type"` // notify, warning
}
