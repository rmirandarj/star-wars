#!make

build_directory=./build
app_name=server

build: 	## Generate application binaries
	@go mod download && \
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -v -installsuffix cgo -o ${build_directory}/${app_name} cmd/main.go || exit $?

dependency:
	cd cmd/; go mod download

test:
	CGO_ENABLED=0 go test ./...

coverall:
	CGO_ENABLED=0 go test ./... -coverprofile ./cover.out && \
		go tool cover -func ./cover.out && \
		rm ./cover.out

test-short:
	CGO_ENABLED=0 go test ./... -short

start:
	docker-compose up --build

stop:
	docker-compose down --remove-orphans

clean:
	go clean -cache

env-up:
	docker-compose up -d postgres

env-down:
	docker-compose down --remove-orphans

run-local: dependency
	go run cmd/main.go