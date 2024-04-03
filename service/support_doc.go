package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
)

func (s *Message) UploadDoc(ctx context.Context, file *multipart.FileHeader) (fileUrl string, err error) {
	fileContent, err := file.Open()
	if err != nil {
		log.Error(err)
		return "", err
	}
	byteContent := make([]byte, file.Size)
	fileContent.Read(byteContent)

	url := OTT_URL + "/ott/" + OTT_VERSION + "/crm/upload"
	client := resty.New()
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetFileReader("data", file.Filename, bytes.NewReader(byteContent)).
		Post(url)
	if err != nil {
		log.Error(err)
		return "", err
	}
	if res.StatusCode() == 200 {
		var result model.OttUploadResponse
		if err := json.Unmarshal([]byte(res.Body()), &result); err != nil {
			log.Error(err)
			return "", err
		}
		if len(result.Data) > 0 {
			return result.Data[0], nil
		}
		return "", nil
	}

	return "", errors.New("upload failed")
}
