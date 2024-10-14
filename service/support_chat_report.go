package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

// collect user's replying message metrics
func handleCollectUserChatReplyMetrics(message model.Message, previousDirection map[string]string, channelMetrics *model.ChannelWorkPerformanceMetrics, receivingTimes []time.Time, receivingTimesConversation map[string][]time.Time) {
	if message.Direction == "send" && previousDirection[message.ConversationId] == "receive" {
		if len(receivingTimesConversation[message.ConversationId]) == 1 {
			channelMetrics.ReceivingTime.AddTimestamp(message.SendTime.Sub(receivingTimes[len(receivingTimes)-1]))
		}
		channelMetrics.ReplyingTime.AddTimestamp(message.SendTime.Sub(receivingTimes[len(receivingTimes)-1]))
	} else if message.Direction == "send" && previousDirection[message.ConversationId] == "send" {
		if len(channelMetrics.ReceivingTime.Timestamps) < 1 {
			channelMetrics.ReceivingTime.AddTimestamp(message.SendTime.Sub(receivingTimes[len(receivingTimes)-1]))
		}
		if len(channelMetrics.ReplyingTime.Timestamps) < len(receivingTimes) {
			// not to add too much reply time metrics
			channelMetrics.ReplyingTime.AddTimestamp(message.SendTime.Sub(receivingTimes[len(receivingTimes)-1]))
		}
	}
}

func (c *ChatReport) generateExportUsersWorkPerformance(ctx context.Context, exportName, fileType string, exportMap *model.ExportMap, chatReport *[]model.ChatWorkReport) (err error) {
	headers := [][]string{
		{
			"Nhân viên hỗ trợ",
			"Tổng",
			"Facebook",
			"",
			"",
			"",
			"",
			"",
			"",
			"Zalo",
			"",
			"",
			"",
			"",
			"",
			"",
		},
		{
			"",
			"",
			"Hỗ trợ (lượt)",
			"Thời gian tiếp nhận",
			"",
			"",
			"Thời gian chờ phản hổi",
			"",
			"",
			"Hỗ trợ (lượt)",
			"Thời gian tiếp nhận",
			"",
			"",
			"Thời gian chờ phản hổi",
			"",
			"",
		},
		{
			"",
			"",
			"",
			"Nhanh nhất",
			"Chậm nhất",
			"Trung bình",
			"Nhanh nhất",
			"Chậm nhất",
			"Trung bình",
			"",
			"Nhanh nhất",
			"Chậm nhất",
			"Trung bình",
			"Nhanh nhất",
			"Chậm nhất",
			"Trung bình",
		},
	}
	rows := make([][]string, 0)
	limitPerPart := 50

	for offset := 0; offset < len(*chatReport); offset += limitPerPart {
		for _, report := range *chatReport {
			userFullName := report.UserId
			if len(report.UserFullname) > 0 {
				userFullName = report.UserFullname
			}
			row := make([]string, 0)
			row = append(row,
				userFullName,
				strconv.Itoa(report.Total),
				strconv.Itoa(report.Facebook.TotalChannels),
				util.ConvertMillisToTimeString(report.Facebook.ReceivingTime.Fastest),
				util.ConvertMillisToTimeString(report.Facebook.ReceivingTime.Slowest),
				util.ConvertMillisToTimeString(report.Facebook.ReceivingTime.Average),
				util.ConvertMillisToTimeString(report.Facebook.ReplyingTime.Fastest),
				util.ConvertMillisToTimeString(report.Facebook.ReplyingTime.Slowest),
				util.ConvertMillisToTimeString(report.Facebook.ReplyingTime.Average),
				strconv.Itoa(report.Zalo.TotalChannels),
				util.ConvertMillisToTimeString(report.Zalo.ReceivingTime.Fastest),
				util.ConvertMillisToTimeString(report.Zalo.ReceivingTime.Slowest),
				util.ConvertMillisToTimeString(report.Zalo.ReceivingTime.Average),
				util.ConvertMillisToTimeString(report.Zalo.ReplyingTime.Fastest),
				util.ConvertMillisToTimeString(report.Zalo.ReplyingTime.Slowest),
				util.ConvertMillisToTimeString(report.Zalo.ReplyingTime.Average),
			)
			rows = append(rows, row)
		}

		if offset < len(*chatReport) {
			percentComplete := (float64(offset) / float64(len(*chatReport))) * 100
			exportMap.TotalRows = len(*chatReport)
			exportMap.Status = fmt.Sprintf("In Progress (%.2f%%)", percentComplete)
		}
	}

	var buf *bytes.Buffer
	if fileType == "xlsx" {
		buf, err = util.HandleExcelStreamWriter(headers, rows, "A1", "A3", "B1", "B3", "C1", "I1", "J1", "P1",
			"C2", "C3", "D2", "F2", "G2", "I2", "J2", "J3", "K2", "M2", "N2", "P2")
		if err != nil {
			log.Error(err)
			return
		}
	} else if fileType == "csv" {
		buf, err = util.HandleCSVStreamWriter(exportName, headers, rows, "A1", "A3", "B1", "B3", "C1", "I1", "J1", "P1",
			"C2", "C3", "D2", "F2", "G2", "I2", "J2", "J3", "K2", "M2", "N2", "P2")
		if err != nil {
			log.Error(err)
			return
		}
	}
	fileUrl, err := uploadFileToStorage(ctx, buf, "/bss-message/v1/share-info/image/", exportName)
	if err != nil {
		log.Error(err)
		return
	}

	exportMap.Url = fileUrl
	exportMap.ExportTimeFinish = util.TimeToString(time.Now())
	exportMap.TotalRows = len(*chatReport)
	exportMap.Status = "Done"
	return
}

func (c *ChatReport) generateExportGeneralMetrics(ctx context.Context, exportName, fileType string, exportMap *model.ExportMap, chatReport *[]model.ChatGeneralReport) (err error) {
	headers := [][]string{
		{
			"Kênh",
			"Tên trang",
			"Tổng hội thoại",
			"Mới",
			"",
			"Đang xử lý",
			"",
			"Đã xử lý",
			"",
		},
		{
			"",
			"",
			"",
			"Số lượng",
			"Tỷ trọng",
			"Số lượng",
			"Tỷ trọng",
			"Số lượng",
			"Tỷ trọng",
		},
	}
	rows := make([][]string, 0)
	limitPerPart := 50

	for offset := 0; offset < len(*chatReport); offset += limitPerPart {
		for _, report := range *chatReport {
			row := make([]string, 0)
			row = append(row,
				report.Channel,
				report.OaName,
				strconv.Itoa(report.TotalConversations),
				strconv.Itoa(report.Fresh.Quantity),
				fmt.Sprintf("%.2f%%", float64(report.Fresh.Percent)),
				strconv.Itoa(report.Processing.Quantity),
				fmt.Sprintf("%.2f%%", float64(report.Processing.Percent)),
				strconv.Itoa(report.Resolved.Quantity),
				fmt.Sprintf("%.2f%%", float64(report.Resolved.Percent)),
			)
			rows = append(rows, row)
		}

		if offset < len(*chatReport) {
			percentComplete := (float64(offset) / float64(len(*chatReport))) * 100
			exportMap.TotalRows = len(*chatReport)
			exportMap.Status = fmt.Sprintf("In Progress (%.2f%%)", percentComplete)
		}
	}

	var buf *bytes.Buffer
	if fileType == "xlsx" {
		buf, err = util.HandleExcelStreamWriter(headers, rows, "A1", "A2", "B1", "B2", "C1", "C2", "D1", "E1", "F1", "G1", "H1", "I1")
		if err != nil {
			log.Error(err)
			return
		}
	} else if fileType == "csv" {
		buf, err = util.HandleCSVStreamWriter(exportName, headers, rows, "A1", "A2", "B1", "B2", "C1", "C2", "D1", "E1", "F1", "G1", "H1", "I1")
		if err != nil {
			log.Error(err)
			return
		}
	}
	fileUrl, err := uploadFileToStorage(ctx, buf, "/bss-message/v1/share-info/image/", exportName)
	if err != nil {
		log.Error(err)
		return
	}

	exportMap.Url = fileUrl
	exportMap.ExportTimeFinish = util.TimeToString(time.Now())
	exportMap.TotalRows = len(*chatReport)
	exportMap.Status = "Done"
	return
}

func GetChatIntegrateSystem(ctx context.Context, authUser *model.AuthUser) (chatIntegrateSystem *model.ChatIntegrateSystem, err error) {
	chatIntegrateSystem = &model.ChatIntegrateSystem{}
	chatIScache := cache.RCache.Get(CHAT_INTEGRATE_SYSTEM + "_" + authUser.SystemId)
	if chatIScache != nil {
		if err = json.Unmarshal([]byte(chatIScache.(string)), chatIntegrateSystem); err != nil {
			log.Error(err)
			return
		}
	} else {
		_, chatIntegrateSystems, errTmp := repository.ChatIntegrateSystemRepo.GetIntegrateSystems(ctx, repository.DBConn, model.ChatIntegrateSystemFilter{
			SystemId: authUser.SystemId}, 1, 0)
		if errTmp != nil {
			log.Error(errTmp)
			err = errTmp
			return
		} else if len(*chatIntegrateSystems) < 1 {
			err = errors.New("invalid system id " + authUser.SystemId)
			log.Error(err)
			return
		}

		chatIntegrateSystem = &(*chatIntegrateSystems)[0]
	}
	return
}

func GetUsersCrm(apiUrl, token string, userIds []string, unitUuid string) (result []model.AuthUserInfo, err error) {
	client := resty.New()
	request := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+token).
		SetQueryParam("limit", "-1").
		SetQueryParam("offset", "0").
		SetQueryParam("unit_uuid", unitUuid)
	if len(unitUuid) < 1 {
		request.SetQueryParamsFromValues(url.Values{
			"user_uuid": userIds,
		})
	}
	res, err := request.Get(apiUrl)
	if err != nil {
		return
	}
	if res.IsError() {
		err = fmt.Errorf("unexpected status code: %d", res.StatusCode())
		return
	}

	bodyRaw := make(map[string]any)
	if err = json.Unmarshal(res.Body(), &bodyRaw); err != nil {
		return
	}
	if err = util.ParseAnyToAny(bodyRaw["data"], &result); err != nil {
		return
	}
	return
}

func usersMapListFromArray(data []model.AuthUserInfo) (result map[string]*model.AuthUserInfo) {
	result = make(map[string]*model.AuthUserInfo)
	for _, user := range data {
		result[user.UserUuid] = &user
		result[user.UserUuid].FirstName = fmt.Sprintf("%s %s %s", user.FirstName, user.MiddleName, user.LastName)
	}
	return
}

func SendExportedFileMetadataToCrm(apiUrl, token string, exportMap *model.ExportMap) (err error) {
	client := resty.New()
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+token).
		SetBody(exportMap).
		Post(apiUrl + "/v1/crm/export/external")
	if err != nil {
		return
	}
	if res.IsError() {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode())
	}

	return
}
