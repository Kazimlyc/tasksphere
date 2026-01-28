# ------------ BUILDER STAGE ------------
FROM golang:1.24 AS BUILDER

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# ------------ RUNNER STAGE ------------
FROM alpine:3.19

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
