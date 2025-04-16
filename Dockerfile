FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend ./
RUN npm run build

FROM golang:1.24-alpine AS go-builder
WORKDIR /app

RUN apk add --no-cache git gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

RUN CGO_ENABLED=1 go build -o tempus ./main.go

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache sqlite

COPY --from=go-builder /app/tempus ./

EXPOSE 8080

CMD ["./tempus"]
