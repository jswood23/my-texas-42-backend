package models

import (
	"encoding/json"
	"io"
)

func DecodeAPIModel[T any](body io.ReadCloser) (T, error) {
	var model T
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return model, err
	}
	defer body.Close()

	err = json.Unmarshal(bodyBytes, &model)
	if err != nil {
		return model, err
	}

	return model, nil
}
