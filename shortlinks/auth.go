package shortlinks

import (
	"net/http"
)

type Auth interface {
	// Wrap allows the Auth driver to inject middleware to set up
	// authentication.
	Wrap(http.Handler) http.Handler

	// User extracts the user from the http.Request.
	User(*http.Request) (string, error)
}
