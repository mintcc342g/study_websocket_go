FROM golang:1.16-alpine
RUN mkdir /study-websocket-go
WORKDIR /study-websocket-go
ADD bin/study-websocket-go bin/study-websocket-go
ADD conf conf
ARG BUILD_PORT
ENV PORT $BUILD_PORT
EXPOSE $BUILD_PORT
ENTRYPOINT ["bin/study-websocket-go"]
