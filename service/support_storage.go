package service

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/storage"
)

func uploadImageToStorageShareInfo(c context.Context, file *multipart.FileHeader) (url string, err error) {
	f, err := file.Open()
	if err != nil {
		log.Error(err)
		return
	}
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		log.Error(err)
		return
	}
	metaData := storage.NewStoreInput(fileBytes, file.Filename)
	isSuccess, err := storage.Instance.Store(c, *metaData)
	if err != nil || !isSuccess {
		log.Error(err)
		return
	}

	input := storage.NewRetrieveInput(file.Filename)
	_, err = storage.Instance.Retrieve(c, *input)
	if err != nil {
		log.Error(err)
		return
	}

	url = API_DOC + "/bss-message/v1/share-info/image/" + input.Path

	return
}
