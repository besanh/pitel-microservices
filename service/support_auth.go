package service

// func SupportAuth(ctx context.Context, db sqlclient.ISqlClientConn, source string) (*model.AuthUser, error) {
// 	result := &model.AuthUser{}
// 	filter := model.AuthSourceFilter{
// 		Source: source,
// 		Status: sql.NullBool{
// 			Valid: true,
// 			Bool:  true,
// 		},
// 	}
// 	total, authSource, err := repository.AuthSourceRepo.GetAuthSource(ctx, db, filter, 1, 0)
// 	if err != nil {
// 		log.Error(err)
// 		return nil, err
// 	} else if total < 1 {
// 		return nil, errors.New("source not found")
// 	}

// }
