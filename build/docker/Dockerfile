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

COPY . /src/payments

WORKDIR /src/payments

RUN make build

FROM alpine:3.8
WORKDIR /
COPY --from=0 /src/payments/target/* /bin/
COPY --from=0 /bin/wait-for /bin/wait-for
ENTRYPOINT ["/bin/payments-api"]
