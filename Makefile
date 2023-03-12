build:
	go build -o fs main.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o fs main.go