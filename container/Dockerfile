# faz o build da aplicação no container
FROM golang:1.23.0 AS builder

WORKDIR /app

COPY *.go .
COPY go.mod go.sum .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o bootstrap

# faz o deploy final da aplicação
FROM alpine:3.20.2

WORKDIR /app

COPY --from=builder /app/bootstrap .

RUN addgroup -S nonroot 
RUN adduser -S nonroot -G nonroot
RUN chown -R nonroot:nonroot /app

USER nonroot

ENTRYPOINT ["/app/bootstrap"]

# build da imagem:
#   docker build -t sqsconsumer:1.0 .
#   docker image tag sqsconsumer:1.0 registry.i9.com/sqsconsumer:1.0
#   docker image tag sqsconsumer:1.0 registry.i9.com/sqsconsumer:latest
#   docker image push --all-tags registry.i9.com/sqsconsumer
