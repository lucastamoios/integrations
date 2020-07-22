BINARY := server
DIST_FOLDER := dist
IMAGE := lucastamoios/integrations

build:
	go build -o $(DIST_FOLDER)/$(BINARY) ./internals/server.go

clean:
	rm -rf $(DIST_FOLDER)

db-shell:
	docker-compose run --rm db psql

docker:
	docker build -t $(IMAGE) .

run:
	docker-compose run --rm app
