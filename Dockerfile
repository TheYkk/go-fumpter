FROM golang:1.17-alpine as builder_fumpt
RUN go install mvdan.cc/gofumpt@latest

FROM golang:1.17-alpine as  builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /bin/gobot

FROM alpine:3.15
RUN apk add git
WORKDIR /app
COPY --from=builder_fumpt /go/bin/gofumpt /app/gofumpt
COPY --from=builder /bin/deployment /app/gobot
ENTRYPOINT ["/app/gobot"]