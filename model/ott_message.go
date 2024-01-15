package model

type OttMessage struct {
	MessageType string            `json:"message_type"`
	EventName   string            `json:"event_name"`
	AppId       string            `json:"app_id"`
	AppName     string            `json:"app_name"`
	OaId        string            `json:"oa_id"`
	UserIdByApp string            `json:"user_id_by_app"`
	UserId      string            `json:"user_id"`
	Username    string            `json:"username"`
	Avatar      string            `json:"avatar"`
	Timestamp   int64             `json:"timestamp"`
	MsgId       string            `json:"msg_id"`
	Content     string            `json:"content"`
	Attachments *[]OttAttachments `json:"attachments"`
}

type OttAttachments struct {
	Payload any    `json:"payload"`
	AttType string `json:"att_type"` // image/audio/video/link/sticker/gif/file
}

type OttPayloadMedia struct {
	Thubnail    string `json:"thumbnail"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Id          string `json:"id"` // only sticker
}

type OttPayloadFile struct {
	Size     string `json:"size"`
	Url      string `json:"url"`
	Name     string `json:"name"`
	Checksum string `json:"checksum"`
	Type     string `json:"type"`
}

type SendMessageToOtt struct {
	Type          string            `json:"type"`
	EventName     string            `json:"event_name" bun:"event_name,type:text,notnull"`
	AppId         string            `json:"app_id"`
	OaId          string            `json:"oa_id"`
	UserIdByApp   string            `json:"user_id_by_app"`
	Uid           string            `json:"uid"`
	SupporterId   string            `json:"supporter_id"`
	SupporterName string            `json:"supporter_name"`
	Timestamp     string            `json:"timestamp"`
	MsgId         string            `json:"msg_id"`
	Text          string            `json:"text"`
	Attachments   []*OttAttachments `json:"attachments"`
}

type MessageOttCache struct {
	EventType   string `json:"event_type"`
	AppId       string `json:"app_id"`
	OaId        string `json:"oa_id"`
	UserIdByApp string `json:"user_id_by_app"`
	UserId      string `json:"user_id"`
}

type OttResponse struct {
	Code int `json:"code"`
	Data struct {
		MsgId   string `json:"msg_id"`
		Uid     string `json:"uid"`
		Message string `json:"message"`
	}
}
