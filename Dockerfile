FROM node:22-alpine AS uibuild

WORKDIR /ui
COPY ui/package.json ui/package-lock.json ./
RUN npm ci
COPY ui/ .
RUN npm run build

FROM golang:1.25-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=uibuild /build/content/ content/
RUN CGO_ENABLED=0 GOOS=linux go build -o segdsp_worker .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /build/segdsp_worker .

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
ENV SQUELCH=-150
ENV SQUELCH_ALPHA=0.001

CMD ["./segdsp_worker"]
