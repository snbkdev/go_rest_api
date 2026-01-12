package helpers

import (
	"context"
	"encoding/json"
	"net/http"
)

func DecodeEncryptRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request EncryptRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeDecryptRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request DecryptRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}