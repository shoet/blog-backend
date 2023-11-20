# ===== build stage ====
FROM golang:1.19.13-bullseye as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -trimpath -ldflags="-w -s" -o cmd/bin/api cmd/api/main.go
RUN go build -trimpath -ldflags="-w -s" -o cmd/bin/cli cmd/cli/main.go

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
