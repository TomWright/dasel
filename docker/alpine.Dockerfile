FROM alpine

ARG daselpath=./dasel

WORKDIR /root
COPY $daselpath /usr/local/bin/dasel
RUN chmod +x /usr/local/bin/dasel

ENTRYPOINT ["/usr/local/bin/dasel"]
CMD []
