FROM golang:alpine AS builder

ENV GOPATH="$HOME/go"

WORKDIR $GOPATH/src/work-at-olist

COPY . $GOPATH/src/work-at-olist

RUN apk update && apk add curl git && \
    curl https://glide.sh/get | sh && \
    glide up && \
    go build

FROM alpine:latest

ENV GOPATH="$HOME/go"

WORKDIR /app

COPY --from=builder $GOPATH/src/work-at-olist .

ENTRYPOINT ["./work-at-olist"]
CMD ["run"]
