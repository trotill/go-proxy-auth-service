FROM golang:1.23.0-alpine as build

ARG BRANCH

RUN apk update
RUN apk add --no-cache curl git make
RUN git clone https://github.com/trotill/go-proxy-auth-service.git /opt/auth --branch $BRANCH&&\
    cd /opt/auth&&\
    make tidy&&\
    make build

FROM alpine:3.21
COPY --from=build /opt/auth/bin/authProxy /app

ENTRYPOINT ["/app"]
CMD /app/authProxy