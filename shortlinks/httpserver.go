package shortlinks

import (
	"fmt"
	"net/http"
	"os"
)

type Server struct {
	DB   DB
	Auth Auth
}

func (s Server) ListenAndServe(listen string) error {
	mux := http.NewServeMux()

	mux.Handle("/", indexHandler(s.DB))
	mux.Handle("/_delete/", deleteHandler(s.DB, s.Auth))
	mux.Handle("/_edit/", editHandler(s.DB, s.Auth))
	mux.Handle("/_favicon", http.HandlerFunc(faviconHandler))

	if dbd, ok := s.DB.(DBDeleted); ok {
		mux.Handle("/_deleted/", deletedHandler(dbd))
	}

	var h http.Handler = mux
	if auth := s.Auth; auth != nil {
		h = auth.Wrap(h)
	}

	fmt.Fprintln(os.Stderr, "rw serving at", listen)
	return http.ListenAndServe(listen, h)
}

func (s Server) PublicListenAndServe(listen string) error {
	mux := http.NewServeMux()

	mux.Handle("/", publicIndexHandler(s.DB))
	mux.Handle("/_favicon", http.HandlerFunc(faviconHandler))

	fmt.Fprintln(os.Stderr, "public serving at", listen)
	return http.ListenAndServe(listen, mux)
}

func _500(w http.ResponseWriter, err error) {
	fmt.Fprintln(os.Stderr, err)
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(500)
	fmt.Fprintln(w, err)
}

func _403(w http.ResponseWriter, err error) {
	fmt.Fprintln(os.Stderr, err)
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(403)
	fmt.Fprintln(w, "forbidden")
}
