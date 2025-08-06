ARG IMAGE=ubuntu:22.04
FROM ${IMAGE}

RUN apt-get update && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

ADD ./service /service
ADD ./service-cli /service-cli
ADD ./openapi3.json /gen/http/openapi3.json

EXPOSE 8000 8080

ENTRYPOINT ["/service"]
