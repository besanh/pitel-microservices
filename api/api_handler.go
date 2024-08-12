package apiHandler

import (
	"github.com/cardinalby/hureg"
	apiv1 "github.com/tel4vn/fins-microservices/api/v1"
)

func InitAPI(api hureg.APIGen) {
	apiv1.RegisterAPIIBKAuth(api)
	apiv1.RegisterAPIIBKUser(api)
	apiv1.RegisterAPIIBKScope(api)
	apiv1.RegisterAPIIBKToken(api)
	apiv1.RegisterAPIIBKBusinessUnit(api)
}
