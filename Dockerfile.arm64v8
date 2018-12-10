FROM arm64v8/golang:alpine as build

COPY qemu-aarch64-static /usr/bin/

RUN apk update

RUN apk add git ca-certificates gcc musl-dev

ADD . /go/src/github.com/racerxdl/segdsp

WORKDIR /go/src/github.com/racerxdl/segdsp

RUN go get -v

RUN CGO_ENABLED=0 GOOS=linux go build -o segdsp_worker

FROM arm64v8/alpine:latest

COPY qemu-aarch64-static /usr/bin/

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/github.com/racerxdl/segdsp/segdsp_worker .
COPY content content

ENV RADIOSERVER=localhost:4050
ENV CENTER_FREQUENCY=106300000
ENV FFT_FREQUENCY=106300000
ENV HTTP_ADDRESS=localhost:8080
ENV DISPLAY_PIXELS=512
ENV DECIMATION_STAGE=3
ENV FFT_DECIMATION_STAGE=0
ENV OUTPUT_RATE=48000
ENV DEMOD_MODE=FM
ENV FS_BANDWIDTH=120000
ENV STATION_NAME=SegDSP
ENV WEB_CAN_CONTROL=false
ENV TCP_CAN_CONTROL=true
ENV RECORD=false
ENV RECORD_METHOD=file
ENV PRESET=none

ENV FM_DEVIATION=75000
ENV FM_TAU=0.000075
ENV FM_SQUELCH=-72
ENV FM_SQUELCH_ALPHA=0.001

CMD ["./segdsp_worker"]

