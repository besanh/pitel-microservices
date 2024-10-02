package util

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tel4vn/fins-microservices/common/log"
)

func ConvertToBytes(message any) ([]byte, error) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Error(err)
	}
	return messageBytes, err
}

func ConvertMillisToTimeString(millis int) string {
	duration := time.Duration(millis) * time.Millisecond
	t := time.Time{}.Add(duration)

	// Format the string as hh:mm:ss.milliseconds
	// Use .Format to format up to seconds, and manually append milliseconds
	return fmt.Sprintf("%s.%03d", t.Format("15:04:05"), millis%1000)
}
