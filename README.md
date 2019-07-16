## Running project

**Run project`go run main go`**

## Build project
**project for macOS -> `go build`**

**for windows x32 -> `GOOS=windows GOARCH=386 go build -o goserver.exe`**

**for windows x64 -> `GOOS=windows GOARCH=amd64 go build -o goserver.exe`**

**for linux -> `GOOS=linux GOARCH=arm go build -o goserver`**

**after build simple double click file or use `./goserver`**

## Routes
**GET /users/{user} (username || all)**

**POST /user?login=123&password=asdas**

**PATCH /user/{id}?login=123&password=123**