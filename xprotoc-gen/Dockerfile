FROM golang:1.24 AS builder

RUN apt-get update && apt-get install -y \
    curl \
    unzip \
    make

WORKDIR /app

ENV CGO_ENABLED=0

COPY Makefile .
COPY tools tools/

RUN go build -C ./tools/swagger-patcher -o /usr/local/bin/patch-swagger

RUN make prepare
RUN rm -rf tools

FROM debian:bullseye-slim

COPY --from=builder /usr/bin/make /bin/
COPY --from=builder /usr/local/bin/protoc /usr/local/bin/
COPY --from=builder /usr/local/bin/patch-swagger /usr/local/bin/
COPY --from=builder /usr/local/include /usr/local/include
COPY --from=builder /app /app

WORKDIR /app
ENV PATH="$PATH:/app/bin"

CMD ["make", "versions"]