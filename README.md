# Metrics Service

## Description
A service for generating and sending metrics.  
Business metrics are obtained from the database.

### There are two types of metrics provided by the service
1) `Prometheus` metrics, metrics, available at `some.example.ru/metrics-app/metrics`:
  - `policies_purchased_today_gauge` - Number of policies purchased today
  - `authentications_today_gauge` - Number of authentications today

<br>

2) Data tables sent to specified email addresses from `fatal-alert@some.example.ru`:
  - Data on certificate issues  
  - Data on paid but not issued policies

## Building and Running the Application
The project was developed with `Go` - `1.22.1`.  

The application should be run in the root directory of the project (where `main.go` is located).  
You can run the application with the command `go run main.go --contour <local/demo/preprod/prod>`.

You can also build the application into an executable file.  
This is done with the command `go build -o ./metrics-app ./main.go` for `Linux/MacOS`
or `go build -o ./metrics-app.exe ./main.go` for `Windows`.  
The application will be built specifically for the OS and processor architecture
on which the build command was run.
If you need to build the application for another case,
you must first set the environment variables `GOOS` and `GOARCH`.

Examples:
- `env GOOS=linux GOARCH=amd64 go build -o metrics-app main.go`
- `env GOOS=darwin GOARCH=arm64 go build -o metrics-app main.go`
- `env GOOS=windows GOARCH=386 go build -o metrics-app.exe main.go`

You can find the full list of available options with the command `go tool dist list`.

## Configuration
The `./config` directory should contain at least one configuration file for the required contour. 
The file name should be formatted as follows: `config.<local/demo/preprod/prod>.json`.   
Example: `config.demo.json`  

When starting the application, you can specify the required contour (as shown in the [Building and Running the Application](#building-and-running-the-application) section).
