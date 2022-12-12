FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /user_service ./cmd/main.go

ENV POSTGRES_PASSWORD=$POSTGRES_PASSWORD
ENV POSTGRES_HOST=$POSTGRES_HOST
ENV POSTGRES_USER=$POSTGRES_USER
ENV POSTGRES_PORT=$POSTGRES_PORT
ENV POSTGRES_DB_NAME=$POSTGRES_DB_NAME

EXPOSE 6001

CMD [ "/user_service" ]