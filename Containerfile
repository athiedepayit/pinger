FROM golang:alpine
# add a few packages in case you need them, I guess?
RUN apk add -U --no-cache curl netcat-openbsd git
RUN git clone https://gitlab.com/adamthiede/pinger && cd pinger && go build && mv ./pinger /pinger
ENTRYPOINT ["/pinger"]
