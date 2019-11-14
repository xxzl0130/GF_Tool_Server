set GOARCH=amd64
set GOOS=linux
go build -o GF_Tool_Server ./main.go ./proxy.go ./chip.go ./html.go