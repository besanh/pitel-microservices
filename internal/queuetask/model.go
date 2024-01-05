package queuetask

import "github.com/tel4vn/fins-microservices/model"

type MessageDeliveryPayload struct {
	OttMessage *model.OttMessage `json:"ott_message"`
}
