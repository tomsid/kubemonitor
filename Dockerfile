FROM alpine:3.2
EXPOSE 80
COPY ./kubemonitor /usr/bin/
ENTRYPOINT ["kubemonitor"]
