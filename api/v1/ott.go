package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"golang.org/x/exp/slices"
)

type OttMessage struct {
	ottMessageService    service.IOttMessage
	connectionAppService service.IChatConnectionApp
	conversationService  service.IConversation
}

func NewOttMessage(r *gin.Engine, messageService service.IOttMessage, connectionApp service.IChatConnectionApp, conversation service.IConversation) {
	handler := &OttMessage{
		ottMessageService:    messageService,
		connectionAppService: connectionApp,
		conversationService:  conversation,
	}

	Group := r.Group("bss-message/v1/ott")
	{
		Group.POST("", handler.GetOttMessage)
		Group.GET("code-challenge/:app_id", handler.GetCodeChallenge)
		Group.POST("ask-info", handler.AskInfo)
	}
}

func (h *OttMessage) GetOttMessage(c *gin.Context) {
	jsonBody := make(map[string]any, 0)
	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	log.Info("ott get message body: ", jsonBody)
	messageType, _ := jsonBody["type"].(string)
	eventName, _ := jsonBody["event_name"].(string)
	appId, _ := jsonBody["app_id"].(string)
	appName, _ := jsonBody["app_name"].(string)
	oaId, _ := jsonBody["oa_id"].(string)
	userIdByApp, _ := jsonBody["user_id_by_app"].(string)
	externalUserId, _ := jsonBody["uid"].(string)
	username, _ := jsonBody["username"].(string)
	avatar, _ := jsonBody["avatar"].(string)
	timestampTmp, _ := jsonBody["timestamp"].(string)
	timestamp, _ := strconv.ParseInt(timestampTmp, 10, 64)
	msgId, _ := jsonBody["msg_id"].(string)
	content, _ := jsonBody["text"].(string)
	attachmentsTmp, _ := jsonBody["attachments"].([]any)
	attachmentsAny := make([]any, 0)
	for item := range attachmentsTmp {
		tmp := attachmentsTmp[item].(map[string]any)
		attType, _ := tmp["att_type"].(string)
		if slices.Contains[[]string](variables.EVENT_NAME_SEND_MESSAGE, attType) {
			attachment := map[string]any{
				"att_type": attType,
				"payload":  tmp["payload"],
			}
			attachmentsAny = append(attachmentsAny, attachment)
		}
	}
	attachments := make([]model.OttAttachments, 0)
	if err := util.ParseAnyToAny(attachmentsAny, &attachments); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	shareInfoTmp, _ := jsonBody["share_info"].(map[string]any)
	shareInfo := model.ShareInfo{}
	if eventName != variables.EVENT_NAME_EXCLUDE["oa_connection"] && shareInfoTmp != nil {
		shareInfoName, _ := shareInfoTmp["name"].(string)
		shareInfoPhoneNumber, _ := shareInfoTmp["phone"].(string)
		shareInfoAddress, _ := shareInfoTmp["address"].(string)
		shareInfoCity, _ := shareInfoTmp["city"].(string)
		shareInfoDistrict, _ := shareInfoTmp["district"].(string)
		shareInfo = model.ShareInfo{
			Fullname:    shareInfoName,
			PhoneNumber: shareInfoPhoneNumber,
			Address:     shareInfoAddress,
			City:        shareInfoCity,
			District:    shareInfoDistrict,
		}
	}

	var message model.OttMessage
	if eventName == variables.EVENT_NAME_EXCLUDE["oa_connection"] {
		oaInfoMessageTmp, _ := jsonBody["oa_info"].(map[string]any)
		oaInfoMessageCode, _ := oaInfoMessageTmp["code"].(float64)
		oaInfoMessage := model.OaInfoMessage{}
		if oaInfoMessageCode == 200 {
			if err := util.ParseAnyToAny(oaInfoMessageTmp, &oaInfoMessage); err != nil {
				c.JSON(response.BadRequestMsg(err))
				return
			}
		} else if oaInfoMessageCode != 200 {
			log.Error("get oa info error: ", oaInfoMessageTmp)
			c.JSON(response.BadRequestMsg("get oa info error"))
			return
		}
		if len(oaInfoMessage.ConnectionId) < 1 {
			c.JSON(response.BadRequestMsg("connection_id is required"))
			return
		}
		connectionAppRequest := model.ChatConnectionAppRequest{
			OaId:     oaId,
			AppId:    appId,
			Id:       oaInfoMessage.ConnectionId,
			OaName:   oaInfoMessage.Name,
			Avatar:   oaInfoMessage.Avatar,
			Cover:    oaInfoMessage.Cover,
			CateName: oaInfoMessage.CateName,
			Status:   "active",
		}
		authUser := model.AuthUser{
			Source: "authen",
		}
		if err := h.connectionAppService.UpdateChatConnectionAppById(c, &authUser, oaInfoMessage.ConnectionId, connectionAppRequest); err != nil {
			c.JSON(response.BadRequestMsg(err))
			return
		}
		c.JSON(response.OKResponse())
	} else if eventName == "submit_info" {
		code, result := h.conversationService.UpdateConversationById(c, &model.AuthUser{}, appId, externalUserId, shareInfo)
		c.JSON(code, result)
		return
	} else {
		message = model.OttMessage{
			MessageType:    messageType,
			EventName:      eventName,
			AppId:          appId,
			AppName:        appName,
			OaId:           oaId,
			ShareInfo:      &shareInfo,
			UserIdByApp:    userIdByApp,
			ExternalUserId: externalUserId,
			Username:       username,
			Avatar:         avatar,
			Timestamp:      timestamp,
			MsgId:          msgId,
			Content:        content,
			Attachments:    &attachments,
		}
		code, result := h.ottMessageService.GetOttMessage(c, message)
		c.JSON(code, result)
		return
	}
}

func (h *OttMessage) GetCodeChallenge(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	appId := c.Param("app_id")
	if len(appId) < 1 {
		c.JSON(response.BadRequestMsg("app_id is required"))
		return
	}

	code, result := h.ottMessageService.GetCodeChallenge(c, res.Data, appId)
	c.JSON(code, result)
}

func (h *OttMessage) AskInfo(c *gin.Context) {
	jsonBody := make(map[string]any, 0)
	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	log.Info("ott get ask info body: ", jsonBody)

	appId, _ := jsonBody["app_id"].(string)
	externalUserId, _ := jsonBody["uid"].(string)
	shareInfoTmp, _ := jsonBody["share_info"].(map[string]any)
	shareInfoName, _ := shareInfoTmp["name"].(string)
	shareInfoPhoneNumber, _ := shareInfoTmp["phone"].(string)
	shareInfoAddress, _ := shareInfoTmp["address"].(string)
	shareInfoCity, _ := shareInfoTmp["city"].(string)
	shareInfoDistrict, _ := shareInfoTmp["district"].(string)
	shareInfo := model.ShareInfo{
		Fullname:    shareInfoName,
		PhoneNumber: shareInfoPhoneNumber,
		Address:     shareInfoAddress,
		City:        shareInfoCity,
		District:    shareInfoDistrict,
	}

	code, result := h.conversationService.UpdateConversationById(c, &model.AuthUser{}, appId, externalUserId, shareInfo)
	c.JSON(code, result)
}
