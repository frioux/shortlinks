package shortlinks

import (
	"net/http"
)

type edit struct {
	Shortlink
	Submit string

	History []History
}

func (e edit) Title() string {
	if e.From == "" {
		return "Create"
	}

	return "Edit " + e.From
}

func editHandler(db DB, auth Auth) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		from := r.URL.Query().Get("from")

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

			if from == "" {
				from = r.Form.Get("from")
			}

			if err := db.InsertHistory(History{
				From: from,
				To:   r.Form.Get("to"),
				Who:  u,

				Description: r.Form.Get("description"),
			}); err != nil {
				_500(w, err)
				return
			}
			if err := db.CreateShortlink(Shortlink{
				To:   r.Form.Get("to"),
				From: from,

				Description: r.Form.Get("description"),
			}); err != nil {
				_500(w, err)
				return
			}
			w.Header().Add("Location", "/")
			w.WriteHeader(302)
			return
		}

		sl, err := db.Shortlink(from)
		if err != nil {
			_500(w, err)
			return
		}

		h, err := db.History(from)
		if err != nil {
			_500(w, err)
			return
		}

		v := edit{
			Shortlink: sl,
			History:   h,

			Submit: "Update",
		}

		if err := tpl.ExecuteTemplate(w, "edit.html", v); err != nil {
			_500(w, err)
			return
		}
	})
}
