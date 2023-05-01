package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

// The error response the handlers will use
type errorResponse struct {
	HttpStatusCode int    `json:"-"`
	Status         string `json:"status"`
	Error          string `json:"error"`
}

// Implements the render.Renderer so we can easily return this as a JSON
func (resp *errorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Will set the response status based on the error
	render.Status(r, resp.HttpStatusCode)
	return nil
}

func newErrorResponse(err error, status int) render.Renderer {
	return &errorResponse{
		HttpStatusCode: status,
		Status:         http.StatusText(status),
		Error:          err.Error(),
	}
}
