FROM golang:1.11-alpine

COPY ./build/docker/wait-for /bin/

RUN apk --no-cache --upgrade add ca-certificates \
    && update-ca-certificates --fresh \
    && apk --no-cache add --upgrade \
    git \
    make \
    gcc \
    musl-dev \
    && chmod +X /bin/wait-for

COPY . /usr/local/src/payments

WORKDIR /usr/local/src/payments

RUN make build

FROM alpine:3.8
WORKDIR /

COPY --from=0 /usr/local/src/payments/target/* /usr/local/payments/
COPY --from=0 /usr/local/src/payments/api/* /usr/local/payments/static/
COPY --from=0 /bin/wait-for /bin/wait-for

WORKDIR /usr/local/payments/

ENTRYPOINT ["/usr/local/payments/server"]
