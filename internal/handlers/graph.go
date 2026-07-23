package handlers

import (
	"encoding/json"
	"sort"
	"fmt"
	"net/http"
)

// this function analyzes each job posting for skills, counts each skill and pairs two skills.  
// Given one skill is present, how many times does another skill co-occur in the same job description
func AnalyzeSkillPair(jobs []Job, sc SkillCategories) NetworkGraphData {
	skillCount := make(map[string]int)
	pairCount := make(map[SkillPair]int)
	md := make(map[string]SkillRefData) // metadata list of skills referenced
	for _, ref := range sc {
		md[ref.Skill] = ref
	}
	// for each unique skill found increment the skillcount by 1
	for _, job := range jobs {
		unique := make(map[string]bool) 
		for _, skill := range job.Skills {
			// if there is a duplicated skill, set to true
			if unique[skill] {
				continue
			}
			unique[skill] = true 
			skillCount[skill]++
		}
		skills := make([]string, 0, len(unique))
		for skill := range unique {
			skills = append(skills, skill)
		}
		// sort the list
		sort.Strings(skills) 
		/*
		/* for every pair found increment the pair counter by 1
		ex: map[
    SkillPair{"Go","Docker"}] = 2
		SkillPair{"skill1", "skill2"}] = cooccurrence_count
		*/
		for i:=0; i<len(skills); i++ {
			for j:=i+1; j<len(skills); j++ {
				pair := SkillPair {
					Skill1: skills[i],
					Skill2: skills[j],
				}
				pairCount[pair]++
			}
		}
	}
	// build the nodes list
	nodes := make([]Node, 0, len(skillCount))
	for skill, count := range skillCount {
		ref := md[skill]
		nodes = append(nodes, Node{
			SkillID: skill,
			Count: count,
			Category: ref.Category,
		})
	}
	// build the links list
	links := make([]Link, 0, len(pairCount))
	// for each pair, count append to the list 
	for pair, count := range pairCount {
		links = append(links, Link{
			SkillSource: pair.Skill1,
			SkillTarget: pair.Skill2,
			Count: count,
		})
	}
	return NetworkGraphData{
		Nodes: nodes,
		Links:links, 
	}
} 

// this function displays skill counts and displays skill pairs for co occurrence network graph data
func (h *TaskHandler) GetSkillGraph(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	country := r.URL.Query().Get("country")
	seniority := r.URL.Query().Get("seniority")
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")

	filepath := fmt.Sprintf("jobs/%s/year%s/%s/%s/%s-jobs.json", country, year, month, role, seniority)
	jobs, err := DownloadJobs(r.Context(), h.Storage, bucketName, filepath)
	sc, err := DownloadSkillRef(r.Context(), h.Storage, bucketName, "metadata/skillRef.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	networkGraph := AnalyzeSkillPair(jobs, sc)
	// write to json
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(networkGraph)
}