package common

import (
	"context"
	"errors"
	"time"

	cacheUtil "github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

const (
	INFO_TEMPLATE   = "info_template"
	EXPIRE_TEMPLATE = 30 * time.Minute
)

func HandleCheckContentMatchTemplate(ctx context.Context, dbCon sqlclient.ISqlClientConn, authUser *model.AuthUser, templateUuid, content string) ([]string, []string, error) {
	templateCache, err := cacheUtil.MCache.Get(INFO_TEMPLATE + "_" + templateUuid)
	if err != nil {
		return nil, nil, err
	} else if templateCache != nil {
		template := templateCache.(*model.TemplateBss)
		keysContent, keysTemplate, ok := handleMatchMessageWithTemplate(content, template.Content)
		if !ok {
			return nil, nil, errors.New("content not match template")
		}
		return keysContent, keysTemplate, nil
	} else {
		template, err := repository.TemplateBssRepo.GetById(ctx, dbCon, templateUuid)
		if err != nil {
			return nil, nil, err
		}
		keysContent, keysTemplate, ok := handleMatchMessageWithTemplate(content, template.Content)
		if !ok {
			return nil, nil, errors.New("content not match template")
		}
		if err := cacheUtil.MCache.SetTTL(INFO_TEMPLATE+"_"+templateUuid, template, EXPIRE_TEMPLATE); err != nil {
			return nil, nil, err
		}
		return keysContent, keysTemplate, nil
	}
}

// Remove conten in {{bracket}} and compare 2 content
// keysContent: {{bracket}} from content message
// keys: {{bracket}} from template
func handleMatchMessageWithTemplate(content string, template string) (keysContent, keys []string, ok bool) {
	contentNew, keysContent, ok := util.CheckTemplate(content, true)
	if ok {
		return keysContent, keys, false
	}
	templateNew, keys, ok := util.CheckTemplate(template, false)
	if ok {
		return keysContent, keys, false
	}
	if keys[len(keys)-1] == "" {
		keys = keys[:len(keys)-1]
	}
	if contentNew == templateNew {
		return keysContent, keys, true
	}
	return keysContent, keys, false
}
