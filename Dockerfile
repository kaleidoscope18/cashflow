# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
ARG WORKDIR
WORKDIR $WORKDIR

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 8080 for the application
EXPOSE 8080

ENV PORT=8080

# Set the entry point to run the Go application
CMD ["./main"]