package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type Profile struct {
	profile service.IProfile
}

func NewProfile(engine *gin.Engine, profileService service.IProfile) {
	handler := &Profile{
		profile: profileService,
	}

	Group := engine.Group("bss-message/v1/profile")
	{
		Group.GET("get-update", handler.GetUpdateProfile)
	}
}

func (handler *Profile) GetUpdateProfile(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	request := model.ProfileRequest{
		AppId:          c.Query("app_id"),
		OaId:           c.Query("oa_id"),
		UserId:         c.Query("user_id"),
		ProfileType:    c.Query("profile_type"),
		ConversationId: c.Query("conversation_id"),
	}

	if err := request.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	code, result := handler.profile.GetUpdateProfileByUserId(c, res.Data, request)
	c.JSON(code, result)
}
