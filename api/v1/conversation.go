package v1

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/pitel-microservices/api"
	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/common/response"
	"github.com/tel4vn/pitel-microservices/common/util"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/service"
	"golang.org/x/exp/slices"
)

type Conversation struct {
	conversationService service.IConversation
}

func NewConversation(engine *gin.Engine, conversationService service.IConversation) {
	handler := &Conversation{
		conversationService: conversationService,
	}
	Group := engine.Group("bss-message/v1/conversation")
	{
		Group.GET("", handler.GetConversations)
		Group.GET("scroll", handler.GetConversationsWithScrollAPI)
		Group.GET("manager", handler.GetConversationsByManager)
		Group.PUT(":app_id/:oa_id/:id", handler.UpdateConversation)
		Group.POST("status", handler.UpdateStatusConversation)
		Group.GET(":app_id/:oa_id/:id", handler.GetConversationById)
		Group.PUT("label/:label_type", handler.PutLabelToConversation)
		Group.PUT("preference", handler.UpdateUserPreferenceConversation)
	}
}

func (handler *Conversation) GetConversations(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.Unauthorized())
		return
	}
	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))

	isDone := sql.NullBool{}
	if len(c.Query("is_done")) > 0 {
		isDone.Valid = true
		isDone.Bool, _ = strconv.ParseBool(c.Query("is_done"))
	}
	major := sql.NullBool{}
	if len(c.Query("major")) > 0 {
		major.Valid = true
		major.Bool, _ = strconv.ParseBool(c.Query("major"))
	}
	following := sql.NullBool{}
	if len(c.Query("following")) > 0 {
		following.Valid = true
		following.Bool, _ = strconv.ParseBool(c.Query("following"))
	}

	filter := model.ConversationFilter{
		AppId:          util.ParseQueryArray(c.QueryArray("app_id")),
		ConversationId: util.ParseQueryArray(c.QueryArray("conversation_id")),
		Username:       c.Query("username"),
		PhoneNumber:    c.Query("phone_number"),
		Email:          c.Query("email"),
		IsDone:         isDone,
		Major:          major,
		Following:      following,
	}

	total, result, err := handler.conversationService.GetConversations(c, res.Data, filter, limit, offset)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.Pagination(result, total, limit, offset))
}

func (handler *Conversation) GetConversationsWithScrollAPI(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.Unauthorized())
		return
	}
	limit := util.ParseLimit(c.Query("limit"))
	scrollId := c.Query("scroll_id")

	isDone := sql.NullBool{}
	if len(c.Query("is_done")) > 0 {
		isDone.Valid = true
		isDone.Bool, _ = strconv.ParseBool(c.Query("is_done"))
	}
	major := sql.NullBool{}
	if len(c.Query("major")) > 0 {
		major.Valid = true
		major.Bool, _ = strconv.ParseBool(c.Query("major"))
	}
	following := sql.NullBool{}
	if len(c.Query("following")) > 0 {
		following.Valid = true
		following.Bool, _ = strconv.ParseBool(c.Query("following"))
	}

	filter := model.ConversationFilter{
		AppId:          util.ParseQueryArray(c.QueryArray("app_id")),
		ConversationId: util.ParseQueryArray(c.QueryArray("conversation_id")),
		Username:       c.Query("username"),
		PhoneNumber:    c.Query("phone_number"),
		Email:          c.Query("email"),
		IsDone:         isDone,
		Major:          major,
		Following:      following,
	}

	total, data, respScrollId, err := handler.conversationService.GetConversationsWithScrollAPI(c, res.Data, filter, limit, scrollId)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	result := map[string]any{
		"conversations": data,
		"scroll_id":     respScrollId,
	}
	c.JSON(response.Pagination(result, total, limit, 0))
}

func (handler *Conversation) UpdateConversation(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.Unauthorized())
		return
	}
	appId := c.Param("app_id")
	if len(appId) < 1 {
		c.JSON(response.BadRequestMsg("app_id is required"))
	}

	oaId := c.Param("oa_id")
	if len(oaId) < 1 {
		c.JSON(response.BadRequestMsg("oa_id is required"))
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	var shareInfo model.ShareInfo
	if err := c.ShouldBindJSON(&shareInfo); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("update conversation payload -> ", shareInfo)

	code, result := handler.conversationService.UpdateConversationById(c, res.Data, appId, oaId, id, shareInfo)
	c.JSON(code, result)
}

func (handler *Conversation) UpdateStatusConversation(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.Unauthorized())
		return
	}
	jsonBody := make(map[string]any, 0)
	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	appId, _ := jsonBody["app_id"].(string)
	conversationId, _ := jsonBody["conversation_id"].(string)
	status, _ := jsonBody["status"].(string)

	log.Info("update status conversation payload -> ", jsonBody)

	if !slices.Contains([]string{"done", "reopen"}, status) {
		c.JSON(response.BadRequestMsg("status is invalid"))
		return
	}

	err := handler.conversationService.UpdateStatusConversation(c, res.Data, appId, conversationId, res.Data.UserId, status)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}

func (handler *Conversation) GetConversationsByManager(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.Unauthorized())
		return
	}
	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))
	isDone := sql.NullBool{}
	if len(c.Query("is_done")) > 0 {
		isDone.Valid = true
		isDone.Bool, _ = strconv.ParseBool(c.Query("is_done"))
	}
	major := sql.NullBool{}
	if len(c.Query("major")) > 0 {
		major.Valid = true
		major.Bool, _ = strconv.ParseBool(c.Query("major"))
	}
	following := sql.NullBool{}
	if len(c.Query("following")) > 0 {
		following.Valid = true
		following.Bool, _ = strconv.ParseBool(c.Query("following"))
	}

	filter := model.ConversationFilter{
		AppId:          util.ParseQueryArray(c.QueryArray("app_id")),
		ConversationId: util.ParseQueryArray(c.QueryArray("conversation_id")),
		Username:       c.Query("username"),
		PhoneNumber:    c.Query("phone_number"),
		Email:          c.Query("email"),
		IsDone:         isDone,
		Major:          major,
		Following:      following,
	}

	total, result, err := handler.conversationService.GetConversationsByHighLevel(c, res.Data, filter, limit, offset)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.Pagination(result, total, limit, offset))
}

func (handler *Conversation) GetConversationById(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.Unauthorized())
		return
	}
	appId := c.Param("app_id")
	if len(appId) < 1 {
		c.JSON(response.BadRequestMsg("app_id is required"))
		return
	}

	oaId := c.Param("oa_id")
	if len(oaId) < 1 {
		c.JSON(response.BadRequestMsg("oa_id is required"))
		return
	}

	conversationId := c.Param("id")
	if len(conversationId) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	result, err := handler.conversationService.GetConversationById(c, res.Data, appId, conversationId)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OK(result))
}

func (handler *Conversation) PutLabelToConversation(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.Unauthorized())
		return
	}
	labelType := c.Param("label_type")
	if len(labelType) < 1 {
		c.JSON(response.BadRequestMsg("label_type is required"))
		return
	}

	request := model.ConversationLabelRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("put label to conversation payload -> ", &request)

	if err := request.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	labelId, err := service.PutLabelToConversation(c, res.Data, labelType, request)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.OK(map[string]any{
		"id": labelId,
	}))
}

func (handler *Conversation) UpdateUserPreferenceConversation(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.Unauthorized())
		return
	}
	preferenceRequest := model.ConversationPreferenceRequest{}
	if err := c.ShouldBindJSON(&preferenceRequest); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	if err := preferenceRequest.Validate(); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	err := handler.conversationService.UpdateUserPreferenceConversation(c, res.Data, preferenceRequest)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}
