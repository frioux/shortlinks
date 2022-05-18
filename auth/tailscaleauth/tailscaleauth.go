package tailscaleauth

import (
	"net/http"

	"tailscale.com/client/tailscale"
)

type Auther struct{}

func (a Auther) Wrap(inner http.Handler) http.Handler {
	return inner
}

func (a Auther) User(r *http.Request) (string, error) {
	u, err := tailscale.WhoIs(r.Context(), r.RemoteAddr)
	if err != nil {
		return "", err
	}

	return u.UserProfile.DisplayName, nil
}
