package model

// Status: 1 success, 2 hoac -11 dang cho chua co kq, 0 fail
type FptWebhook struct {
	SmsId   int    `json:"smsid"`
	Status  int    `json:"Status"`
	Telco   string `json:"Telco"` // viettel, vina, mobi, htc, beeline, itel
	Error   string `json:"Error"`
	MtCount int    `json:"mt_count"`
}
