package codec

import (
	"encoding/base64"
	"errors"
)

func Encode(token string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(token))
}

func Decode(encoded string) (string, error) {
	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", errors.New("invalid base64-encoded token")
	}
	return string(data), nil
}
