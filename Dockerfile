FROM ubuntu:16.04

COPY ./artifacts/lofai /usr/local/bin

EXPOSE 8000
ENTRYPOINT [ "/usr/local/bin/lofai" ]
