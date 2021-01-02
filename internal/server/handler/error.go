package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorParameter struct {
	key   string
	value interface{}
}

const (
	ParameterPath = "path"
)

func CreateErrorParameters(parameters ...ErrorParameter) map[string]interface{} {
	params := make(map[string]interface{}, len(parameters))
	for _, el := range parameters {
		params[el.key] = el.value
	}
	return params
}

func HandleError(w http.ResponseWriter, r *http.Request, statusCode int, e error, errorParams map[string]interface{}) {
	writeJsonHeaders(w)
	w.WriteHeader(statusCode)

	resp := make(map[string]interface{})
	resp["message"] = e.Error()
	resp["statusCode"] = statusCode
	resp["errorParameters"] = errorParams
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
