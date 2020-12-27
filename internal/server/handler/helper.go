package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func MapBodyToStruct(r *http.Request, o interface{}) error {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		return errors.New("content-type is not application/json")
	}

	var unmarshalErr *json.UnmarshalTypeError
	decoder := json.NewDecoder(r.Body)
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
