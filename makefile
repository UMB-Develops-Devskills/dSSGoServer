build:
	go build -o main ./cmd/devskills

run:
	go run ./cmd/devskills/main.go

routes:
	go run ./cmd/devskills/main.go -routes

# script
script:
	go run cmd/import/main.go data/usa/jobs-default-usa-june-2026-page1.json