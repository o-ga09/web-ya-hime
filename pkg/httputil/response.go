package httputil

import (
	"encoding/json"
	"net/http"
)

func Response(w *http.ResponseWriter, status int, message ...interface{}) {
	if len(message) == 0 {
		(*w).WriteHeader(status)
		return
	}

	var data interface{}
	if len(message) == 1 {
		data = message[0]
	} else {
		data = message
	}

	json, err := json.Marshal(data)
	if err != nil {
		(*w).Header().Set("Content-Type", "application/json")
		(*w).WriteHeader(http.StatusInternalServerError)
		(*w).Write([]byte(`{"message": "error marshalling json"}`))
		return
	}

	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(status)
	(*w).Write(json)
}
