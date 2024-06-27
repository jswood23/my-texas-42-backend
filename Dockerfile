FROM golang:1.21-alpine AS builder
# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY . ./
RUN go vet ./...

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -x -o /my-texas-42-backend

EXPOSE 8080

# Run
CMD [ "/my-texas-42-backend" ]

#go build my-texas-42-backend
#127.0.0.1 5432 mytexas42-staging postgres \"+WYk.dRb/B4Vp`P f4yjqDCa7eze