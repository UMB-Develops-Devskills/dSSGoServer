package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/storage"
	//"github.com/jackc/pgx/v5/pgxpool"
)

// this function downloads the skills reference from firebase /metadata folder and stores them as a []SkillRefData
func DownloadSkillRef(ctx context.Context, client *storage.Client, bucketName string, filepath string) (SkillCategories, error) {
	bkt := client.Bucket(bucketName)
	object := bkt.Object(filepath)
	
	reader, err := object.NewReader(ctx) 
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	var sc SkillCategories
	err = json.NewDecoder(reader).Decode(&sc)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

// this function loops through each item in skillsrefdata to match with the target 'skill' string. 
func FindSkill(sc SkillCategories, skill string) *SkillRefData {
	for _, item := range sc {
		if item.Skill == skill {
			return &item
		}
	}
	return nil
}

// this function analyzes the []jobs and the skillsref.json file and returns a skilltrendresponse struct to display skills, specific skill counts and total count of skills. 
func AnalayzeSkills(jobs []Job, sc SkillCategories, skill string, category string, subcat string) (SkillTrendResponse) {
	counts := make(map[string] int)
	md := make(map[string]SkillRefData) // metadata map of skills referenced
	var trends []SkillTrend
	for _, ref := range sc {
		md[ref.Skill] = ref
	}
	totalCount := 0
	for _, job := range jobs {
		for _, skill := range job.Skills {
			ref, ok := md[skill]
			if !ok {
				continue
			}
			if skill != "" && ref.Skill != skill {
				continue
			}
			if category != "" && ref.Category != category {
				continue
			}
			if subcat != "" && ref.Subcategory != subcat {
				continue
			}
			counts[skill]++
			totalCount++
		}
	}
	var totalSkills []string
	for skill, count := range counts {
		ref := md[skill]
		totalSkills = append(totalSkills, skill)
		trends = append(trends, SkillTrend{
			Category: ref.Category,
			Subcategory: ref.Subcategory,
			Skill: skill, 
			Count: count,
		})
	}
	return SkillTrendResponse {
		TotalCount: totalCount, 
		TotalSkills: totalSkills,
		Trends: trends, 
	}
}

// This function displays skill counts from specific role, seniority (entry+mid+intern), country, time period (month + year)
// skills will be organized via skill_categories to show frontend and backend categories and subcategories (languages, frameworks, mobile, db, cloud, devops, data engineering, ai/machine learning, version control, testing)
// example query: /api/trends/skills?country=usa&year=2026&month=july&role=frontend-engineer&seniority=mid
func (h *TaskHandler) GetSkillTrends(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	country := r.URL.Query().Get("country")
	seniority := r.URL.Query().Get("seniority")
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")
	skill := r.URL.Query().Get("skill")
	category := r.URL.Query().Get("category")
	subcategory := r.URL.Query().Get("subcategory")

	filepath := fmt.Sprintf("jobs/%s/year%s/%s/%s/%s-jobs.json", country, year, month, role, seniority)
	jobs, err := DownloadJobs(r.Context(), h.Storage, bucketName, filepath)
	sc, err := DownloadSkillRef(r.Context(), h.Storage, bucketName, "metadata/skillRef.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	trends := AnalayzeSkills(jobs, sc, skill, category, subcategory)
	// write to json
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(trends)
}