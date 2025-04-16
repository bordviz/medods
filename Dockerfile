FROM golang:1.24.1

WORKDIR /medods

COPY go.mod go.sum ./
RUN go mod tidy
COPY . .

RUN go build -o main cmd/main.go

EXPOSE 8080

CMD ["./main"]
