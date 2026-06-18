package main

import (
	//"fmt"
	//"log"
	"net/http"
	//"os"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	//"github.com/go-chi/docgen"
	"github.com/go-chi/render"
	//"github.com/joho/godotenv"
	"github.com/bettaburger/dSSGoServer/internal/handlers"
)

func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := chi.NewRouter()

	// standard middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// endpoints
	r.Get("/", func(w http.ResponseWriter, r*http.Request) {w.Write([]byte("This is the devskills root server"))})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { // health endpoint, for potential cloud run service
		w.WriteHeader(http.StatusOK) 
		w.Write([]byte("ok"))
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {w.Write([]byte("pong"))})
	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {panic("test")})

	// REST routes 
	r.Route("/api", func(r chi.Router) {
		//r.With(paginate).Get("/", handlers.GetJobData)
		r.Get("/jobs", handlers.GetJobData) // GET /api/jobs
		// add more methods later
	})
	
	http.ListenAndServe(":8000", r) 
}