FROM golang:1.20-alpine

ENV ID=""
ENV CONFIG=""
ENV PORT=""

RUN mkdir /workspace
WORKDIR /workspace
COPY .  .
RUN go build
ENTRYPOINT /workspace/subgen --id $ID --config $CONFIG --port $PORT 
