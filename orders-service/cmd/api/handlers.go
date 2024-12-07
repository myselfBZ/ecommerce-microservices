package main

import "net/http"

var successMessage = map[string]bool{
	"success": true,
}

type apiResponse struct {
	Err    string `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
	status int
}

type apiHandler func(http.ResponseWriter, *http.Request) *apiResponse

func makeHTTPHandler(f apiHandler, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		apiResp := f(w, r)
		if apiResp.Err != "" {
			writeJSON(w, apiResp, apiResp.status)
			return
		}
		writeJSON(w, apiResp, apiResp.status)
	}
}
