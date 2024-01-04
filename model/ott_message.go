package model

type GetOttMessage struct {
	Type        string            `json:"type" bun:"type,type:text,notnull"`
	EventName   string            `json:"event_name" bun:"event_name,type:text,notnull"` // text/image/audio/video/file/link/sticker/gif
	AppId       string            `json:"app_id" bun:"app_id,type:text,notnull"`
	AppName     string            `json:"app_name" bun:"app_name,type:text"`
	OaId        string            `json:"oa_id" bun:"oa_id,type:text"`
	UserIdByApp string            `json:"user_id_by_app" bun:"user_id_by_app,type:text"`
	UserId      string            `json:"user_id" bun:"user_id,type:text,notnull"`   // uid, id user zalo
	Username    string            `json:"username" bun:"username,type:text,notnull"` // username zalo
	Avatar      string            `json:"avatar" bun:"avatar,type:text"`
	Timestamp   int64             `json:"timestamp" bun:"timestamp,type:numeric,notnull"`
	MsgId       string            `json:"msg_id" bun:"msg_id,type:text,notnull"`
	Text        string            `json:"text" bun:"text,type:text"`
	Attachments []*OttAttachments `json:"attachments"`
}

type OttAttachments struct {
	Payload any    `json:"payload" bun:"payload,type:text,notnull"`
	AttType string `json:"att_type" bun:"att_type,type:text,notnull"` // image/audio/video/link/sticker/gif/file
}

type OttPayloadMedia struct {
	Thubnail    string `json:"thumbnail"`
	Description string `json:"description"`
	Url         string `json:"url" bun:"url,type:text,notnull"`
	Id          string `json:"id"` // only sticker
}

type OttPayloadFile struct {
	Size     string `json:"size" bun:"size,type:text,notnull"`
	Url      string `json:"url" bun:"url,type:text,notnull"`
	Name     string `json:"name" bun:"name,type:text,notnull"`
	Checksum string `json:"checksum" bun:"checksum,type:text,notnull"`
	Type     string `json:"type" bun:"type,type:text,notnull"`
}

type SendMessageToOtt struct {
	Type          string            `json:"type" bun:"type,type:text,notnull"`
	EventName     string            `json:"event_name" bun:"event_name,type:text,notnull"`
	AppId         string            `json:"app_id" bun:"app_id,type:text,notnull"`
	OaId          string            `json:"oa_id" bun:"oa_id,type:text"`
	UserIdByApp   string            `json:"user_id_by_app" bun:"user_id_by_app,type:text"`
	Uid           string            `json:"uid" bun:"uid,type:text,notnull"`
	SupporterId   string            `json:"supporter_id" bun:"supporter_id,type:text,notnull"`
	SupporterName string            `json:"supporter_name" bun:"supporter_name,type:text"`
	Timestamp     string            `json:"timestamp" bun:"timestamp,type:text,notnull"`
	MsgId         string            `json:"msg_id" bun:"msg_id,type:text,notnull"`
	Text          string            `json:"text" bun:"text,type:text"`
	Attachments   []*OttAttachments `json:"attachments"`
}
