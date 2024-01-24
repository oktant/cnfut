FROM golang:1.20-alpine3.16

RUN apk add --no-cache git
WORKDIR /app/cnfut
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ./out/cnfut .
EXPOSE 8080
CMD ["./out/cnfut"]