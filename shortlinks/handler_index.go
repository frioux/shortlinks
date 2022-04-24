package shortlinks

import (
	"net/http"
	"strings"
)

type index struct {
	Shortlinks []Shortlink
}

func (i index) Title() string  { return "go links" }
func (i index) To() string     { return "" }
func (i index) From() string   { return "" }
func (i index) Submit() string { return "Create" }

func indexHandler(db DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			sl, err := db.AllShortlinks()
			if err != nil {
				_500(w, err)
				return
			}
			v := index{Shortlinks: sl}

			if err := tpl.ExecuteTemplate(w, "index.html", v); err != nil {
				_500(w, err)
				return
			}
			return
		}

		sl, err := db.Shortlink(strings.TrimPrefix(r.URL.Path, "/"))
		if err != nil {
			_500(w, err)
			return
		}
		w.Header().Add("Location", sl.To)
		w.WriteHeader(302)

	})
}
