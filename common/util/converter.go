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
	// Convert the milliseconds to a time.Duration (time.Duration expects nanoseconds)
	duration := time.Duration(millis) * time.Millisecond

	// Extract hours, minutes, seconds, and milliseconds
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	milliseconds := millis % 1000 // Get the remainder to capture only milliseconds

	// Format the string as hh:mm:ss.milliseconds
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}
