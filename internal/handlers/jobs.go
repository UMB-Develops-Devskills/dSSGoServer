package handlers

import (
	"log"
	"net/http"
	//"os"
	"github.com/joho/godotenv"
)

func GetJobData(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//dataLakeKey := os.Getenv("DATA_LAKE_API")
}