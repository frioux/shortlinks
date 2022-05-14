package shortlinks

import (
	"net/http"
)

func deleteHandler(db DB, auth Auth) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var u string
			if auth != nil {
				var err error
				u, err = auth.User(r)
				if err != nil {
					_403(w, err)
					return
				}
			}
			if err := r.ParseForm(); err != nil {
				_500(w, err)
				return
			}

			if err := db.DeleteShortlink(r.Form.Get("from"), u); err != nil {
				_500(w, err)
				return
			}
		}

		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Location", "/")
		w.WriteHeader(303)
	})
}
