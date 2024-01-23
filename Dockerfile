FROM golang:1.21.3-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's working directory
COPY . .

# Build the Go application
RUN go build -o main .

# Command to run the executable
CMD ["./main"]
