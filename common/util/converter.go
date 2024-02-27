package util

import (
	"encoding/json"

	"github.com/tel4vn/fins-microservices/common/log"
)

func ConvertToBytes(message interface{}) ([]byte, error) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Error(err)
	}
	return messageBytes, err
}
