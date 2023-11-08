FROM golang:1.21.3-alpine

WORKDIR /app

ENV GO111MODULE=on
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY . ./

RUN go build -o /GatewayService

EXPOSE 10000

CMD [ "/GatewayService" ]