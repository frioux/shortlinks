package shortlinks

import (
	"fmt"
	"net/http"
	"os"
)

type Server struct {
	DB DB
}

func (s Server) ListenAndServe(listen string) error {
	http.Handle("/", indexHandler(s.DB))
	http.Handle("/_edit/", editHandler(s.DB))
	http.Handle("/_favicon", http.HandlerFunc(faviconHandler))

	fmt.Fprintln(os.Stderr, "serving at", listen)
	return http.ListenAndServe(listen, nil)
}

func _500(w http.ResponseWriter, err error) {
	fmt.Fprintln(os.Stderr, err)
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(500)
	fmt.Fprintln(w, err)
}
