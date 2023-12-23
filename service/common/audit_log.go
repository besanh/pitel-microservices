package common

import (
	"encoding/json"
	"time"

	"github.com/tel4vn/fins-microservices/model"
)

func HandleAuditLogInboxMarketing(auditLog model.LogInboxMarketing) (string, error) {
	logInfo := map[string]model.LogInboxMarketing{}
	t := time.Now().Format(time.RFC3339)
	logInfo[t] = auditLog
	result, err := json.Marshal(logInfo)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
