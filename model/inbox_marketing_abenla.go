package model

type AbenlaCheckConnectionResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AbenlaSendMessageResponse struct {
	SmsPerMessage int    `json:"SmsPerMessage"`
	Code          int    `json:"code"`
	Message       string `json:"message"`
}
