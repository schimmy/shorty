FROM debian:jessie
COPY build/shorty /usr/bin/shorty
CMD ["shorty"]
