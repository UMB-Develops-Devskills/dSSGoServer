package handlers 

import (
	"fmt"
	"net/http"
	"encoding/json"
)

func AnalyzeHistory(histories []HistoryFile, skill string) (HistoryGraphResponse) {
	var hist []History
	var totalSkills [] string
	var month, year, location string 
	var totalJobs int
	for _, history := range histories {
		month = history.Month
		year = history.Year
		location = history.Location
		totalJobs = history.TotalJobs
		for _, historySkill := range history.HistorySkillCounts {
			totalSkills = append(totalSkills, historySkill.Skill)
			// if skill, ignore the rest
			if skill != "" && historySkill.Skill != skill {
				continue
			}
			hist = append(hist, History{
				Month: history.Month, 
				Year: history.Year,
				Skill: historySkill.Skill,
				Count: historySkill.Count,
			})
		}
	}
	return HistoryGraphResponse{
		Skill: skill,
		TotalSkills: totalSkills,
		TotalJobs: totalJobs,
		Location: location, 
		Month: month,
		Year: year, 
		History: hist,
	}
}

// this function displays the history of a specific location, year, month and skill 
// the history line graph will show the count of each skill out of total_jobs
func (h *TaskHandler) GetHistoryData(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	year := r.URL.Query().Get("year")
	month := r.URL.Query().Get("month")
	skill := r.URL.Query().Get("skill")
	var histories []HistoryFile
	filepath := fmt.Sprintf("history/%s/%s/%s.json", location, year, month)
	hf, err := DownloadHistory(r.Context(), h.Storage, bucketName, filepath)
	if err != nil {
		http.Error(w, "500 error, could not find history filepath to firebase " + err.Error(), http.StatusInternalServerError)
		return
	}
	histories = append(histories, hf)
	hist := AnalyzeHistory(histories, skill)
	// write to json
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(hist)
}