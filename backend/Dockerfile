# ===== build stage ====
FROM golang:1.19.13-bullseye as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -trimpath -ldflags="-w -s" -o app cmd/api/main.go

# ===== deploy stage ====
FROM golang:1.19.13-alpine as deploy

RUN apk update

COPY --from=build /app/app ./app

RUN go install github.com/rubenv/sql-migrate/...@latest
COPY --from=build /app/migrations ./migrations
RUN sql-migrate up

CMD ["./app"]

# ===== dev ====
FROM golang:1.19.13-bullseye as dev

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest
CMD ["air"]
