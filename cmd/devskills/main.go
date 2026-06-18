package main

import (
	//"fmt"
	"log"
	"net/http"
	"os"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	//"github.com/go-chi/docgen"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
	"github.com/bettaburger/dSSGoServer/internal/handlers"
	"github.com/bettaburger/dSSGoServer/internal/cloudsql"

)

func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()
	
	// Connect to Cloud sql 
	db, cleanup, err := cloudsql.ConnectDB(ctx)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer cleanup()
	defer db.Close()

	handler := &handlers.TaskHandler{DB : db}

	// standard middleware
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// endpoints
	r.Get("/", func(w http.ResponseWriter, r*http.Request) {w.Write([]byte("This is the devskills root server"))})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { // health endpoint, for cloud run service
		w.WriteHeader(http.StatusOK) 
		w.Write([]byte("ok"))
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {w.Write([]byte("pong"))})
	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {panic("test")})

	// REST routes /api
	r.Route("/api", func(r chi.Router) {
		//r.With(paginate).Get("/", handlers.GetJobData)
		r.Get("/jobs", handler.GetJobData) // GET /api/jobs
		// add more methods later
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	http.ListenAndServe(":"+port, r)
}