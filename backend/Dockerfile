# Use the official Golang image as a base
FROM golang:1.20

WORKDIR /backend

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["go", "run", "server.go"]
