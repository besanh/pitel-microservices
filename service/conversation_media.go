package service

import (
	"context"
	"net/url"
	"path"
	"regexp"
	"time"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IConversationMedia interface {
		GetConversationMedias(ctx context.Context, authUser *model.AuthUser, filter model.ConversationMediaFilter, limit, offset int) (int, *[]model.ConversationMedia, error)
		InsertConversationMedias(ctx context.Context, authUser *model.AuthUser, message model.Message) error
	}
	ConversationMedia struct{}
)

var ConversationMediaService IConversationMedia

func NewConversationMedia() IConversationMedia {
	return &ConversationMedia{}
}

func (s *ConversationMedia) GetConversationMedias(ctx context.Context, authUser *model.AuthUser, filter model.ConversationMediaFilter, limit, offset int) (total int, medias *[]model.ConversationMedia, err error) {
	filter.TenantId = authUser.TenantId
	total, medias, err = repository.ConversationMediaRepo.GetConversationMedias(ctx, repository.DBConn, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func (s *ConversationMedia) InsertConversationMedias(ctx context.Context, authUser *model.AuthUser, message model.Message) (err error) {
	entities := make([]model.ConversationMedia, 0)
	for _, item := range message.Attachments {
		if item == nil {
			log.Error("not found attachment in message " + message.MessageId)
			continue
		}
		data := model.ConversationMediaRequest{
			TenantId:               authUser.TenantId,
			ConversationId:         message.ConversationId,
			ExternalConversationId: message.ExternalConversationId,
			ConversationType:       message.MessageType,
			MessageId:              message.MessageId,
			MediaType:              item.AttType,
			SendTimestamp:          time.UnixMilli(message.SendTimestamp),
		}
		if item.Payload != nil {
			data.MediaUrl = item.Payload.Url
			data.MediaHeader = item.Payload.Name
			data.MediaSize = item.Payload.Size

			if len(data.MediaHeader) < 1 {
				mediaUrl, err := url.Parse(data.MediaUrl)
				if err != nil {
					log.Error(err)
					return err
				}
				data.MediaHeader = path.Base(mediaUrl.Path)
			}
		}
		if err = data.Validate(); err != nil {
			log.Error(err)
			return
		}
		tmp := model.ConversationMedia{
			Base:                   model.InitBase(),
			TenantId:               data.TenantId,
			ConversationId:         data.ConversationId,
			ExternalConversationId: data.ExternalConversationId,
			ConversationType:       data.ConversationType,
			MessageId:              data.MessageId,
			MediaType:              data.MediaType,
			MediaHeader:            data.MediaHeader,
			MediaUrl:               data.MediaUrl,
			MediaSize:              data.MediaSize,
			SendTimestamp:          data.SendTimestamp,
		}

		entities = append(entities, tmp)
	}

	// Find all matches in the text
	data := model.ConversationMediaRequest{
		TenantId:               authUser.TenantId,
		ConversationId:         message.ConversationId,
		ExternalConversationId: message.ExternalConversationId,
		ConversationType:       message.MessageType,
		MessageId:              message.MessageId,
		MediaType:              string(model.MediaLink),
		SendTimestamp:          time.UnixMilli(message.SendTimestamp),
	}
	re := regexp.MustCompile(`https?://[^\s]+`)
	matches := re.FindAllString(message.Content, -1)
	for _, embeddedLink := range matches {
		data.MediaUrl = embeddedLink
		if err = data.Validate(); err != nil {
			log.Error(err)
			return
		}

		tmp := model.ConversationMedia{
			Base:                   model.InitBase(),
			TenantId:               data.TenantId,
			ConversationId:         data.ConversationId,
			ExternalConversationId: data.ExternalConversationId,
			ConversationType:       data.ConversationType,
			MessageId:              data.MessageId,
			MediaType:              data.MediaType,
			MediaHeader:            data.MediaHeader,
			MediaUrl:               data.MediaUrl,
			MediaSize:              data.MediaSize,
			SendTimestamp:          data.SendTimestamp,
		}

		entities = append(entities, tmp)
	}

	tx, err := repository.ConversationMediaRepo.BeginTx(ctx, repository.DBConn, nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer tx.Rollback()

	if len(entities) > 0 {
		if err = repository.ConversationMediaRepo.TxBulkInsert(ctx, tx, entities); err != nil {
			log.Error(err)
			return
		}
	}

	if err = repository.ConversationMediaRepo.CommitTx(ctx, tx); err != nil {
		log.Error(err)
		return
	}

	return
}
