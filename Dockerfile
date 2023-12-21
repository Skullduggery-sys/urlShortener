FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash musl-dev postgresql-client

# Копирование wait-for-postgres.sh и сделать его исполняемым
COPY wait-for-postgres.sh ./

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .

RUN go build -o ./bin/app ./cmd/urlShortener

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /
COPY --from=builder /usr/local/src/wait-for-postgres.sh /
COPY /config/local.yaml /local.yaml

# Для Windows
RUN apk add dos2unix
RUN dos2unix wait-for-postgres.sh

RUN chmod +x /wait-for-postgres.sh
RUN apk --no-cache add postgresql-client

CMD ["/app"]
