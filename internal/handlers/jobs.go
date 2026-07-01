package handlers

import (
	"encoding/json"
	//"log"
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

// This function requests all the jobs from the database 
func (h *TaskHandler) GetJobData(w http.ResponseWriter, r *http.Request) {
	/*err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}*/
	rows, err := h.DB.Query(r.Context(), "SELECT id, title, company, job_function, role, locations, skills FROM jobs")
	if err != nil {
		http.Error(w, "failed to query jobs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	// jobs array
	var jobs []Job
	for rows.Next() {
		var j Job
		if err := rows.Scan(&j.ID, &j.Title, &j.Company, &j.JobFunction, &j.Role, &j.Location, &j.Skills); err != nil {
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

// This function creates a new job query into the db


// This function returns jobs by [specified query]
// This function returns jobs by ID
// This function returns jobs by Role
// This function returns jobs by title 
// This function returns jobs by location