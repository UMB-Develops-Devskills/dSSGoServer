build:
	go build -o main ./cmd/devskills

run:
	go run ./cmd/devskills/main.go

routes:
	go run ./cmd/devskills/main.go -routes

# script for generating monthly data for history graph
# gcloud storage cp <storage_location> .
script:
	go run ./scripts/generateMonth.go
