FROM debian:jessie
ADD static/ /var/www
WORKDIR /var
COPY bin/shorty /usr/bin/shorty
CMD ["shorty"]
