# go-file-data-stream-reader
A project to test the processing of large files in parts instead of a whole

# How to run

## Run with Docker
Requirements:
- docker installed

Use the make file present at the root and run the commands in order:
```
make build
make run
```
It will build an image with the code and run the project for the fist time, after that use the following command to run the project again with the same docker instance:
```
make start
```
If you want to remove both the docker instance and the image run:
```
make remove
```

## Run locally
Requirements:
- go installed

run the following commands:
```
go mod vendor
go run main.go
```
They are going to donwnload the necessary dependencies and run the project.