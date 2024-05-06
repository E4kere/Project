# Step 1: Use Golang official image as a builder
FROM golang:1.22.1 as builder

# Step 2: Set the working directory to the Guns project folder
WORKDIR /app/Guns

# Step 3: Copy go.mod and go.sum into the working directory
COPY Guns/go.mod Guns/go.sum ./

# Step 4: Download dependencies
RUN go mod download

# Step 5: Copy the entire project into the working directory
COPY Guns/ ./

# Step 6: Build the application (adjust the binary name if needed)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/GunsBinary ./cmd/main.go

# Step 7: Use a smaller image to run the app
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Step 8: Set the working directory inside the container
WORKDIR /root/

# Step 9: Copy the pre-built binary file from the previous stage
COPY --from=builder /app/GunsBinary .

# Step 10: Copy the migration files
COPY --from=builder /app/Guns/migrations ./migrations

# Step 11: Command to run the executable (adjust the binary name if needed)
CMD ["./GunsBinary"]
