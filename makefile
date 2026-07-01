build:
	go build -o main ./cmd/devskills

run:
	go run ./cmd/devskills/main.go

routes:
	go run ./cmd/devskills/main.go -routes

# script
script:
	go run cmd/import/main.go data/india/jobs-default-india-june-2026-page4.json