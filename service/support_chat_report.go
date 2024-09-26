package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"github.com/tel4vn/fins-microservices/service/common"
)

func (c *ChatReport) generateExportUsersWorkPerformance(ctx context.Context, tenantId, userId, exportName, fileType string, exportMap *model.ExportMap, chatReport *[]model.ChatWorkReport) (err error) {
	headers := make([]string, 0)
	headers = append(headers,
		"Nhân viên hỗ trợ",
		"Tổng",
		"Facebook.Hỗ trợ (lượt)",
		"Facebook.Thời gian tiếp nhận.Nhanh nhất",
		"Facebook.Thời gian tiếp nhận.Chậm nhất",
		"Facebook.Thời gian tiếp nhận.Trung bình",
		"Facebook.Thời gian chờ phản hổi.Nhanh nhất",
		"Facebook.Thời gian chờ phản hổi.Chậm nhất",
		"Facebook.Thời gian chờ phản hổi.Trung bình",
		"Zalo.Hỗ trợ (lượt)",
		"Zalo.Thời gian tiếp nhận.Nhanh nhất",
		"Zalo.Thời gian tiếp nhận.Chậm nhất",
		"Zalo.Thời gian tiếp nhận.Trung bình",
		"Zalo.Thời gian chờ phản hổi.Nhanh nhất",
		"Zalo.Thời gian chờ phản hổi.Chậm nhất",
		"Zalo.Thời gian chờ phản hổi.Trung bình",
	)
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
				strconv.Itoa(report.Facebook.ReceivingTime.Fastest),
				strconv.Itoa(report.Facebook.ReceivingTime.Slowest),
				strconv.Itoa(report.Facebook.ReceivingTime.Average),
				strconv.Itoa(report.Facebook.ReplyingTime.Fastest),
				strconv.Itoa(report.Facebook.ReplyingTime.Slowest),
				strconv.Itoa(report.Facebook.ReplyingTime.Average),
				strconv.Itoa(report.Zalo.TotalChannels),
				strconv.Itoa(report.Zalo.ReceivingTime.Fastest),
				strconv.Itoa(report.Zalo.ReceivingTime.Slowest),
				strconv.Itoa(report.Zalo.ReceivingTime.Average),
				strconv.Itoa(report.Zalo.ReplyingTime.Fastest),
				strconv.Itoa(report.Zalo.ReplyingTime.Slowest),
				strconv.Itoa(report.Zalo.ReplyingTime.Average),
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
		buf, err = common.HandleExcelStreamWriter(headers, rows)
		if err != nil {
			log.Error(err)
			return
		}
	} else if fileType == "csv" {
		buf, err = common.HandleCSVStreamWriter(exportName, headers, rows)
		if err != nil {
			log.Error(err)
			return
		}
	}
	url, err := uploadFileToStorage(ctx, buf, "/bss-message/v1/share-info/image/", exportName)
	if err != nil {
		log.Error(err)
		return
	}

	exportMap.Url = url
	exportMap.ExportTimeFinish = util.TimeToString(time.Now())
	exportMap.TotalRows = len(*chatReport)
	exportMap.Status = "Done"
	return
}

func (c *ChatReport) generateExportGeneralMetrics(ctx context.Context, tenantId, userId, exportName, fileType string, exportMap *model.ExportMap, chatReport *[]model.ChatGeneralReport) (err error) {
	headers := make([]string, 0)
	headers = append(headers,
		"Kênh",
		"Tên trang",
		"Tổng hội thoại",
		"Mới.Số lượng",
		"Mới.Tỷ trọng",
		"Đang xử lý.Số lượng",
		"Đang xử lý.Tỷ trọng",
		"Đã xử lý.Số lượng",
		"Đã xử lý.Tỷ trọng",
	)
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
				fmt.Sprintf("%.2f%%", float64(report.Fresh.Percent)/100),
				strconv.Itoa(report.Processing.Quantity),
				fmt.Sprintf("%.2f%%", float64(report.Processing.Percent)/100),
				strconv.Itoa(report.Resolved.Quantity),
				fmt.Sprintf("%.2f%%", float64(report.Resolved.Percent)/100),
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
		buf, err = common.HandleExcelStreamWriter(headers, rows)
		if err != nil {
			log.Error(err)
			return
		}
	} else if fileType == "csv" {
		buf, err = common.HandleCSVStreamWriter(exportName, headers, rows)
		if err != nil {
			log.Error(err)
			return
		}
	}
	url, err := uploadFileToStorage(ctx, buf, "/bss-message/v1/share-info/image/", exportName)
	if err != nil {
		log.Error(err)
		return
	}

	exportMap.Url = url
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

func GetUsersCrm(apiUrl, token string) (result []model.AuthUserInfo, err error) {
	client := resty.New()
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+token).
		SetQueryParam("limit", "-1").
		SetQueryParam("offset", "0").
		Get(apiUrl)
	if err != nil {
		return
	}
	if res.StatusCode() < 200 || res.StatusCode() >= 300 {
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
	if res.StatusCode() < 200 || res.StatusCode() >= 300 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode())
	}

	return
}
