FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the Go application
# -o /app/main specifies the output file name and location.
# CGO_ENABLED=0 disables Cgo to produce a statically linked binary.
# -ldflags "-s -w" strips debugging information, reducing the binary size.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" -o /app/main .

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3000

CMD ["/app/main"]