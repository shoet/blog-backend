ARG DB_REPLICA_REMOTE_PATH
ARG LITESTREAM_ACCESS_KEY_ID
ARG LITESTREAM_SECRET_ACCESS_KEY

# ===== build stage ====
FROM golang:1.19.13-bullseye as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ADD https://github.com/benbjohnson/litestream/releases/download/v0.3.9/litestream-v0.3.9-linux-amd64-static.tar.gz litestream.tar.gz
RUN tar -xzf litestream.tar.gz -C ./

RUN go build -trimpath -ldflags="-w -s" -o app cmd/api/main.go

# ===== deploy stage ====
FROM golang:1.19.13-alpine as deploy

ARG DB_REPLICA_REMOTE_PATH
ARG LITESTREAM_ACCESS_KEY_ID
ARG LITESTREAM_SECRET_ACCESS_KEY

RUN apk update

COPY --from=build /app/litestream /usr/local/bin/litestream
COPY --from=build /app/restore.sh ./restore.sh
RUN chmod +x ./restore.sh
RUN ./restore.sh

COPY --from=build /app/litestream /usr/local/bin/litestream
COPY --from=build /app/app ./app

CMD ["litestream", "replicate", "-exec", "./app", "./database.sqlite", "${DB_REPLICA_REMOTE_PATH}"]

# ===== dev ====
FROM golang:1.19.13-bullseye as dev

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest
CMD ["air"]
