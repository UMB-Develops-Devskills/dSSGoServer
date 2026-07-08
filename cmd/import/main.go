package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"github.com/UMB-Develops-Devskills/dSSGoServer/internal/cloudsql"
)

/* filename structure
EX: jobs-default-usa-june-2026-page1..n.json
See /data for files
*/
type DataFile struct {
	Country string
	Month string	// what month the data was collected
	Year string	// what year the data was collected
}

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
	Company string `json:"company_name"`
	JobFunction string `json:"job_function"`
	Role string `json:"role"`
	Seniority []string `json:"seniority"`
	Location []string `json:"locations"`
	Skills []string `json:"required_skills"`
}

// this function parses the filepath, retrieves the file and returns Datafile{}
func ParseDataFile(filename string) (DataFile, error) {
	base := filepath.Base(filename)
	base = strings.TrimSuffix(base, filepath.Ext(base)) // rm extension
	// []string{"jobs", ...}
	str := strings.Split(base, "-")
	if len(str) != 6 {
		return DataFile{}, fmt.Errorf("filename is not in proper structure: %s", base)
	}
	return DataFile {Country: strings.ToLower(str[2]), Month: strings.ToLower(str[3]), Year: str[4]}, nil
}

func main() {
	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("failed to read the file %v", err)
	}
	metadata, err := ParseDataFile(filename)
	if err != nil {
		log.Printf("failed to parse the file %v", err)
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
		_, err := db.Exec(ctx, `INSERT INTO jobs (country, month, year, id, title, company, job_function, role, seniority, locations, skills) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) ON CONFLICT (id, month, year) DO NOTHING`, metadata.Country, metadata.Month, metadata.Year, j.ID, j.Title, j.Company, j.JobFunction, j.Role, j.Seniority, j.Location, j.Skills)
		if err != nil {
			log.Printf("failed importing script %s: %v", j.ID, err)
			continue
		}
		i++
	}
	log.Printf("%d of %d jobs", i, len(resp.Jobs))
}
