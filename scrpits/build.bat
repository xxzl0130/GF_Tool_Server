set GOARCH=386
set GOOS=windows
packr build -o GF_Tool_Server.exe ./main.go ./proxy.go ./chip.go ./html.go