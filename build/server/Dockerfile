FROM golang:1.18 AS dependencies
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download

FROM dependencies AS build
COPY . /go/src/app
WORKDIR /go/src/app
RUN make build-server

FROM debian:buster-slim
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/app/bin/* /app/
COPY --from=build /go/src/app/scripts/wait-for-it.sh /app/
RUN chmod +x /app/*
EXPOSE 50051/tcp
CMD /app/gk serve
