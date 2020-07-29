BINARY := server
DIST_FOLDER := dist
IMAGE := lucastamoios/integrations

.PHONY: db

build:
	go build -o $(DIST_FOLDER)/$(BINARY) ./internals/server.go

clean:
	rm -rf $(DIST_FOLDER)

db:
	docker-compose up -d db

db-shell:
	docker-compose run --rm db psql

docker:
	docker build -t $(IMAGE) .

migrate:
	docker-compose run --rm app go run ./cmd/migrate

run:
	docker-compose run --rm app

shell:
	docker-compose run --rm app sh
