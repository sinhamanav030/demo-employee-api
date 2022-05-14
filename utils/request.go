package utils

import (
	"encoding/json"
	"net/http"
)

func SanitizeRequest(r *http.Request, req interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(req); err != nil {
		return err
	}

	// if validationResult := ValidateStruct(req); !validationResult.Success {
	// 	return validationResult.OriginalError
	// }

	return nil
}
