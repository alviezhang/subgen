FROM debian:12-slim

ENV ID=""
ENV CONFIG=""
ENV PORT=""

COPY ./subgen /usr/local/bin
ENTRYPOINT [subgen --id $ID --config $CONFIG --port $PORT] 