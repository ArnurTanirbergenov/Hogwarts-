# Start from a base image, e.g., golang:1.18-alpine
FROM golang:1.21.1 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

COPY ./ ./

COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o letsgo ./cmd/api/

FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

RUN ls -la

COPY --from=builder /app/letsgo .
COPY --from=builder /app/migrations ./migrations


# Command to run the executable
CMD ["./letsgo"]

# # Use an official PostgreSQL image
# FROM postgres:latest

# # Set necessary environmet variables (optional)
# ENV POSTGRES_DB=hogwarts
# ENV POSTGRES_USER=hogwarts
# ENV POSTGRES_PASSWORD=paSSword

# # Copy initialization scripts if needed (optional)
# COPY ./init.sql /docker-entrypoint-initdb.d/