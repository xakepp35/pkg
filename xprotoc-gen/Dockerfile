FROM golang:1.24-bullseye AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl unzip make ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app
ENV CGO_ENABLED=0

COPY Makefile .
COPY tools tools/
RUN make prepare
RUN rm -rf tools

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    make ca-certificates libstdc++6 zlib1g && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/local/bin/protoc /usr/local/bin/
COPY --from=builder /usr/local/include /usr/local/include
COPY --from=builder /app /app

WORKDIR /app
ENV PATH="/app/bin:${PATH}"

CMD ["make", "versions"]
