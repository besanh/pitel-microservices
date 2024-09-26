package service

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"

	"github.com/tel4vn/fins-microservices/common/constant"
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

func uploadFileToStorage(c context.Context, buffer *bytes.Buffer, prefix, fileName string) (url string, err error) {
	metaData := storage.NewStoreInput(buffer.Bytes(), fileName)
	isSuccess, err := storage.Instance.PresignedStore(c, *metaData, constant.OBJECT_EXPIRE_TIME)
	if err != nil || !isSuccess {
		log.Error(err)
		return
	}

	input := storage.NewRetrieveInput(fileName)
	url, err = storage.Instance.PresignedRetrieve(c, *input, constant.OBJECT_EXPIRE_TIME)
	if err != nil {
		log.Error(err)
		return
	}

	url = API_DOC + prefix + input.Path
	return
}
