package swagger

import (
	_ "embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed swagger.json
var swaggerJSON []byte

//go:embed index.html
var indexHTML []byte

func RegisterRoutes(r chi.Router) {
	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(swaggerJSON)
	})

	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	r.Get("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(indexHTML)
	})
}
