package handlers

import (
	"cloud.google.com/go/storage"
)

const bucketName = "devskills-499815.firebasestorage.app"


// outlines the folder structure data of the files in firebase bucket
type BucketData struct {
	Location string `json:"location"`
	Year string `json:"year"`
	Month string `json:"month"`
	Role string `json:"role"`
	Seniority string `json:"seniority"`
}

// returns the bucketdata response for frontend dataset filter in json
type BucketResponse struct {
	BD []BucketData `json:"bucketdata"`
}

type JobFile struct {
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
	Skills []string `json:"required_skills"`
}

type SkillRefData struct {
	Skill string `json:"skill"`
	Category string `json:"category"`
	Subcategory string `json:"subcategory"`
}

type SkillCategories []SkillRefData

// describes the skill trends
type SkillTrend struct {
	Category string `json:"category"`
	Subcategory string `json:"subcategory"`
	Skill string `json:"skill"`
	Count int `json:"count"`
}

// describes the outer skill trends response
type SkillTrendResponse struct {
	TotalCount int `json:"total_count_skills"`
	TotalSkills []string `json:"total_mentioned_skills"`
	Trends []SkillTrend `json:"trends"`
}

// This struct holds dependencies for http handlers
type TaskHandler struct {
	Storage *storage.Client
}