package feign

import (
	"encoding/json"
	"io"
)

func decodeResponse(body io.Reader, target interface{}) error {
	if target == nil {
		return nil
	}
	return json.NewDecoder(body).Decode(target)
}

