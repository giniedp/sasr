GOOS=darwin GOARCH=amd64 go build -o bin/darwin/sasrd tool/main.go
GOOS=linux GOARCH=amd64 go build -o bin/linux/sasrd tool/main.go
GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui -o bin/windows/sasrd tool/main.go