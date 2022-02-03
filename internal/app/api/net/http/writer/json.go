package writer

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, v interface{}) {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	w.Write(bytes)
}
