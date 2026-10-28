FROM golang:1.26 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build  -o /bin/app ./cmd/api

FROM scratch

WORKDIR /app

# The binary and the per-client config the app reads at startup.
COPY --from=build /bin/app /app/app
COPY --from=build /src/mock_db.json /app/mock_db.json

# Bind to all interfaces inside the container, and keep the writable SQLite
# database under the working directory.
ENV RATE_LIMITER_HOST=0.0.0.0 \
    RATE_LIMITER_PORT=8000 \
    RATE_LIMITER_MOCK_DB=/app/mock_db.json \
    RATE_LIMITER_SQLITE_DB=/app/ratelimiter.db

EXPOSE 8000

ENTRYPOINT ["/app/app"]
