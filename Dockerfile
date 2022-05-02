FROM alpine:latest
RUN apk add --no-cache go

COPY . .

ENTRYPOINT ["/entrypoint.sh"]
