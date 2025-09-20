# Frontend build stage
FROM node:20-alpine AS frontend-builder

WORKDIR /app

# Copy package files from root directory (where the custom Svelte app is)
COPY package.json package-lock.json* ./
RUN npm install

# Copy the entire source including src/ directory and vite.config.js
COPY src/ ./src/
COPY index.html ./
COPY vite.config.js ./

# Build the frontend from the root (which includes src/)
RUN npm run build

# Backend build stage
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary and frontend dist from builders
COPY --from=backend-builder /app/main .
COPY --from=frontend-builder /app/dist /app/dist

# Create required directories
RUN mkdir -p uploads

# Set permissions
RUN chmod +x /app/main
RUN chown -R nobody:nobody /app

# Switch to non-root user
USER nobody

# Expose port
EXPOSE 3000

# Run the binary
CMD ["./main"]