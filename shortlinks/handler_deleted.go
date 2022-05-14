package shortlinks

import (
	"net/http"
)

type deleted struct {
	Shortlinks []Shortlink
}

func (d deleted) Title() string { return "deleted links" }

func deletedHandler(db DBDeleted) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sl, err := db.DeletedShortlinks()
		if err != nil {
			_500(w, err)
			return
		}
		v := deleted{Shortlinks: sl}

		if err := tpl.ExecuteTemplate(w, "deleted.html", v); err != nil {
			_500(w, err)
			return
		}
		return
	})
}
