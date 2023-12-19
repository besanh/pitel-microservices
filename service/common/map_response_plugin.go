package common

import (
	"strconv"

	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
)

func HandleMapResponsePlugin(plugin string, statusCode int, res any) model.ResponseInboxMarketing {
	resStandard := model.ResponseInboxMarketing{}
	if plugin == "abenla" {
		resAbenla := model.AbenlaSendMessageResponse{}
		if err := util.ParseAnyToAny(res, &resAbenla); err != nil {
			resStandard.Code = "3" // fail
			resStandard.Message = constants.STATUS["fail"]
			return resStandard
		}
		code := strconv.Itoa(resAbenla.Code)
		resStandard.Code = constants.STANDARD_CODE_ABENLA_TO_TEL4VN[code]
		resStandard.Message = constants.MESSAGE_STATUS_ABENLA[resAbenla.Code]
		return resStandard
	} else if plugin == "incom" {
		resIncom := model.IncomSendMessageResponse{}
		if err := util.ParseAnyToAny(res, &resIncom); err != nil {
			resStandard.Code = "3" // fail
			resStandard.Message = constants.STATUS["fail"]
			return resStandard
		}
		resStandard.Code = constants.STANDARD_CODE_INCOM_TO_TEL4VN[resIncom.Status]
		status, _ := strconv.Atoi(resIncom.Status)
		resStandard.Message = constants.MESSAGE_STATUS_INCOM[status]
		return resStandard
	} else if plugin == "fpt" {
		if statusCode != 200 {
			resFpt := model.FptResponseError{}
			if err := util.ParseAnyToAny(res, &resFpt); err != nil {
				resStandard.Code = "3" // fail
				resStandard.Message = constants.STATUS["fail"]
				return resStandard
			}
			code := strconv.Itoa(resFpt.Err)
			resStandard.Code = constants.STANDARD_CODE_FPT_TO_TEL4VN[code]
			resStandard.Message = constants.MESSAGE_STATUS_FPT[resFpt.Err]
			return resStandard
		} else {
			resFpt := model.FptSendMessageResponse{}
			if err := util.ParseAnyToAny(res, &resFpt); err != nil {
				resStandard.Code = "3" // fail
				resStandard.Message = constants.STATUS["fail"]
				return resStandard
			}
			resStandard.Code = constants.STANDARD_CODE["1"]
			resStandard.Message = constants.STANDARD_CODE["success"]
			return resStandard
		}
	}
	resStandard.Code = "3" // fail
	resStandard.Message = constants.STATUS["fail"]
	return resStandard
}
