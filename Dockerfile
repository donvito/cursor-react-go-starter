# Build frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Build backend
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
COPY backend/ ./backend/
WORKDIR /app/backend
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:3.19
WORKDIR /app

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates python3

# Copy frontend build
COPY --from=frontend-builder /app/frontend/dist /app/frontend/dist

# Copy backend binary
COPY --from=backend-builder /app/backend/main /app/backend/main

# Create data directory for BadgerDB
RUN mkdir -p /app/backend/data && \
    chmod 777 /app/backend/data

# Expose ports
EXPOSE 8080 5173

# Create startup script
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'cd /app/backend && ./main &' >> /app/start.sh && \
    echo 'cd /app/frontend/dist && python3 -m http.server 5173' >> /app/start.sh && \
    chmod +x /app/start.sh

# Set environment variables
ENV FRONTEND_URL=http://localhost:5173
ENV BACKEND_URL=http://localhost:8080

CMD ["/app/start.sh"] 