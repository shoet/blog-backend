# ===== build stage ====
FROM golang:1.24.2-bullseye as builder

WORKDIR /app

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/gomod-cache \
    go mod download

COPY . .

RUN --mount=type=cache,target=/gomod-cache \
    --mount=type=cache,target=/go-cache \
    go build -trimpath -ldflags="-w -s" -o cmd/bin/api cmd/api/main.go

RUN --mount=type=cache,target=/gomod-cache \
    --mount=type=cache,target=/go-cache \
  go build -trimpath -ldflags="-w -s" -o cmd/bin/cli cmd/cli/main.go

# ===== local development stage ====
FROM golang:1.24.2-bullseye as dev

WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.51.0
CMD ["air"]

# ===== deploy stage ====
FROM golang:1.24.2-bullseye as production

ARG PORT

WORKDIR /app

RUN apt update -y

COPY --from=builder /app/cmd/bin/api .
COPY --from=builder /app/cmd/bin/cli .

COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:0.9.0 /lambda-adapter /opt/extensions/lambda-adapter

ENV PORT=${PORT:-8080}
ENV READINESS_CHECK_PATH=/health

EXPOSE ${PORT}

CMD ["/app/api"]

