package model

type InboxMarketingRequest struct {
	RoutingConfig string   `json:"routing_config"`
	Content       string   `json:"content"`
	PhoneNumber   string   `json:"phone_number"`
	Template      string   `json:"template"` // template uuid
	Channel       []string `json:"channel"`
}
