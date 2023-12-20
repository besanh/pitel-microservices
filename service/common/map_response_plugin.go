package common

import (
	"strconv"

	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"golang.org/x/exp/slices"
)

var (
	CODE_ABENLA_SUCESS = []string{"106"}
	CODE_INCOM_SUCCESS = []string{"1"}
)

func HandleMapResponsePlugin(plugin string, id string, statusCode int, res any) model.ResponseInboxMarketing {
	resStandard := model.ResponseInboxMarketing{
		Id: id,
	}
	if plugin == "abenla" {
		resAbenla := model.AbenlaSendMessageResponse{}
		if err := util.ParseAnyToAny(res, &resAbenla); err != nil {
			resStandard.Code = "3" // fail
			resStandard.Message = constants.MESSAGE_TEL4VN["fail"]
			return resStandard
		}
		code := strconv.Itoa(resAbenla.Code)
		if standardCode, exist := constants.STANDARD_CODE_ABENLA_TO_TEL4VN[code]; exist {
			resStandard.Code = standardCode
		} else {
			resStandard.Code = "3" // fail
			resStandard.Message = constants.MESSAGE_TEL4VN["fail"]
		}

		if message, exist := constants.MESSAGE_ABENLA[resAbenla.Code]; exist {
			resStandard.Message = message
		} else {
			resStandard.Message = constants.MESSAGE_TEL4VN["fail"]
		}

		if slices.Contains[[]string](CODE_ABENLA_SUCESS, code) {
			resStandard.Status = "Sucess"
		} else {
			resStandard.Status = "Fail"
		}
		resStandard.Quantity = resAbenla.SmsPerMessage

		return resStandard
	} else if plugin == "incom" {
		resIncom := model.IncomSendMessageResponse{}
		if err := util.ParseAnyToAny(res, &resIncom); err != nil {
			resStandard.Code = "3" // fail
			resStandard.Message = constants.MESSAGE_TEL4VN["fail"]
			return resStandard
		}
		if standardCode, exist := constants.STANDARD_CODE_INCOM_TO_TEL4VN[resIncom.Status]; exist {
			resStandard.Code = standardCode
		} else {
			resStandard.Code = "3" // fail
			resStandard.Message = constants.MESSAGE_TEL4VN["fail"]
		}
		status, _ := strconv.Atoi(resIncom.Status)
		if message, exist := constants.MESSAGE_INCOM[status]; exist {
			resStandard.Message = message
		} else {
			resStandard.Message = constants.MESSAGE_TEL4VN["fail"]
		}

		if slices.Contains[[]string](CODE_INCOM_SUCCESS, resIncom.Status) {
			resStandard.Status = "Sucess"
		} else {
			resStandard.Status = "Fail"
		}

		return resStandard
	} else if plugin == "fpt" {
		if statusCode != 200 {
			resFpt := model.FptResponseError{}
			if err := util.ParseAnyToAny(res, &resFpt); err != nil {
				resStandard.Code = "3" // fail
				resStandard.Message = constants.MESSAGE_TEL4VN["fail"]
				return resStandard
			}
			code := strconv.Itoa(resFpt.Err)
			if standardCode, exist := constants.STANDARD_CODE_FPT_TO_TEL4VN[code]; exist {
				resStandard.Code = standardCode
			} else {
				resStandard.Code = "3" // fail
				resStandard.Message = constants.MESSAGE_TEL4VN["fail"]
			}

			if message, exist := constants.MESSAGE_FPT[resFpt.Err]; exist {
				resStandard.Message = message
			} else {
				resStandard.Message = constants.MESSAGE_TEL4VN["fail"]
			}

			resStandard.Status = "Fail"

			return resStandard
		} else {
			resFpt := model.FptSendMessageResponse{}
			if err := util.ParseAnyToAny(res, &resFpt); err != nil {
				resStandard.Code = "3" // fail
				resStandard.Message = constants.MESSAGE_TEL4VN["fail"]
				return resStandard
			}
			resStandard.Code = constants.STANDARD_CODE["1"]
			resStandard.Message = constants.MESSAGE_TEL4VN["success"]
			resStandard.Status = "Success"
			return resStandard
		}
	}
	resStandard.Code = "3" // fail
	resStandard.Message = constants.MESSAGE_TEL4VN["fail"]
	return resStandard
}
