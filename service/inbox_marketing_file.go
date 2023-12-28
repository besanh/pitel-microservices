package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"github.com/xuri/excelize/v2"
)

type ResponseExportFile struct {
	FileName string
	Message  string
	Status   string
}

func (s *InboxMarketing) PostExportReportInboxMarketing(ctx context.Context, authUser *model.AuthUser, fileType string, filter model.InboxMarketingFilter) (int, any) {
	total, _, err := repository.InboxMarketingESRepo.GetReport(ctx, authUser.TenantId, ES_INDEX, -1, 0, filter)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if total > 30000 {
		return response.ServiceUnavailableMsg("please export maximum 30k rows")
	}
	timeStr := util.TimeToStringLayout(time.Now(), "2006_01_02_15_04_05")
	exportName := "Report_Inbox_Marketing_" + timeStr + "." + fileType
	exportMap := []string{exportName, util.TimeToString(time.Now()), "", "0", "In Progress", authUser.TenantId, "inbox_marketing"}
	if err := SetExportValue(authUser.TenantId, exportName, exportMap); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	switch fileType {
	case "csv":
		go s.generateCSVReportInboxMarketing(authUser.TenantId, ES_INDEX, exportName, exportMap, filter)
	case "xlsx":
		go s.generateExcelReportInboxMarketing(authUser.TenantId, ES_INDEX, exportName, exportMap, filter)
	}

	return response.Created(map[string]any{
		"export_name": exportName,
		"status":      "In Progress",
	})
}

func (s *InboxMarketing) generateCSVReportInboxMarketing(tenantId, index, exportName string, exportMap []string, filter model.InboxMarketingFilter) {
	rows, total, err := s.handleProcessDataExcelize(tenantId, index, exportName, exportMap, filter)
	if err != nil {
		log.Error(err)
		return
	}
	w, f, err := getCSVWriter(exportName)
	if err != nil {
		log.Error(err)
		return
	}
	defer f.Close()
	bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
	f.Write(bomUtf8)
	var errExport error
	if errExport != nil {
		log.Error(errExport)
		exportMap[3] = util.TimeToString(time.Now())
		exportMap[4] = "Error"
		if err := SetExportValue(tenantId, exportName, exportMap); err != nil {
			log.Error(err)
			return
		}
		return
	}

	for _, row := range rows {
		records := make([]string, len(row))
		for k, v := range row {
			records[k] = fmt.Sprintf("%v", v)
		}
		if err = w.Write(records); err != nil {
			log.Error(err)
			return
		}
	}
	w.Flush()
	if err := f.Close(); err != nil {
		log.Error(err)
		return
	}

	exportMap[2] = util.TimeToString(time.Now())
	exportMap[3] = fmt.Sprintf("%d", total)
	exportMap[4] = "Done"
	if err := SetExportValue(tenantId, exportName, exportMap); err != nil {
		log.Error(err)
		return
	}
}

func (s *InboxMarketing) generateExcelReportInboxMarketing(tenantId, index, exportName string, exportMap []string, filter model.InboxMarketingFilter) {
	rows, total, err := s.handleProcessDataExcelize(tenantId, index, exportName, exportMap, filter)
	if err != nil {
		log.Error(err)
		return
	}
	idx := 1
	file := excelize.NewFile()
	streamWriter, err := file.NewStreamWriter("sheet1")
	if err != nil {
		log.Error(err)
		return
	}
	for _, row := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, idx)
		if err := streamWriter.SetRow(cell, row); err != nil {
			log.Error(err)
			return
		}
		idx += 1
	}
	if err := streamWriter.Flush(); err != nil {
		log.Error(err)
		return
	}
	exportDir := constants.EXPORT_DIR + "inbox_marketing/"
	_ = os.MkdirAll(filepath.Dir(exportDir), 0755)
	if err := file.SaveAs(exportDir + exportName); err != nil {
		log.Error(err)
		return
	}
	exportMap[2] = util.TimeToString(time.Now())
	exportMap[3] = fmt.Sprintf("%d", total)
	exportMap[4] = "Done"
	if err := SetExportValue(tenantId, exportName, exportMap); err != nil {
		log.Error(err)
		return
	}
}

func (s *InboxMarketing) handleProcessDataExcelize(tenantId, index, exportName string, exportMap []string, filter model.InboxMarketingFilter) ([][]any, int, error) {
	filter.Limit = 1
	ctx := context.Background()
	total, _, err := repository.InboxMarketingESRepo.GetReport(ctx, tenantId, index, filter.Limit, filter.Offset, filter)
	if err != nil {
		exportMap[2] = util.TimeToString(time.Now())
		exportMap[4] = "Error"
		if err := SetExportValue(tenantId, exportName, exportMap); err != nil {
			log.Error(err)
			return nil, total, err
		}
	}

	headers := map[string]ColumnInfo{
		"phone_number":   {Index: 0, Name: "Phone Number"},
		"template_code":  {Index: 1, Name: "Template Code"},
		"channel":        {Index: 2, Name: "Channel"},
		"status":         {Index: 3, Name: "Status"},
		"network":        {Index: 4, Name: "Network"},
		"quantity":       {Index: 5, Name: "Quantity"},
		"error_code":     {Index: 6, Name: "Error Code"},
		"send_time":      {Index: 7, Name: "Send Time"},
		"is_charged_zns": {Index: 8, Name: "Charged(zns)"},
		"param":          {Index: 9, Name: "Param"},
	}

	rows := make([][]any, 0)
	for int(filter.Offset) < total {
		filter.Limit = 200
		total, inboxMarketings, err := repository.InboxMarketingESRepo.GetReport(ctx, tenantId, index, filter.Limit, filter.Offset, filter)
		if err != nil {
			exportMap[3] = util.TimeToString(time.Now())
			exportMap[4] = "Error"
			if err := SetExportValue(tenantId, exportName, exportMap); err != nil {
				log.Error(err)
				return nil, total, err
			}
		}
		for _, inboxMarketing := range inboxMarketings {
			row := s.handleRowData(inboxMarketing, &headers)
			rows = append(rows, row)
		}
		filter.Offset += filter.Limit
		percentComplete := (float64(filter.Offset) / float64(total)) * 100
		exportMap[3] = fmt.Sprintf("%d", total)
		exportMap[4] = "In Progress (" + fmt.Sprintf("%.2f", percentComplete) + "%)"
		if err := SetExportValue(tenantId, exportName, exportMap); err != nil {
			log.Error(err)
			return nil, total, err
		}
	}
	columnsMap := make(map[int]string)
	for _, v := range headers {
		columnsMap[v.Index] = v.Name
	}
	columns := sortMapAny(columnsMap)
	rowsCopy := make([][]any, len(rows)+1)
	rowsCopy[0] = columns
	if len(rows) > 0 {
		copy(rowsCopy[1:], rows)
	}
	return rowsCopy, total, nil
}

func sortMapAny(tmpMap map[int]string) []any {
	var result []any
	keys := make([]int, 0, len(tmpMap))
	for key := range tmpMap {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		result = append(result, tmpMap[key])
	}
	return result
}

func (s *InboxMarketing) handleRowData(inboxMarketing model.InboxMarketingLogReport, headersAdd *map[string]ColumnInfo) []any {
	headers := *headersAdd
	row := make([]any, len(headers))
	row[headers["phone_number"].Index] = inboxMarketing.PhoneNumber
	row[headers["template_code"].Index] = inboxMarketing.TemplateCode
	row[headers["channel"].Index] = inboxMarketing.Channel
	row[headers["status"].Index] = inboxMarketing.Status
	network := constants.NETWORKS[strconv.Itoa(inboxMarketing.TelcoId)]
	if len(network) > 0 {
		row[headers["network"].Index] = network
	} else {
		row[headers["network"].Index] = "unknown"
	}
	row[headers["quantity"].Index] = inboxMarketing.Quantity
	row[headers["error_code"].Index] = inboxMarketing.ErrorCode
	chargedZns := ""
	if inboxMarketing.IsChargedZns {
		chargedZns = "charged"
	} else {
		if inboxMarketing.Channel == "zns" {
			chargedZns = "not charged"
		}
	}
	row[headers["is_charged_zns"].Index] = chargedZns
	row[headers["send_time"].Index] = inboxMarketing.CreatedAt
	row[headers["param"].Index] = string(inboxMarketing.ListParam)
	headersAdd = &headers

	return row
}

func getCSVWriter(fileName string) (*csv.Writer, *os.File, error) {
	exportDir := constants.EXPORT_DIR + "inbox_marketing/"
	_ = os.MkdirAll(filepath.Dir(exportDir), 0755)
	f, err := os.Create(exportDir + fileName)
	if err != nil {
		return nil, nil, err
	}
	w := csv.NewWriter(f)
	w.UseCRLF = true
	return w, f, nil
}
