FROM alpine:3.22.1@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1

COPY .go-version .go-version

RUN apk add --no-cache go=$(cat .go-version)-r0

COPY . .

ENTRYPOINT ["/entrypoint.sh"]
