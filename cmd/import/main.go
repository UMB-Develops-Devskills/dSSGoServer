package main

import (
	"context"
	"encoding/json"
	"os"
	"log"
	"github.com/bettaburger/dSSGoServer/internal/cloudsql"
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
	ID string `json:"id"` 
	Title string `json:"title"`
	Company string `json:"company"`
	JobFunction string `json:"job_function"`
	Role string `json:"role"`
	Location []string `json:"locations"`
	Skills []string `json:"skills"`
}

func main() {
	filename := os.Args[1]

	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("failed to read the file %v", err)
	}
	var resp Response 
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Printf("json: %v", err)
	}
	ctx := context.Background()
	// Connect to Cloud sql 
	db, cleanup, err := cloudsql.ConnectDB(ctx)
	if err != nil {
		log.Printf("DB connection failed: %v", err)
	}
	defer cleanup()
	defer db.Close()

	i := 0
	for _, j := range resp.Jobs {
		_, err := db.Exec(ctx, `INSERT INTO jobs (id, title, company, job_function, role, locations, skills) VALUES ($1,$2,$3,$4,$5,$6,$7)`, j.ID, j.Title, j.Company, j.JobFunction, j.Role, j.Location, j.Skills)
		if err != nil {
			log.Printf("failed importing script %s: %v", j.ID, err)
			continue
		}
		i++
	}
	log.Printf("%d of %d jobs", i, len(resp.Jobs))
}
