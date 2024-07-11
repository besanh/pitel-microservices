package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IChatVendor interface {
		IRepo[model.ChatVendor]
		GetVendors(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatVendorFilter, limit, offset int) (total int, result *[]model.ChatVendor, err error)
	}
	Vendor struct {
		Repo[model.ChatVendor]
	}
)

var VendorRepo IChatVendor

func NewVendor() IChatVendor {
	return &Vendor{}
}

func (repo *Vendor) GetVendors(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatVendorFilter, limit, offset int) (total int, result *[]model.ChatVendor, err error) {
	result = new([]model.ChatVendor)
	query := db.GetDB().NewSelect().
		Model(result)
	if len(filter.VendorName) > 0 {
		query.Where("vendor_name = ?", filter.VendorName)
	}
	if len(filter.VendorType) > 0 {
		query.Where("vendor_type = ?", filter.VendorType)
	}
	if filter.Status.Valid {
		query.Where("status = ?", filter.Status.Bool)
	}
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	total, err = query.ScanAndCount(ctx)
	if err == sql.ErrNoRows {
		return 0, result, nil
	} else if err != nil {
		return 0, result, err
	}
	return total, result, nil
}
