package common

import (
	"bufio"
	"context"
	"io"
	"os"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/storage"
)

func UploadImageToStorage(ctx context.Context, filePath string) (bool, error) {
	var fileBytes []byte
	stat, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}

	fileReader, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Error(err)
		}
	}(fileReader.Name())
	defer func(fileReader *os.File) {
		err = fileReader.Close()
		if err != nil {
			log.Error(err)
		}
	}(fileReader)

	// Read the file into a byte slice
	fileBytes = make([]byte, stat.Size())
	bufferReader := bufio.NewReader(fileReader)
	_, err = bufferReader.Read(fileBytes)
	if err != nil && err != io.EOF {
		return false, err
	}

	// Store file to storage
	metaData := storage.NewStoreInput(fileBytes, filePath)
	isSuccess, err := storage.Instance.Store(ctx, *metaData)
	if err != nil || !isSuccess {
		return false, err
	}

	return true, nil
}
