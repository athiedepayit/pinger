FROM golang:alpine AS builder
RUN apk add git && git clone https://gitlab.com/adamthiede/pinger && cd pinger && go build && mv ./pinger /pinger

FROM alpine:latest
RUN apk add -U --no-cache curl netcat-openbsd git
COPY --from=builder /pinger /bin/pinger

ENTRYPOINT ["pinger"]
