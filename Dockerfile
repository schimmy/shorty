FROM google/debian:wheezy

ADD ./static /root/shortener/static
COPY build/shortener /root/shortener/shortener

EXPOSE 80
WORKDIR /root/shortener
CMD ["/root/shortener/shortener"]
