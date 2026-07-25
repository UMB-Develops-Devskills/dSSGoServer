package main

import (
	"context"
	"encoding/json"
	"log"
	"fmt"
	"cloud.google.com/go/storage"
	"github.com/UMB-Develops-Devskills/dSSGoServer/internal/handlers"
	"github.com/UMB-Develops-Devskills/dSSGoServer/internal/firebase"
)

const bucketName = "devskills-499815.firebasestorage.app"

type SkillCount struct {
	Skill string `json:"skill"`
	Count int `json:"count"`
}

type MonthlySkillCounts struct {
	Location string `json:"location"`
	Year string`json:"year"`
	Month string `json:"month"`
	TotalJobs int `json:"total_jobs"`
	Skills []SkillCount `json:"skills"`
}

func GenerateMonth(ctx context.Context, client *storage.Client, buckeName string, location string, year string, month string) (error) {
	counts := make(map[string]int)
	totalJobs := 0
	filepath := fmt.Sprintf("jobs/%s/year%s/%s/", location, year, month)
	files, err := handlers.GetJobFiles(ctx, client, bucketName, filepath)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Found %d files", len(files))
	// loop through each file and download the job postings then increment totaljobs and skill counter 
	for _, file := range files {
		log.Println(file)
		jobs, err := handlers.DownloadJobs(ctx, client, bucketName, file)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, job := range jobs {
			totalJobs++
			for _, skill := range job.Skills {
				counts[skill]++
			}
		}
	}
	var skills []SkillCount
	for skill, count := range counts {
		skills = append(skills, SkillCount{
			Skill: skill,
			Count: count,
		})
	}
	result := MonthlySkillCounts{
		Location: location,
		Year: year,
		Month: month,
		TotalJobs: totalJobs,
		Skills: skills,
	}
	// write to json
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Println(err)
	}
	fpOutput := fmt.Sprintf("history/%s/%s/%s.json", location, year, month) // where to send the file to
	return handlers.UploadFile(ctx, client, bucketName, fpOutput, data)
}

func main() {
	// *** edit these arrays to generate a month of total skill count for a specific location, year and month
	locations := []string{"usa", "california"}
	years := []string {"2026"}
	months := []string {"july"}
	// Connect to Firebase
	ctx := context.Background()
	client, err := firebase.ConnectFireBaseStorage(ctx)
	if err != nil {
		log.Printf("firebase connection failed: %v", err)
	}
	defer client.Close()
	for _, location := range locations {
		for _, year := range years {
			for _, month := range months {
				err := GenerateMonth(ctx, client, bucketName, location, year, month)
				if err != nil {
					log.Println(location, err)
					continue 
				}
				log.Printf("generated file: %s %s %s", location, year, month)
			}
		}
	}
	log.Println("finished generating file(s), type gcloud storage cp <storage_location> . to download/copy the file locally")
}