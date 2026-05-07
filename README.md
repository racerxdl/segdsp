# SegDSP

Go SDR demodulator. Connects to a [radioserver](https://github.com/racerxdl/radioserver) instance for IQ data and provides a web-based spectrum display with demodulated audio output.

[![CI](https://github.com/racerxdl/segdsp/actions/workflows/ci.yml/badge.svg)](https://github.com/racerxdl/segdsp/actions/workflows/ci.yml) [![Apache License](https://img.shields.io/badge/license-Apache-blue.svg)](https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)) [![Go Report](https://goreportcard.com/badge/github.com/racerxdl/segdsp)](https://goreportcard.com/report/github.com/racerxdl/segdsp)

## Features

- Wideband and narrowband FM demodulation with de-emphasis
- AM demodulation with AGC
- Real-time FFT waterfall via WebSocket
- Squelch with automatic recording
- Cross-platform single binary (static content embedded via `go:embed`)
- Multi-arch Docker images

## Quick Start

```bash
go build -o segdsp .

# FM radio
segdsp -channelFrequency 106300000 -demodMode FM -fmDeviation 75000 -filterBandwidth 120000 -radioserver localhost:4050

# AM radio
segdsp -channelFrequency 145570000 -demodMode AM -filterBandwidth 15000 -radioserver localhost:4050
```

**Requires** a running [radioserver](https://github.com/racerxdl/radioserver) instance.

## Docker

Images are published to Docker Hub on every push to `master`:

```bash
docker run -e RADIOSERVER=host.docker.internal:4050 racerxdl/segdsp
```

Available tags: `latest`, `amd64-latest`, `arm64v8-latest`, `arm32v6-latest`

## Presets

Instead of setting individual flags, use a preset:

```bash
segdsp -preset wbfm -channelFrequency 106300000
segdsp -preset nbfm -channelFrequency 145570000
segdsp -preset am   -channelFrequency 145570000
```

## Configuration

| Argument              | Environment variable    | Type   | Default    | Description                                                       |
|-----------------------|-------------------------|--------|------------|-------------------------------------------------------------------|
| `-channelFrequency`   | `CENTER_FREQUENCY`      | number | 106300000  | Channel (IQ) Center Frequency in Hz                               |
| `-radioserver`        | `RADIOSERVER`           | string | localhost:4050 | radioserver address                                           |
| `-demodMode`          | `DEMOD_MODE`            | string | FM         | Demodulator mode: `FM`, `AM`                                      |
| `-preset`             | `PRESET`                | string | none       | Preset: `wbfm`, `nbfm`, `am`                                      |
| `-outputRate`         | `OUTPUT_RATE`           | number | 48000      | Output audio sample rate in Hz                                     |
| `-filterBandwidth`    | `FS_BANDWIDTH`          | number | 120000     | First stage filter bandwidth in Hz                                 |
| `-decimationStage`    | `DECIMATION_STAGE`      | number | 3          | Channel decimation stage (decimation = 2^d)                        |
| `-fftFrequency`       | `FFT_FREQUENCY`         | number | 106300000  | FFT center frequency in Hz                                         |
| `-fftDecimationStage` | `FFT_DECIMATION_STAGE`  | number | 1          | FFT decimation stage                                               |
| `-displayPixels`      | `DISPLAY_PIXELS`        | number | 512        | FFT display width in pixels                                        |
| `-httpAddr`           | `HTTP_ADDRESS`          | string | localhost:8080 | HTTP service address                                           |
| `-fmDeviation`        | `FM_DEVIATION`          | number | 75000      | FM max deviation in Hz                                             |
| `-fmTau`              | `FM_TAU`                | number | 0.000075   | FM de-emphasis tau (0 to disable)                                  |
| `-squelch`            | `SQUELCH`               | number | -150       | Squelch threshold in dB                                            |
| `-squelchAlpha`       | `SQUELCH_ALPHA`         | number | 0.001      | Squelch filter alpha                                               |
| `-amAudioCut`         | `AM_AUDIO_CUT`          | number | 5000       | AM audio low-pass cut in Hz                                        |
| `-stationName`        | `STATION_NAME`          | string | SegDSP     | Station name / callsign                                            |
| `-record`             | `RECORD`                | bool   | false      | Record audio when not squelched                                    |
| `-recordMethod`       | `RECORD_METHOD`         | string | file       | Recording method (`file`)                                          |
| `-webCanControl`      | `WEB_CAN_CONTROL`       | bool   | false      | Allow web UI to control settings                                   |
| `-tcpCanControl`      | `TCP_CAN_CONTROL`       | bool   | true       | Allow TCP clients to control settings                               |
| `-cpuprofile`         |                         | string |            | Write CPU profile to file                                          |

## Development

```bash
go build -o segdsp .
go test -v -race ./...
golangci-lint run
```

## License

Apache 2.0
