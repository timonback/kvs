package util

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

func MapBodyToStruct(body io.Reader, header http.Header, o interface{}) error {
	headerContentTtype := header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		return errors.New("content-type is not application/json")
	}

	var unmarshalErr *json.UnmarshalTypeError
	decoder := json.NewDecoder(body)
	err := decoder.Decode(o)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			return errors.New("Wrong Type provided for field " + unmarshalErr.Field + " at position " + strconv.FormatInt(unmarshalErr.Offset, 10))
		} else {
			return errors.New("internal error during request parsing")
		}
	}
	return nil
}
