package domain

import (
	"encoding/json"
)

type HTTPError struct {
	Err string `json:"error"`
}

func NewHTTPError(err error) string {
	httpErr := HTTPError{
		Err: err.Error(),
	}

	httpErrorBytes, _ := json.MarshalIndent(&httpErr, "", "    ")
	return string(httpErrorBytes)
}
