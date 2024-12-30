package emails

import "github.com/go-chi/chi"

func EmailsRoutes() chi.Router {
	r := chi.NewRouter()
	email := Email{}

	r.Get("/", email.GetAllEmails)
	r.Get("/search", email.SearchEmails)

	return r
}
