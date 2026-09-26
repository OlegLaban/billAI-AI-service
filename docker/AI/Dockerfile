FROM alpine:3.24.2

RUN apk add --no-cache curl git ollama util-linux
RUN echo '' > "/tmp/test.txt"

CMD ["ollama", "serve"]
