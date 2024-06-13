package service

import (
	"context"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/storage"
	"io"
	"mime/multipart"
	"net/url"
	"path"
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

func removeFileFromStorageShareInfo(c context.Context, fileUrl string) error {
	fileName, err := splitFileNameFromUrl(fileUrl)
	if err != nil {
		return err
	}

	input := storage.NewRetrieveInput(fileName)
	return storage.Instance.RemoveFile(c, *input)
}

func splitFileNameFromUrl(fileUrl string) (string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(fileUrl)
	if err != nil {
		return "", err
	}
	fileName := path.Base(parsedURL.Path)
	return fileName, nil
}
