FROM golang:1.21.1-bullseye

RUN go install github.com/beego/bee/v2@latest

ENV GO111MODULE=on

ENV APP_HOME /go/src/mathapp
RUN mkdir -p "$APP_HOME"
RUN export GIN_MODE=release
RUN mkdir ../logs
WORKDIR "$APP_HOME"

EXPOSE 8000
CMD ["bee", "run", "-buildvcs=false"]