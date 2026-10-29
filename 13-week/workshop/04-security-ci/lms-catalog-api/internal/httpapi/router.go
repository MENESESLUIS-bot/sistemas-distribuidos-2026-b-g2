// Package httpapi wires the /api/v1/books routes of lms-catalog-api with the
// per-route scopes proposed in the [security] story for STRIDE row E.
package httpapi

import (
	"crypto/rsa"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MENESESLUIS-bot/w13/lms-catalog-api/internal/httpapi/middleware"
)

// Handlers are the book use cases; the real service backs them with
// catalog_db, the test backs them with stubs.
type Handlers struct {
	Create, Update, Delete, LoanCopy, ReturnCopy http.HandlerFunc
}

func NewRouter(publicKey *rsa.PublicKey, internalSecret string, h Handlers) http.Handler {
	r := chi.NewRouter()
	r.Route("/api/v1/books", func(books chi.Router) {
		books.Use(middleware.RequireAuth(publicKey, internalSecret))

		// Administrator only (HU-04, HU-09).
		books.Group(func(admin chi.Router) {
			admin.Use(middleware.RequireScope(middleware.ScopeBooksWrite))
			admin.Post("/", h.Create)
			admin.Patch("/{id}", h.Update)
			admin.Delete("/{id}", h.Delete)
		})

		// lms-circulation-api only (HU-06 loan-copy, HU-07 return-copy).
		books.Group(func(internal chi.Router) {
			internal.Use(middleware.RequireScope(middleware.ScopeCopiesWrite))
			internal.Post("/{id}/loan-copy", h.LoanCopy)
			internal.Post("/{id}/return-copy", h.ReturnCopy)
		})
	})
	return r
}
