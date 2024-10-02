package routes

import (
	"net/http"

	"github.com/MeherKandukuri/studioClasses_API/handlers"
	"github.com/go-chi/chi"
)

func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Post("/classes", handlers.PostCreateClass)
	mux.Post("/bookings", handlers.PostCreateBooking)

	return mux
}
