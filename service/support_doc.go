package service

import (
	"mime/multipart"

	"github.com/tel4vn/fins-microservices/model"
)

func (s *Message) UploadDoc(authUser *model.AuthUser, data model.MessageRequest, file *multipart.FileHeader) (fileUrl string, err error) {
	return
}
