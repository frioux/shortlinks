package shortlinks

import (
	"net/http"
	"strings"
	"fmt"
	"os"
)

func publicIndexHandler(db PublicDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			sl, err := db.AllShortlinks()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				w.Header().Add("Content-Type", "text/plain")
				w.WriteHeader(500)
				fmt.Fprintln(w, "couldn't load links")
				return
			}
			v := index{Shortlinks: sl}

			if err := tpl.ExecuteTemplate(w, "public_index.html", v); err != nil {
				fmt.Fprintln(os.Stderr, err)
				w.Header().Add("Content-Type", "text/plain")
				w.WriteHeader(500)
				fmt.Fprintln(w, "couldn't execute template")
				return
			}
			return
		}

		sl, err := db.Shortlink(strings.TrimPrefix(r.URL.Path, "/"))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(500)
			fmt.Fprintln(w, "couldn't load link")
			return
		}
		w.Header().Add("Location", sl.To)
		w.WriteHeader(302)

	})
}
