FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git for fetching dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install LibreOffice for DOCX conversion
RUN apk add --no-cache libreoffice ttf-dejavu

COPY --from=builder /app/main .
COPY --from=builder /app/static ./static

EXPOSE 19080

CMD ["./main"]
