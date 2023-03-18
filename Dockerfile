# Start from a Debian-based image with Go installed
FROM golang:1.16

# Set the working directory to /app
WORKDIR /chatapp

# Copy the Go modules file and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go get golang.org/x/net
RUN go mod tidy


# Copy the rest of the application source code
COPY . .

RUN go mod download golang.org/x/net
# Build the Go chatapp
RUN go build -o main .

# Expose port 8000 for the web server
EXPOSE 8000

# Run the Go app
CMD ["./main"]

# build the Docker image by running the following command
# docker build -t chatapp .


# After the build is complete, you can run the Docker container using the following command:
# docker run -p 8000:8000 chatapp