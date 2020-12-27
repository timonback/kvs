package handler

import (
	"encoding/json"
	"net/http"
)

func HandleError(w http.ResponseWriter, r *http.Request, e error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := make(map[string]interface{})
	resp["message"] = e.Error()
	resp["statusCode"] = statusCode
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
