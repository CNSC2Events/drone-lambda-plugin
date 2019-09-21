FROM golang:1.13.0-alpine AS BUILDER

WORKDIR /builder

ADD . /builder

RUN go build -o main main.go


FROM alpine

WORKDIR /bin

RUN apk update && \
    apk add \
    ca-certificates && \
    rm -rf /var/cache/apk/*

ENV AWS_SDK_LOAD_CONFIG=true

COPY --from=BUILDER /builder/main /bin/main

ENTRYPOINT ["/bin/main"]

