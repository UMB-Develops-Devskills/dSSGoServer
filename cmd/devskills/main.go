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
	"github.com/UMB-Develops-Devskills/dSSGoServer/internal/handlers"
	"github.com/UMB-Develops-Devskills/dSSGoServer/internal/firebase"

)

var routes = flag.Bool("routes", false, "Generate router documentation")

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

	// generate a routes.json file
	if *routes {
		fmt.Println(docgen.JSONRoutesDoc(r))
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "github.com/bettaburger/dSSGoServer",
		}))
		return
	}
	
	// Connect to Firebase
	ctx := context.Background()
	client, err := firebase.ConnectFireBaseStorage(ctx)
	if err != nil {
		log.Printf("firebase connection failed: %v", err)
	}
	defer client.Close()

	handler := &handlers.TaskHandler{Storage : client}

	// REST routes /api
	r.Route("/api", func(r chi.Router) {
		//r.With(paginate).Get("/", handlers.GetJobData)
		//r.Get("/jobs", handler.GetJobData) // GET /api/jobs
		r.Get("/keys", handler.GetKeyData) // GET /api/keys, call this in frontend to display the filters
		r.Route("/trends", func(r chi.Router) {
			r.Get("/skills", handler.GetSkillTrends) // GET /api/trends/skills
			r.Get("/graph", handler.GetSkillGraph) // GET /api/trends/graph
			r.Get("/history", handler.GetHistoryData) // GET /api/trends/graph
		})
	})

	//r.Mount("/admin", adminRouter())

	// devskills instance runs on default port localhost:8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on :%s", port)
    if err := http.ListenAndServe(":"+port, r); err != nil {
      log.Printf("Server failed: %v", err)
  }
}
