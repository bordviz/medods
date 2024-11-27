FROM golang:1.23

RUN mkdir /medods
WORKDIR /medods

COPY . .
RUN chmod a+x docker/*.sh
ENTRYPOINT ["/medods/docker/app.sh"]