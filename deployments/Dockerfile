FROM golang:1.23.0-alpine as build

ARG BRANCH

RUN apk update
RUN apk add --no-cache curl git make
RUN git clone https://github.com/trotill/docker-black-hole.git /opt/dbh --branch $BRANCH&&\
    cd /opt/dbh&&\
    go install github.com/swaggo/swag/cmd/swag@v1.16.4&&\
    make tidy&&\
    make swagger&&\
    make build

FROM alpine:3.21
COPY --from=build /opt/dbh/bin/blackHole /app

ENTRYPOINT ["/app"]
CMD /app/blackHole