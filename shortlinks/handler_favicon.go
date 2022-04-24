package shortlinks

import (
	"fmt"
	"net/http"
)

func faviconHandler(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "image/svg+xml")
	rw.Header().Set("Cache-Control", "Cache-Control: public, max-age=604800, immutable")
	fmt.Fprintln(rw, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text y=".9em" font-size="90">⚡</text></svg>`)
}
