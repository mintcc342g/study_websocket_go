FROM golang:1.16-alpine
RUN mkdir /study_websocket_go
WORKDIR /study_websocket_go
ADD bin/study_websocket_go bin/study_websocket_go
ADD conf conf
ARG BUILD_PORT
ENV PORT $BUILD_PORT
EXPOSE $BUILD_PORT
ENTRYPOINT ["bin/study_websocket_go"]
