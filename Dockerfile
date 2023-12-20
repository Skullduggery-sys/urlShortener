FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .

RUN go build -o ./bin/app ./cmd/urlShortener

FROM alpine AS runner
COPY --from=builder /usr/local/src/bin/app /
COPY /config/local.yaml /local.yaml

CMD ["/app"]
