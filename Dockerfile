FROM centos:7
LABEL author="yashmankad"

RUN mkdir -p /opt/test_service
WORKDIR /opt/test_service
COPY bin/test_service .
COPY config/config.yml .

ENTRYPOINT ./test_service -c=/opt/test_service/config.yml
EXPOSE 8000
EXPOSE 8001