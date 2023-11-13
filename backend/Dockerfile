# ===== build stage ====
FROM golang:1.19.13-bullseye as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -trimpath -ldflags="-w -s" -o app cmd/api/main.go

# ===== deploy stage ====
FROM golang:1.19.13-bullseye as deploy

WORKDIR /app

RUN apt update -y

COPY --from=builder /app/app .

# RUN go install github.com/rubenv/sql-migrate/...@latest
# COPY --from=build /app/_tools/migrations ./migrations
# COPY --from=build /app/_tools/dbconfig.yml ./dbconfig.yml
# RUN sql-migrate up

CMD ["/app/app"]

# ===== dev ====
FROM golang:1.19.13-bullseye as dev

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest
CMD ["air"]
