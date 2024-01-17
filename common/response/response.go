package response

import (
	"errors"
	"net/http"

	validator "github.com/bufbuild/protovalidate-go"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ERR_TOKEN_IS_EMPTY       = "token is empty"
	ERR_TOKEN_IS_INVALID     = "token is invalid"
	ERR_TOKEN_IS_EXPIRED     = "token is expired"
	ERR_EMPTY_CONN           = "empty connection"
	ERR_DATA_NOT_FOUND       = "data not found"
	ERR_DATA_INVALID         = "data is invalid"
	ERR_EXAMPLE_NOT_FOUND    = "example not found"
	ERR_EXAMPLE_INVALID      = "example is invalid"
	ERR_CAMPAIGN_NOT_FOUND   = "campaign not found"
	ERR_CAMPAIGN_INVALID     = "campaign is invalid"
	ERR_MASTERDATA_NOT_FOUND = "masterdata not found"
	ERR_MASTERDATA_INVALID   = "masterdata is invalid"
	ERR_INSERT_FAILED        = "insert failed"
	ERR_GET_FAILED           = "get failed"
	ERR_PUT_FAILED           = "put failed"
	ERR_PATCH_FAILED         = "patch failed"
	ERR_DELETE_FAILED        = "delete failed"
	ERR_VALIDATION_FAILED    = "validation failed"
)

var MAP_ERR_RESPONSE = map[string]struct {
	GRPC_Code codes.Code
	Code      string
	Message   string
}{
	ERR_TOKEN_IS_EMPTY: {
		GRPC_Code: codes.Unauthenticated,
		Code:      "ERR_UNAUTHORIZE",
		Message:   ERR_TOKEN_IS_EMPTY,
	},
	ERR_TOKEN_IS_INVALID: {
		GRPC_Code: codes.Unauthenticated,
		Code:      "ERR_UNAUTHORIZE",
		Message:   ERR_TOKEN_IS_INVALID,
	},
	ERR_TOKEN_IS_EXPIRED: {
		GRPC_Code: codes.Unauthenticated,
		Code:      "ERR_UNAUTHORIZE",
		Message:   ERR_TOKEN_IS_EXPIRED,
	},
	ERR_EMPTY_CONN: {
		GRPC_Code: codes.OK,
		Code:      "ERR_EMPTY_CONN",
		Message:   ERR_EMPTY_CONN,
	},
	ERR_EXAMPLE_NOT_FOUND: {
		GRPC_Code: codes.OK,
		Code:      "ERR_EXAMPLE_NOT_FOUND",
		Message:   ERR_EXAMPLE_NOT_FOUND,
	},
	ERR_EXAMPLE_INVALID: {
		GRPC_Code: codes.OK,
		Code:      "ERR_EXAMPLE_INVALID",
		Message:   ERR_EXAMPLE_INVALID,
	},
	ERR_DATA_NOT_FOUND: {
		GRPC_Code: codes.NotFound,
		Code:      "ERR_DATA_NOT_FOUND",
		Message:   ERR_DATA_NOT_FOUND,
	},
	ERR_DATA_INVALID: {
		GRPC_Code: codes.Unavailable,
		Code:      "ERR_DATA_INVALID",
		Message:   ERR_DATA_INVALID,
	},
	ERR_INSERT_FAILED: {
		GRPC_Code: codes.Unavailable,
		Code:      "ERR_INSERT_FAILED",
		Message:   ERR_INSERT_FAILED,
	},
	ERR_GET_FAILED: {
		GRPC_Code: codes.Unavailable,
		Code:      "ERR_GET_FAILED",
		Message:   ERR_GET_FAILED,
	},
	ERR_PUT_FAILED: {
		GRPC_Code: codes.Unavailable,
		Code:      "ERR_PUT_FAILED",
		Message:   ERR_PUT_FAILED,
	},
	ERR_PATCH_FAILED: {
		GRPC_Code: codes.Unavailable,
		Code:      "ERR_PATCH_FAILED",
		Message:   ERR_PATCH_FAILED,
	},
	ERR_DELETE_FAILED: {
		GRPC_Code: codes.Unavailable,
		Code:      "ERR_DELETE_FAILED",
		Message:   ERR_DELETE_FAILED,
	},
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func HandleGRPCErrResponse(err error) (code codes.Code, response any) {
	if e, ok := status.FromError(err); ok {
		if e, ok := MAP_ERR_RESPONSE[e.Message()]; ok {
			return e.GRPC_Code, ErrorResponse{
				Code:    e.Code,
				Message: e.Message,
				Error:   err.Error(),
			}
		}
	}

	return codes.Unavailable, map[string]any{
		"message": "internal server error",
		"code":    "SERVICE_UNAVAILABLE",
		"error":   err.Error(),
	}
}

type ValidatorFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func HandleValidatorPBError(err error) (isString bool, val any) {
	var valErr *validator.ValidationError
	if ok := errors.As(err, &valErr); ok {
		tmp := make(map[string]any)
		for _, v := range valErr.ToProto().GetViolations() {
			tmp[v.GetFieldPath()] = v.GetMessage()
		}
		val = tmp
	} else {
		isString = true
		if err.Error() == "EOF" {
			val = "body is invalid"
		} else {
			val = err.Error()
		}
	}
	return
}

// response for gin api

func OKResponse() (int, any) {
	return http.StatusOK, map[string]any{
		"message": "SUCCESS",
		"code":    http.StatusText(http.StatusOK),
	}
}

func BadRequest() (int, any) {
	return http.StatusBadRequest, map[string]any{
		"error":   http.StatusText(http.StatusBadRequest),
		"code":    http.StatusText(http.StatusBadRequest),
		"message": http.StatusText(http.StatusBadRequest),
	}
}

func BadRequestMsg(msg any) (int, any) {
	return http.StatusBadRequest, map[string]any{
		"error":   http.StatusText(http.StatusBadRequest),
		"code":    http.StatusText(http.StatusBadRequest),
		"message": msg,
	}
}

func NotFound() (int, any) {
	return http.StatusNotFound, map[string]any{
		"error":   http.StatusText(http.StatusNotFound),
		"code":    http.StatusText(http.StatusNotFound),
		"message": http.StatusText(http.StatusNotFound),
	}
}

func NotFoundMsg(msg any) (int, any) {
	return http.StatusNotFound, map[string]any{
		"error":   http.StatusText(http.StatusNotFound),
		"code":    http.StatusText(http.StatusNotFound),
		"message": msg,
	}
}

func Forbidden() (int, any) {
	return http.StatusForbidden, map[string]any{
		"error":   "Do not have permission for the request.",
		"code":    http.StatusText(http.StatusForbidden),
		"message": http.StatusText(http.StatusForbidden),
	}
}

func Unauthorized() (int, any) {
	return http.StatusUnauthorized, map[string]any{
		"error":   http.StatusText(http.StatusUnauthorized),
		"code":    http.StatusText(http.StatusUnauthorized),
		"message": http.StatusText(http.StatusUnauthorized),
	}
}

func ServiceUnavailable() (int, any) {
	return http.StatusServiceUnavailable, map[string]any{
		"error":   http.StatusText(http.StatusServiceUnavailable),
		"code":    http.StatusText(http.StatusServiceUnavailable),
		"message": http.StatusText(http.StatusServiceUnavailable),
	}
}

func ServiceUnavailableMsg(msg any) (int, any) {
	return http.StatusServiceUnavailable, map[string]any{
		"error":   http.StatusText(http.StatusServiceUnavailable),
		"code":    http.StatusText(http.StatusServiceUnavailable),
		"message": msg,
	}
}

func ResponseXml(field, val string) (int, any) {
	return http.StatusOK, gin.H{field: val}
}

func Created(data any) (int, any) {
	result := map[string]any{
		"code":    http.StatusCreated,
		"message": "SUCCESS",
		"data":    data,
	}

	return http.StatusCreated, result
}

func Pagination(data, total, limit, offset any) (int, any) {
	return http.StatusOK, map[string]any{
		"data":   data,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}
}

func OK(data any) (int, any) {
	return http.StatusOK, data
}
