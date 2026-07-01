package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	//"os"
	//"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	//"github.com/joho/godotenv"
)

// Datajoblake response
type Response struct {
	Found int `json:"found"`
	Page int `json:"page"`
	PerPage int `json:"per_page"`
	Jobs []Job `json:"jobs"`
}

// job data general, usa based 
type Job struct {
	//Country string `json:"file_country"`
 //Month string `json:"file_month"`
  //Year string `json:"file_year"`
	ID string `json:"id"` 
	Title string `json:"title"`
	Company string `json:"company"`
	JobFunction string `json:"job_function"`
	Role string `json:"role"`
	Location []string `json:"locations"`
	Skills []string `json:"skills"`
}

// This struct holds dependencies for http handlers
type TaskHandler struct {
	DB *pgxpool.Pool
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// This function requests jobs from the database 
func (h *TaskHandler) GetJobData(w http.ResponseWriter, r *http.Request) {
	/*err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}*/
	role := r.URL.Query().Get("role")
	//title := r.URL.Query().Get("title") // ex. senior, intern...
	//location := r.URL.Query().Get("locations")
	company := r.URL.Query().Get("company")
	country := r.URL.Query().Get("country")
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")
	// run examples: 
	// http://localhost:8080/api/jobs?country=usa&role=Program+Manager&month=june&year=2026
	// http://localhost:8080/api/jobs?country=usa&company=SpaceX
	query := `SELECT id, title, company, job_function, role, locations, skills FROM jobs 
	WHERE ($1::text IS NULL OR country = $1)
  AND ($2::text IS NULL OR company = $2)
  AND ($3::text IS NULL OR role = $3)
  AND ($4::text IS NULL OR month = $4)
  AND ($5::text IS NULL OR year = $5);`
	
	rows, err := h.DB.Query(r.Context(), query, nullIfEmpty(country), nullIfEmpty(company), nullIfEmpty(role), nullIfEmpty(month), nullIfEmpty(year))
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "failed to query jobs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	// jobs array
	var jobs []Job
	for rows.Next() {
		var j Job
		if err := rows.Scan(&j.ID, &j.Title, &j.Company, &j.JobFunction, &j.Role, &j.Location, &j.Skills); err != nil {
			log.Printf("scan error: %v", err)
			http.Error(w, "failed to scan job", http.StatusInternalServerError)
			return
		}
		jobs = append(jobs, j)
	}
	// if no jobs exist return an empty array
	if jobs == nil {
		jobs = []Job{}
	}
	// write to json
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}