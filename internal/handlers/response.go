package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Datajoblake response
type Response struct {
	Found int `json:"found"`
	Page int `json:"page"`
	PerPage int `json:"per_page"`
	Jobs []Job `json:"jobs"`
}

// job data general
type Job struct {
	//Country string `json:"file_country"`
 //Month string `json:"file_month"`
  //Year string `json:"file_year"`
	ID string `json:"id"` 
	Title string `json:"title"`
	Company string `json:"company"`
	JobFunction string `json:"job_function"`
	Role string `json:"role"`
	Seniority []string `json:"seniority"`
	Location []string `json:"locations"`
	Skills []string `json:"skills"`
}

// skill trends struct 
type SkillTrend struct {
	Skill string `json:"skill"`
	Count int `json:"count"`
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