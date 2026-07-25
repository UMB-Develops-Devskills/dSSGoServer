package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	//"os"
	//"github.com/go-chi/chi/v5"
	//"github.com/jackc/pgx/v5/pgxpool"
	//"github.com/joho/godotenv"
)

// this function supports /api/keys handler by listing objects in the bucket
func GetBucketData(ctx context.Context, client *storage.Client, bucketName string) ([]BucketData, error) {
	bkt := client.Bucket(bucketName) 
	query := &storage.Query{Prefix: "jobs/"}
	it := bkt.Objects(ctx, query)
	// this array holds the objects
	var objects []BucketData
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		// folders
		if strings.HasSuffix(attrs.Name, "/") {
			continue 
		}
		// files 
		if !strings.HasSuffix(attrs.Name, "-jobs.json") {
			continue
		}
		parts := strings.Split(attrs.Name, "/")
		// ex: "jobs/california/year2026/july/machine-learning-engineer/senior-jobs.json" -> 6 parts
		if len(parts) != 6 {
			continue
		}
		year := strings.TrimPrefix(parts[2], "year") // returns the string specific year
		seniority := strings.TrimSuffix(parts[5], "-jobs.json") // returns the string specific role

		objects = append(objects, BucketData{
			Location: parts[1], 
			Year: year,
			Month: parts[3],
			Role: parts[4],
			Seniority: seniority,
			}) 
	}
	return objects, nil
}

// this function handles retrieving the object keys and values in the bucket and parses them in json
func (h *TaskHandler) GetKeyData(w http.ResponseWriter, r *http.Request) {
	objects, err := GetBucketData(r.Context(), h.Storage, bucketName)
	if err != nil {
		http.Error(w, "500 error, could not find key filepath to firebase " + err.Error(), http.StatusInternalServerError)
		return
	}
	// write to json
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(objects)
}