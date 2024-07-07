FROM --platform=linux/amd64 golang:alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main

WORKDIR /app
RUN cp /build/main .

FROM --platform=linux/amd64 alpine
COPY --from=builder /app/main .

ENTRYPOINT [ "/main" ]%