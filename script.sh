go mod init <project_name>

go run .

GOOS=windows GOARCH=amd64 go build -o myprogram.exe

go test ./...

go test -cover ./...

go test -v ./...

go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

go get github.com/gocolly/colly