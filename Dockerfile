FROM alpine:3.1
EXPOSE 80
COPY ./kubemonitor /usr/bin/
ENTRYPOINT ["kubemonitor"]