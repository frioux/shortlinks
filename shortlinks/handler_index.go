package shortlinks

import (
	"github.com/hbollon/go-edlib"
	"net/http"
	"sort"
	"strings"
)

type index struct {
	Shortlinks []Shortlink
}

type search struct {
	Path       string
	Shortlinks []Shortlink
}

func (i index) Title() string  { return "go links" }
func (i index) To() string     { return "" }
func (i index) From() string   { return "" }
func (i index) Submit() string { return "Create" }
func (i index) Description() string { return "" }

func (s search) Title() string  { return "go links" }
func (s search) To() string     { return "" }
func (s search) From() string   { return "" }
func (s search) Submit() string { return "Create" }
func (s search) Description() string { return "" }

type scoredShortlink struct {
	shortlink Shortlink
	score     int
}
type scoredShortLinksWrapper []scoredShortlink

func (s scoredShortLinksWrapper) Len() int           { return len(s) }
func (s scoredShortLinksWrapper) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s scoredShortLinksWrapper) Less(i, j int) bool { return s[i].score < s[j].score }

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
func possibleMatches(shortlinks []Shortlink, path string, resultSize int) []Shortlink {
	scoredShortlinks := make([]scoredShortlink, len(shortlinks))

	for i, sl := range shortlinks {
		score := edlib.DamerauLevenshteinDistance(path, sl.From)
		scoredShortlinks[i] = scoredShortlink{shortlink: sl, score: score}
	}

	sort.Sort(scoredShortLinksWrapper(scoredShortlinks))

	for i, sl := range scoredShortlinks {
		shortlinks[i] = sl.shortlink
	}

	return shortlinks[:min(resultSize, len(shortlinks))]
}

func indexHandler(db PublicDB) http.Handler {
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
		} else {
			path := strings.TrimPrefix(r.URL.Path, "/")
			sl, err := db.Shortlink(path)
			if err != nil {
				_500(w, err)
				return
			}

			emptyShortLink := Shortlink{}
			if sl == emptyShortLink {
				sls, err := db.AllShortlinks()
				if err != nil {
					_500(w, err)
					return
				}
				v := search{Path: path, Shortlinks: possibleMatches(sls, path, 20)}

				w.WriteHeader(404)
				if err := tpl.ExecuteTemplate(w, "search.html", v); err != nil {
					_500(w, err)
					return
				}
				return
			} else {
				w.Header().Add("Location", sl.To)
				w.WriteHeader(302)
			}
		}
	})
}
