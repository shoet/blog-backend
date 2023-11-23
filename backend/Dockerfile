# ===== build stage ====
FROM golang:1.19.13-bullseye as builder

WORKDIR /app

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go-mod-cache \
    go mod download

COPY . .

RUN --mount=type=cache,target=/gomod-cache \
    --mount=type=cache,target=/go-cache \
    go build -trimpath -ldflags="-w -s" -o cmd/bin/api cmd/api/main.go

RUN --mount=type=cache,target=/gomod-cache \
    --mount=type=cache,target=/go-cache \
  go build -trimpath -ldflags="-w -s" -o cmd/bin/cli cmd/cli/main.go

# ===== deploy stage ====
FROM golang:1.19.13-bullseye as deploy

WORKDIR /app

RUN apt update -y

COPY --from=builder /app/cmd/bin/api .
COPY --from=builder /app/cmd/bin/cli .

CMD ["/app/api"]

# ===== dev ====
FROM golang:1.19.13-bullseye as dev

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest
CMD ["air"]
