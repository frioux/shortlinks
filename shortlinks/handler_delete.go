package shortlinks

import (
	"net/http"
)

func deleteHandler(db DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				_500(w, err)
				return
			}

			if err := db.DeleteShortlink(r.Form.Get("from")); err != nil {
				_500(w, err)
				return
			}
		}

		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Location", "/")
		w.WriteHeader(303)
	})
}

