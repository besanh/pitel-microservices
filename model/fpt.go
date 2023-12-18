package model

type FptGetTokenResponseSuccess struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type FptGetTokenResponseError struct {
	Err              int    `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// CSKH CS
type FptRequestCS struct {
	AccessToken string `json:"access_token"`
	SessionId   string `json:"session_id"`
	BrandName   string `json:"BrandName"`
	Phone       string `json:"Phone"`
	Message     string `json:"Message"`
	RequestId   string `json:"RequestId"`
}

type FptSendMessageResponse struct {
	MessageId string `json:"MessageId"`
	Phone     string `json:"Phone"`
	BrandName string `json:"BrandName"`
	Message   string `json:"Message"`
	PartnerId string `json:"PartnerId"`
	Telco     string `json:"Telco"`
}

type FptCampaignRequest struct {
	AccessToken  string `json:"access_token"`
	SessionId    string `json:"session_id"`
	CampaignName string `json:"CampaignName"`
	BrandName    string `json:"BrandName"`
	Message      string `json:"Message"`
	ScheduleTime string `json:"ScheduleTime"`
	Quota        string `json:"Quota"`
}

type FptCampaignResponse struct {
	CampaignCode string `json:"campaign_code"`
}

type FptSendMessageQCRequest struct {
	AccessToken  string `json:"access_token"`
	SessionId    string `json:"session_id"`
	CampaignCode string `json:"CampaignCode"`
	PhoneList    string `json:"PhoneList"` // delimiter by ","
}

type FptSendMessageQCResponse struct {
	NumMessageSent int    `json:"NumMessageSent"`
	NumRemainQuota int    `json:"NumRemainQuota"`
	BatchId        string `json:"BatchId"`
}
