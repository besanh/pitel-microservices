package model

type Event struct {
	EventName string     `json:"event_name"`
	EventData *EventData `json:"event_data"`
}

type EventData struct {
	Message      any `json:"message"`
	Conversation any `json:"conversation"`
	ShareInfo    any `json:"share_info"`
}

type WsEvent struct {
	EventName string         `json:"event_name"`
	EventData map[string]any `json:"event_data"`
}
