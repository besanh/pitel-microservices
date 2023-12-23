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
	ERR_TOKEN_IS_EMPTY    = "token is empty"
	ERR_TOKEN_IS_INVALID  = "token is invalid"
	ERR_TOKEN_IS_EXPIRED  = "token is expired"
	ERR_EMPTY_CONN        = "empty connection"
	ERR_EXAMPLE_NOT_FOUND = "campaign not found"
	ERR_EXAMPLE_INVALID   = "campaign is invalid"
	ERR_DATA_NOT_FOUND    = "data not found"
	ERR_DATA_INVALID      = "data is invalid"
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

func ResponseXml(field, val string) (int, any) {
	return http.StatusOK, gin.H{field: val}
}
