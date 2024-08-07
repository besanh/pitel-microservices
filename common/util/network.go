package util

import (
	"github.com/tel4vn/fins-microservices/common/variables"
)

// Map network
func MapNetworkPlugin(externalTelcoId string) string {
	telcoId := externalTelcoId
	if telco, exist := variables.MAP_NETWORK[externalTelcoId]; exist {
		telcoId = telco
	} else {
		telcoId = externalTelcoId
	}

	return telcoId
}
