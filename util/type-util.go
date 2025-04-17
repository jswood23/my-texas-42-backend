package util

import "encoding/json"

func ConvertStringMapToType[T any](m map[string]interface{}) (T, error) {
	var t T
	b, err := json.Marshal(m)
	if err != nil {
		return t, err
	}
	err = json.Unmarshal(b, &t)
	if err != nil {
		return t, err
	}
	return t, nil
}
