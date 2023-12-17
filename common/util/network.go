package util

import "github.com/tel4vn/fins-microservices/common/constants"

// Map network
func MapNetworkPlugin(externalTelcoId string) string {
	telcoId := externalTelcoId
	if telco, exist := constants.MapNetwork[externalTelcoId]; exist {
		telcoId = telco
	} else {
		telcoId = externalTelcoId
	}

	return telcoId
}
