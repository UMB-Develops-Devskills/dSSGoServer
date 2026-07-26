package handlers

import (
	"encoding/json"
	"sort"
	"fmt"
	"net/http"
)

// this function analyzes each job posting for skills, counts each skill and pairs two skills.  
// Given one skill is present, how many times does another skill co-occur in the same job description
func AnalyzeSkillPair(jobs []Job, sc SkillCategories, skill string) NetworkGraphDataResponse {
	skillCount := make(map[string]int)
	pairCount := make(map[SkillPair]int)
	totalSkills := []string{}
	md := make(map[string]SkillRefData) // metadata list of skills referenced
	for _, ref := range sc {
		md[ref.Skill] = ref
	}
	// for each unique skill found increment the skillcount by 1
	for _, job := range jobs {
		unique := make(map[string]bool) 
		for _, jobSkill := range job.Skills {
			// if there is a duplicated skill, set to true
			if unique[jobSkill] {
				continue
			}
			unique[jobSkill] = true 
		}
		if skill != "" && !unique[skill] {
			continue
		}
		skills := make([]string, 0, len(unique))
		for jobSkill := range unique {
			skillCount[jobSkill]++
			skills = append(skills, jobSkill)
		}
		// sort the list
		sort.Strings(skills) 
		/*
		/* for every pair found increment the pair counter by 1
		ex: map[
    SkillPair{"Go","Docker"}] = 2
		SkillPair{"skill1", "skill2"}] = cooccurrence_count
		*/
		for _, other := range skills {
			if other == skill {
				continue 
			}
			pair := SkillPair {
				Skill1: skill,
				Skill2: other, 
			}
			pairCount[pair]++
		}
	}
	// build the nodes list
	nodes := make([]Node, 0, len(skillCount))
	for skill, count := range skillCount {
		totalSkills = append(totalSkills, skill)
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
	return NetworkGraphDataResponse{
		Skill: skill,
		TotalSkills: totalSkills,
		Nodes: nodes,
		Links:links, 
	}
} 

// this function displays skill counts and displays skill pairs for co occurrence network graph data
func (h *TaskHandler) GetSkillGraph(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	location := r.URL.Query().Get("location")
	seniority := r.URL.Query().Get("seniority")
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")
	skill := r.URL.Query().Get("skill")

	filepath := fmt.Sprintf("jobs/%s/year%s/%s/%s/%s-jobs.json", location, year, month, role, seniority)
	jobs, err := DownloadJobs(r.Context(), h.Storage, bucketName, filepath)
	if err != nil {
		http.Error(w, "500 error, could not find jobs filepath to firebase " + err.Error(), http.StatusInternalServerError)
		return
	}
	sc, err := DownloadSkillRef(r.Context(), h.Storage, bucketName, "metadata/skillRef.json")
	if err != nil {
		http.Error(w, "500 error, could not find skillref filepath to firebase " + err.Error(), http.StatusInternalServerError)
		return
	}
	networkGraph := AnalyzeSkillPair(jobs, sc, skill)
	// write to json
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(networkGraph)
}