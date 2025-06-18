install:
	go mod vendor
	go mod download

build:
	go build -o api

run:
	./api

docker-up:
	docker compose -f docker-compose-local.yml up -d

docker-down:
	docker compose -f docker-compose-local.yml down

run-air:
	air