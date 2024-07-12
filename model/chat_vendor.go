package model

import (
	"errors"
	"mime/multipart"
	"slices"

	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/uptrace/bun"
)

type ChatVendor struct {
	*Base
	bun.BaseModel `bun:"table:chat_vendor,alias:cv"`
	VendorName    string `json:"vendor_name" bun:"vendor_name,type:text,notnull"`
	VendorType    string `json:"vendor_type" bun:"vendor_type,type:text,notnull"`
	Logo          string `json:"logo" bun:"logo,type:text,notnull"`
	Status        bool   `json:"status" bun:"status,type:boolean,nullzero,default:false"`
}

type ChatVendorRequest struct {
	VendorName string                `json:"vendor_name" form:"vendor_name""`
	VendorType string                `json:"vendor_type" form:"vendor_type"`
	Logo       string                `json:"logo"`
	Status     bool                  `json:"status" form:"status"`
	File       *multipart.FileHeader `json:"file" form:"file"`
}

func (m *ChatVendorRequest) Validate() error {
	if len(m.VendorName) < 1 {
		return errors.New("vendor name is required")
	}

	if len(m.VendorType) < 1 {
		return errors.New("vendor type is required")
	}

	if !slices.Contains(variables.VENDOR_TYPES, m.VendorType) {
		return errors.New("vendor type " + m.VendorType + " is not supported")
	}

	return nil
}
