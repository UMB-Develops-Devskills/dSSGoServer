package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	//"github.com/jackc/pgx/v5/pgxpool"
)

// this function requests the count of work location and work type (remote) from specific role, seniority (entry+mid+intern), country, company, time period (month + year)
// example query:
// /api/trends/locations
// /api/trends/locations?role=Software+Engineer&country=usa&month=june&year=2026
func (h *TaskHandler) GetWorkLocationTrend(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	country := r.URL.Query().Get("country")
	seniority := r.URL.Query().Get("seniority")
	company := r.URL.Query().Get("company") // ex. all US. companies
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")
	limit := r.URL.Query().Get("limit")

	query := `SELECT location, COUNT(*)
	FROM jobs,
	UNNEST(locations) AS location
	WHERE ($1::text IS NULL OR role = $1)
	AND ($2::text IS NULL OR country = $2)
	AND ($3::text IS NULL OR $3 = ANY(seniority))
	AND ($4::text IS NULL OR company = $4)
	AND ($5::text IS NULL OR month = $5)
	AND ($6::text IS NULL OR year = $6)
	GROUP BY location
	ORDER BY COUNT(*) DESC
	LIMIT $7;`

	rows, err := h.DB.Query(r.Context(), query, nullIfEmpty(role), nullIfEmpty(country), nullIfEmpty(seniority), nullIfEmpty(company), nullIfEmpty(month), nullIfEmpty(year), nullIfEmpty(limit))
	if err != nil {
		log.Printf("location trends query error: %v", err)
		http.Error(w, "failed to query location trends", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	// location trends array
	var locations []LocationTrend
	for rows.Next() {
		var l LocationTrend
		if err := rows.Scan(&l.Location, &l.Count); err != nil {
			log.Printf("scan error: %v", err)
			http.Error(w, "failed to scan locations", http.StatusInternalServerError)
			return
		}
		locations = append(locations, l)
	}
	// if no locations exist return an empty array
	if locations == nil {
		locations = []LocationTrend{}
	}
	// write to json
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(locations)
}