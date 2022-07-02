FROM golang:1.18 AS builder

WORKDIR /app

COPY . .

RUN make build_backend_linux_amd64


FROM scratch

COPY --from=builder /app/build/bin/backend-linux-amd64 /

ENTRYPOINT [ "/backend-linux-amd64" ]
