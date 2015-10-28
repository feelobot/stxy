build:
	GOOS=linux GOARCH=amd64 go build -o bin/stxy-linux-amd64 && GOOS=darwin GOARH=amd64 go build -o bin/stxy-darwin-amd64

