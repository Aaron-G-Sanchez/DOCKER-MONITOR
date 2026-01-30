FROM golang:alpine

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-monitor ./cmd/server

CMD [ "/docker-monitor" ]
