build:
	go build -o main ./cmd/devskills

run:
	go run ./cmd/devskills/main.go

routes:
	go run ./cmd/devskills/main.go -routes

# script for testing there are total of 10 pages
script:
	go run cmd/import/main.go data/july2026/canada/jobs-default-canada-july-2026-page10.json
