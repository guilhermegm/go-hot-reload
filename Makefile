build:
	GOOS=linux GOARCH=amd64 go build -o gohr_linux
	GOOS=darwin GOARCH=amd64 go build -o gohr_macos
