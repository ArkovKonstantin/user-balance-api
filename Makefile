run:
	mkdir -p "bin" && GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/main main.go
	docker-compose up --build
test:
	go test -v ./repository/
