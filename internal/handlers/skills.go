package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	//"github.com/jackc/pgx/v5/pgxpool"
)

// This function requests skill counts from specific role, seniority (entry+mid+intern), country, company, time period (month + year)
// skills will be organized via skill_categories to show frontend and backend categories and subcategories (languages, frameworks, mobile, db, cloud, devops, data engineering, ai/machine learning, version control, testing)
// query example:
// /api/trends/skills?country=usa&month=june&year=2026
// /api/trends/skills?role=Software+Engineer&country=usa&month=june&year=2026
func (h *TaskHandler) GetSkillTrends(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	country := r.URL.Query().Get("country")
	seniority := r.URL.Query().Get("seniority")
	company := r.URL.Query().Get("company") // ex. all US. companies
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")
	category := r.URL.Query().Get("category")
	limit := r.URL.Query().Get("limit")

	query := `SELECT sc.category, sc.subcategory, u.skill, COUNT(*)
	FROM jobs j
	CROSS JOIN UNNEST(j.skills) AS u(skill)
	JOIN skill_categories sc ON sc.skill = u.skill
	WHERE ($1::text IS NULL OR role = $1)
	AND ($2::text IS NULL OR country = $2)
	AND ($3::text IS NULL OR $3 = ANY(j.seniority))
	AND ($4::text IS NULL OR company = $4)
	AND ($5::text IS NULL OR month = $5)
	AND ($6::text IS NULL OR year = $6)
	AND ($7::text IS NULL OR sc.category = $7)
	GROUP BY sc.category, sc.subcategory, u.skill
	ORDER BY COUNT(*) DESC
	LIMIT $8;`

	rows, err := h.DB.Query(r.Context(), query, nullIfEmpty(role), nullIfEmpty(country), nullIfEmpty(seniority), nullIfEmpty(company), nullIfEmpty(month), nullIfEmpty(year), nullIfEmpty(limit), nullIfEmpty(category))
	if err != nil {
		log.Printf("Skill trends query error: %v", err)
		http.Error(w, "failed to query skills trends", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	// skill trends array
	var skills []SkillTrend
	for rows.Next() {
		var s SkillTrend
		if err := rows.Scan(&s.Category, &s.Subcategory, &s.Skill, &s.Count); err != nil {
			log.Printf("scan error: %v", err)
			http.Error(w, "failed to scan skills", http.StatusInternalServerError)
			return
		}
		skills = append(skills, s)
	}
	// if no skills exist return an empty array
	if skills == nil {
		skills = []SkillTrend{}
	}
	// write to json
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(skills)
}