FROM ubuntu:16.04

COPY lofai /usr/local/bin

EXPOSE 8000
ENTRYPOINT [ "/usr/local/bin/lofai" ]
