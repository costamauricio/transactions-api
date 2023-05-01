package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

// Parse the request body into the provided struct that implements a render.Binder interface
func parseRequestBody(r *http.Request, payload render.Binder) error {
	err := render.Bind(r, payload)
	if err != nil {
		// If the error is from the JSON unmarshal, we set a custom message
		if typeError, isTypeError := err.(*json.UnmarshalTypeError); isTypeError {
			err = errors.New("invalid type " + typeError.Value + " for '" + typeError.Field + "' field")
		}

		return err
	}

	return nil
}
