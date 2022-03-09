FROM golang:1.17-alpine as builder

COPY ./ /app
WORKDIR /app

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app main.go

CMD ["./app"]