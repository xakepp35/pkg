FROM golang:1.24 AS builder

# Установка необходимых пакетов
RUN apt-get update && apt-get install -y \
    curl \
    unzip \
    make \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

ENV GOPRIVATE=git.muzpan.com
ENV CGO_ENABLED=0

COPY Makefile .

RUN --mount=type=secret,id=netrc,target=/root/.netrc \
    make prepare

FROM debian:bullseye-slim

COPY --from=builder /usr/bin/make /bin/
COPY --from=builder /usr/local/bin/protoc /usr/local/bin/
COPY --from=builder /usr/local/include /usr/local/include
COPY --from=builder /app /app

WORKDIR /app
ENV PATH="$PATH:/app/bin"

CMD ["make", "versions"]