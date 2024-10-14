package service

import (
	"context"
	"fmt"
	"time"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

const batchSize = 2000

type (
	IChatReport interface {
		GetChatWorkReports(ctx context.Context, authUser *model.AuthUser, filter model.MessageFilter, limit, offset int, token string) (total int, result *[]model.ChatWorkReport, err error)
		GetChatGeneralReports(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit, offset int, token string) (total int, result *[]model.ChatGeneralReport, err error)
		ExportWorkReports(ctx context.Context, authUser *model.AuthUser, filter model.MessageFilter, fileType, token string) (exportName string, err error)
		ExportGeneralReports(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, fileType, token string) (exportName string, err error)
	}
	ChatReport struct{}
)

var ChatReportService IChatReport

func NewChatReport() IChatReport {
	return &ChatReport{}
}

func (c *ChatReport) GetChatWorkReports(ctx context.Context, authUser *model.AuthUser, filter model.MessageFilter, limit, offset int, token string) (total int, result *[]model.ChatWorkReport, err error) {
	chatWorkReports := make([]model.ChatWorkReport, 0)
	userWorkMetrics := make(map[string]*model.ChatWorkReport)
	receivingTimesConversation := make(map[string][]time.Time)
	previousDirection := make(map[string]string)
	filter.IsSortedAscending = true // sort messages from oldest to newest
	userIds := make([]string, 0)

	for offsetLogs := 0; limit == -1 || offsetLogs < limit; offsetLogs += batchSize {
		totalMessages, messages, errTmp := repository.MessageESRepo.GetMessages(ctx, authUser.TenantId, ES_INDEX_MESSAGE, filter, batchSize, offsetLogs)
		if errTmp != nil {
			err = errTmp
			log.Error(err)
			return
		}
		if len(*messages) < 1 || offsetLogs > totalMessages {
			break
		}

		for _, message := range *messages {
			if len(message.ConversationId) < 1 {
				continue
			}

			userReport, ok := userWorkMetrics[message.SupporterId]
			if !ok {
				userReport = &model.ChatWorkReport{
					UserId:             message.SupporterId,
					Total:              0,
					Facebook:           model.ChannelWorkPerformanceMetrics{},
					Zalo:               model.ChannelWorkPerformanceMetrics{},
					ConversationExists: make(map[string]bool),
				}
				userWorkMetrics[message.SupporterId] = userReport
				if len(message.SupporterId) > 0 {
					userIds = append(userIds, message.SupporterId)
				}
			}

			if _, ok := receivingTimesConversation[message.ConversationId]; !ok {
				receivingTimesConversation[message.ConversationId] = []time.Time{message.SendTime}
			}
			// customer's message timestamp that's waiting for a reply
			if message.Direction == "receive" && previousDirection[message.ConversationId] == "send" {
				receivingTimesConversation[message.ConversationId] = append(receivingTimesConversation[message.ConversationId], message.SendTime)
			}

			switch message.MessageType {
			case "facebook":
				if _, ok := userReport.ConversationExists[message.ConversationId]; !ok {
					userReport.Facebook.TotalChannels++
					userReport.ConversationExists[message.ConversationId] = true
				}

				handleCollectUserChatReplyMetrics(message, previousDirection, &userReport.Facebook, receivingTimesConversation[message.ConversationId], receivingTimesConversation)
			case "zalo":
				if _, ok := userReport.ConversationExists[message.ConversationId]; !ok {
					userReport.Zalo.TotalChannels++
					userReport.ConversationExists[message.ConversationId] = true
				}

				handleCollectUserChatReplyMetrics(message, previousDirection, &userReport.Zalo, receivingTimesConversation[message.ConversationId], receivingTimesConversation)
			}

			previousDirection[message.ConversationId] = message.Direction
		}
	}

	chatIntegrateSystem, err := GetChatIntegrateSystem(ctx, authUser)
	if err != nil {
		return
	}

	usersInfo, errTmp := GetUsersCrm(chatIntegrateSystem.InfoSystem.ApiGetUserUrl, token, userIds, filter.UnitUuid)
	if errTmp != nil {
		log.Error(errTmp)
	}
	usersList := usersMapListFromArray(usersInfo)

	for _, report := range userWorkMetrics {
		if len(report.UserId) < 1 {
			continue
		}
		if len(filter.UnitUuid) > 0 {
			// if user id ain't in this list users -> they ain't in this unit. So we need to filter out them
			if _, ok := usersList[report.UserId]; !ok {
				continue
			}
		}

		// fill user's full name
		if usersList[report.UserId] != nil {
			report.UserFullname = usersList[report.UserId].FirstName
		}
		report.Total = report.Facebook.TotalChannels + report.Zalo.TotalChannels
		report.Facebook.ReceivingTime.CalculateMetrics()
		report.Facebook.ReplyingTime.CalculateMetrics()
		report.Zalo.ReceivingTime.CalculateMetrics()
		report.Zalo.ReplyingTime.CalculateMetrics()
		chatWorkReports = append(chatWorkReports, *report)
	}

	return len(chatWorkReports), &chatWorkReports, nil
}

func (c *ChatReport) GetChatGeneralReports(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit, offset int, token string) (total int, result *[]model.ChatGeneralReport, err error) {
	chatGeneralReports := make([]model.ChatGeneralReport, 0)
	pageMetrics := make(map[string]*model.ChatGeneralReport)

	for offsetLogs := 0; limit == -1 || offsetLogs < limit; offsetLogs += batchSize {
		totalConversations, conversations, errTmp := repository.ConversationESRepo.GetConversations(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, filter, limit, offset)
		if errTmp != nil {
			err = errTmp
			log.Error(err)
			return
		}
		if len(*conversations) < 1 || offsetLogs > totalConversations {
			break
		}

		for _, conversation := range *conversations {
			pageReport, ok := pageMetrics[conversation.OaId]
			if !ok {
				pageReport = &model.ChatGeneralReport{
					Channel:            conversation.ConversationType,
					OaName:             conversation.OaName,
					TotalConversations: 0,
					Fresh:              model.QuantityRatio{},
					Processing:         model.QuantityRatio{},
					Resolved:           model.QuantityRatio{},
					ConversationExists: make(map[string]bool),
				}
				pageMetrics[conversation.OaId] = pageReport
			}

			if _, ok := pageReport.ConversationExists[conversation.ConversationId]; !ok {
				pageReport.TotalConversations++
				pageReport.ConversationExists[conversation.ConversationId] = true
			}
			if conversation.IsDone {
				pageReport.Resolved.Quantity++
				continue
			}

			messageFilter := model.MessageFilter{
				TenantId:       conversation.TenantId,
				ConversationId: conversation.ConversationId,
				Direction:      "send",
			}
			totalSentMessages, _, errTmp2 := repository.MessageESRepo.GetMessages(ctx, authUser.TenantId, ES_INDEX_MESSAGE, messageFilter, 1, 0)
			if errTmp2 != nil {
				err = errTmp2
				log.Error(err)
				return
			}
			if totalSentMessages > 0 {
				pageReport.Processing.Quantity++
			} else {
				pageReport.Fresh.Quantity++
			}
		}
	}

	for _, report := range pageMetrics {
		report.Fresh.Percent = (report.Fresh.Quantity * 100) / report.TotalConversations
		report.Processing.Percent = (report.Processing.Quantity * 100) / report.TotalConversations
		report.Resolved.Percent = 100 - report.Fresh.Percent - report.Processing.Percent
		chatGeneralReports = append(chatGeneralReports, *report)
	}

	return len(chatGeneralReports), &chatGeneralReports, nil
}

func (c *ChatReport) ExportWorkReports(ctx context.Context, authUser *model.AuthUser, filter model.MessageFilter, fileType, token string) (exportName string, err error) {
	timeStr := util.TimeToStringLayout(time.Now(), "2006_01_02_15_04_05")
	exportName = fmt.Sprintf("User_Chat_Working_Performance_Export_%s.%s", timeStr, fileType)
	exportMap := &model.ExportMap{
		Name:             exportName,
		ExportTime:       util.TimeToString(time.Now()),
		ExportTimeFinish: "",
		TotalRows:        0,
		Status:           "In Progress",
		DomainUuid:       authUser.TenantId,
		UserUuid:         authUser.UserId,
		Type:             fileType,
		Folder:           "",
		Url:              "",
	}

	ct, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	go func() {
		defer cancel()
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
			}
		}()
		// get report data
		_, workReport, err := c.GetChatWorkReports(ct, authUser, filter, -1, 0, token)
		if err != nil {
			panic(err)
		}

		// generate excel file
		if err = c.generateExportUsersWorkPerformance(ct, exportName, fileType, exportMap, workReport); err != nil {
			panic(err)
		}

		chatIntegrateSystem, err := GetChatIntegrateSystem(ct, authUser)
		if err != nil {
			panic(err)
		}

		// send to crm to ack this exported file
		if err = SendExportedFileMetadataToCrm(chatIntegrateSystem.InfoSystem.ApiUrl, token, exportMap); err != nil {
			panic(err)
		}
	}()

	return
}

func (c *ChatReport) ExportGeneralReports(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, fileType, token string) (exportName string, err error) {
	timeStr := util.TimeToStringLayout(time.Now(), "2006_01_02_15_04_05")
	exportName = fmt.Sprintf("Chat_General_Metrics_Export_%s.%s", timeStr, fileType)
	exportMap := &model.ExportMap{
		Name:             exportName,
		ExportTime:       util.TimeToString(time.Now()),
		ExportTimeFinish: "",
		TotalRows:        0,
		Status:           "In Progress",
		DomainUuid:       authUser.TenantId,
		UserUuid:         authUser.UserId,
		Type:             fileType,
		Folder:           "",
		Url:              "",
	}

	ct, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	go func() {
		defer cancel()
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
			}
		}()
		// get report data
		_, workReport, err := c.GetChatGeneralReports(ct, authUser, filter, -1, 0, token)
		if err != nil {
			panic(err)
		}

		// generate excel file
		if err = c.generateExportGeneralMetrics(ct, exportName, fileType, exportMap, workReport); err != nil {
			panic(err)
		}

		chatIntegrateSystem, err := GetChatIntegrateSystem(ct, authUser)
		if err != nil {
			panic(err)
		}

		// send to crm to ack this exported file
		if err = SendExportedFileMetadataToCrm(chatIntegrateSystem.InfoSystem.ApiUrl, token, exportMap); err != nil {
			panic(err)
		}
	}()

	return
}
