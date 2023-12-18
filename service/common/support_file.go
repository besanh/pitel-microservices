package common

import (
	"strings"

	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/internal/redis"
)

func SetExportValue(domainUuid string, exportName string, exportMap []string) error {
	exportValue := strings.Join(exportMap, ";")
	dataTmp := []interface{}{exportName, exportValue}
	_, err := redis.Redis.HSet(constants.EXPORT_KEY+domainUuid, dataTmp)
	if err != nil {
		return err
	}
	return nil
}
