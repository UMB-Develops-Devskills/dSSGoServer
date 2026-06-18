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

// job data general
type Job struct {
	ID uint16 `json:"id"`
	Title string `json:"Title"`
	Role string `json:"role"`
	Location string `json:"location"`
	Time string `json:"time"`
}

// This struct holds dependencies for http handlers
type TaskHandler struct {
	DB *pgxpool.Pool
}

// This function returns all the jobs from the database 
func (h *TaskHandler) GetJobData(w http.ResponseWriter, r *http.Request) {
	/*err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}*/
	//dataLakeKey := os.Getenv("DATA_LAKE_API")
	rows, err := h.DB.Query(r.Context(), "SELECT id, title, role, location, time FROM jobs")
	if err != nil {
		http.Error(w, "failed to query jobs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	// jobs array
	var jobs []Job
	for rows.Next() {
		var j Job
		if err := rows.Scan(&j.ID, &j.Title, &j.Role, &j.Location, &j.Time); err != nil {
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