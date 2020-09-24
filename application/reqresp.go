package application

import (
	"encoding/json"
	"net/http"
)

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
