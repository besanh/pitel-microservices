package model

type OttMessage struct {
	TenantId       string            `json:"tenant_id"`
	MessageType    string            `json:"message_type"`
	EventName      string            `json:"event_name"`
	AppId          string            `json:"app_id"`
	AppName        string            `json:"app_name"`
	OaId           string            `json:"oa_id"`
	ShareInfo      *ShareInfo        `json:"share_info"`
	UserIdByApp    string            `json:"user_id_by_app"`
	ExternalUserId string            `json:"external_user_id"`
	Username       string            `json:"username"`
	Avatar         string            `json:"avatar"`
	Timestamp      int64             `json:"timestamp"`
	MsgId          string            `json:"msg_id"`
	Content        string            `json:"content"`
	QuoteMsgId     string            `json:"quote_msg_id"`
	ReactIcon      string            `json:"react_icon"`
	Attachments    *[]OttAttachments `json:"attachments"`
}

type OttAttachments struct {
	Payload *OttPayloadMedia `json:"payload"`
	AttType string           `json:"att_type"` // image/audio/video/link/sticker/gif/file
}

type OttAttachmentsPayload struct {
	Url string `json:"url"`
}

type OttPayloadMedia struct {
	Thubnail    string `json:"thumbnail"`
	Description string `json:"description"`
	Url         string `json:"url"`
	UrlOttBk    string `json:"url_ott_bk"`
	Id          string `json:"id"`
	Title       string `json:"title"`
}

type OttPayloadFile struct {
	Size     string `json:"size"`
	Url      string `json:"url"`
	Name     string `json:"name"`
	Checksum string `json:"checksum"`
	Type     string `json:"type"`
}

type SendMessageToOtt struct {
	Type          string `json:"type"`
	EventName     string `json:"event_name"`
	AppId         string `json:"app_id"`
	OaId          string `json:"oa_id"`
	UserIdByApp   string `json:"user_id_by_app"`
	Uid           string `json:"uid"`
	SupporterId   string `json:"supporter_id"`
	SupporterName string `json:"supporter_name"`
	Timestamp     string `json:"timestamp"`
	MsgId         string `json:"msg_id"`
	Text          string `json:"text"`
}

type SendMessageToOttWithAttachment struct {
	Type          string            `json:"type"`
	EventName     string            `json:"event_name"`
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

type OttResponse struct {
	Code string `json:"code"`
	Data struct {
		MsgId string `json:"msg_id"`
		Uid   string `json:"uid"`
	} `json:"data"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type OttAccount struct {
	Type        string `json:"type"`
	AppId       string `json:"app_id"`
	AppName     string `json:"app_name"`
	AppSecret   string `json:"app_secret"`
	OaId        string `json:"oa_id"`
	OaName      string `json:"oa_name"`
	AccessToken string `json:"access_token"`
	Status      string `json:"status"`
}

type OttCodeChallenge struct {
	Data string `json:"data"`
}

type OttShareInfo struct {
	AppId       string `json:"app_id"`
	OaId        string `json:"oa_id"`
	Uid         string `json:"uid"`
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	Address     string `json:"address"`
	City        string `json:"city"`
	District    string `json:"district"`
	PhoneNumber string `json:"phone_number"`
}

type OttUploadResponse struct {
	Data []string `json:"data"`
}
