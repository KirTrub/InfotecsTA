FROM golang:1.24.6-alpine
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go mod download
RUN go build -o main ./cmd/InfotecsTA/main.go
EXPOSE 8081
CMD ["./main"]