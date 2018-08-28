FROM alpine:3.7
LABEL author="sha1n"
LABEL description="an alpine based crazed server used for testing purposes"

RUN apk update

# Install bash for kube console
RUN apk add bash
RUN apk add curl

# Installing ca-certs for HTTPS client support
RUN apk add ca-certificates
RUN update-ca-certificates

COPY bin/dummy-loader /app/dummy-loader
COPY web/help.tpl /app/server/web/help.tpl
COPY scripts /app/scripts

WORKDIR /app

RUN chmod +x dummy-loader

EXPOSE 8080

CMD ["/app/dummy-loader"]

