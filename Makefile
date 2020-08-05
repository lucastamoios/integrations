DIST_FOLDER := dist
IMAGE := gcr.io/experiments-283423/lucastamoios/integrations
MIGRATE_CMD := ./cmd/migrate
VERSION ?= $(shell git rev-parse --short HEAD)

.PHONY: db

build:
	go build -o $(DIST_FOLDER)/migrate $(MIGRATE_CMD)
	go build -o $(DIST_FOLDER)/server ./internals/server.go

clean:
	rm -rf $(DIST_FOLDER)

db:
	docker-compose up -d db

db-shell:
	docker-compose run --rm db psql

deploy:
	sed -i 's/<VERSION>/$(VERSION)/' k8s/api/deployment.yaml
	kubectl apply -f k8s/api/
	kubectl apply -f k8s/postgres/

docker:
	docker build -t $(IMAGE) .

docker-push:
	docker tag $(IMAGE):latest $(IMAGE):$(VERSION)
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest

migrate:
	docker-compose run --rm app go run $(MIGRATE_CMD)

run:
	docker-compose run --rm app

shell:
	docker-compose run --rm app sh
