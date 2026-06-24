package main

// This package starts the API server 

import (
	"fmt"
	"log"
	//"time"
	"net/http"
	"os"
	"flag"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	//"github.com/go-chi/docgen"
	"github.com/go-chi/render"
	//"github.com/joho/godotenv"
	"github.com/go-chi/docgen"
	"github.com/bettaburger/dSSGoServer/internal/handlers"
	"github.com/bettaburger/dSSGoServer/internal/cloudsql"

)

var routes = flag.Bool("routes", false, "Generate router documentation")

func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

func main() {
	/*err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}*/
	flag.Parse()
	// standard middleware
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	//r.Use(middleware.Timeout(120 * time.Second)) // timeout value of 2 minute, could mess w container deployment 

	// endpoints
	r.Get("/", func(w http.ResponseWriter, r*http.Request) {w.Write([]byte("This is the devskills root server"))})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { // health endpoint, for cloud run service
		w.WriteHeader(http.StatusOK) 
		w.Write([]byte("ok"))
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {w.Write([]byte("pong"))})
	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {panic("test")})

	if *routes {
		fmt.Println(docgen.JSONRoutesDoc(r))
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "github.com/bettaburger/dSSGoServer",
		}))
		return
	}

	ctx := context.Background()
	// Connect to Cloud sql 
	db, cleanup, err := cloudsql.ConnectDB(ctx)
	if err != nil {
		log.Printf("DB connection failed: %v", err)
	}
	defer cleanup()
	defer db.Close()

	handler := &handlers.TaskHandler{DB : db}

	// REST routes /api
	r.Route("/api", func(r chi.Router) {
		//r.With(paginate).Get("/", handlers.GetJobData)
		r.Get("/jobs", handler.GetJobData) // GET /api/jobs
		// add more methods later
	})

	//r.Mount("/admin", adminRouter())

	// cloudsql runs on default port 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on :%s", port)
    if err := http.ListenAndServe(":"+port, r); err != nil {
      log.Printf("Server failed: %v", err)
  }
}
