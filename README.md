[![Build Status](https://api.travis-ci.org/racerxdl/segdsp.svg?branch=master)](https://travis-ci.org/racerxdl/segdsp) [![Apache License](https://img.shields.io/badge/license-Apache-blue.svg)](https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)) [![Go Report](https://goreportcard.com/badge/github.com/racerxdl/segdsp)](https://goreportcard.com/report/github.com/racerxdl/segdsp)

# SegDSP - Future of RF Monitoring (WIP)


## Docker Images

All `master` branch builds is uploaded to dockerhub with the name `racerxdl/segdsp`. The following archs are available:

- `x86` - `racerxdl/segdsp:latest` - For any x86 machine
- `amd64` - `racerxdl/segdsp:amd64-latest` - For 64 bit x86 machines
- `arm32v6` - `racerxdl/segdsp:arm32v6-latest` - For 32 bit raspberry pies and equivalents
- `arm64v8` - `racerxdl/segdsp:arm64v8-latest` - For 64 bit raspberry pies and equivalents

## Binary Release

Please check the [Releases](https://github.com/racerxdl/segdsp/releases) for binary releases.

## Nice Sample using SegDSP as a Library

Here is a nice sample using segdsp as a library: https://github.com/racerxdl/segdsp-sample

## Running

SegDSP is pretty straightforward to run if you know what you want to capture. It accepts both Environment Variables (suitable for docker containers) or just normal command line arguments.

## Examples

### WBFM Demodulator

```bash
# Argument Mode
segdsp -channelFrequency 106300000 -demodMode FM -fmDeviation 75000 -filterBandwidth 120000 -fftFrequency 106300000 -decimationStage 3 -stationName PU2NVX

# Environment Mode
CENTER_FREQUENCY="106300000" DEMOD_MODE="FM" FM_DEVIATION="75000" FFT_FREQUENCY="106300000" FS_BANDWIDTH="120000" DECIMATION_STAGE="3" STATION_NAME="PU2NVX" segdsp
```

### NBFM Demodulator

```bash
# Argument Mode
segdsp -channelFrequency 145570000 -demodMode FM -fmDeviation 5000 -filterBandwidth 15000 -fftFrequency 145570000 -decimationStage 5 -stationName PU2NVX

# Environment Mode
CENTER_FREQUENCY="145570000" DEMOD_MODE="FM" FM_DEVIATION="5000" FS_BANDWIDTH="15000" FFT_FREQUENCY="145570000" DECIMATION_STAGE="5" STATION_NAME="PU2NVX" segdsp
```

## Arguments

| Argument              | Environment variable    | Type   | Possible Values  | Description                                                       | Default Value   |
|-----------------------|-------------------------|--------|------------------|-------------------------------------------------------------------|-----------------|
| `-channelFrequency`   | `CENTER_FREQUENCY`      | number |                  | Channel (IQ) Center Frequency in Hz                               | 106300000       |
| `-cpuprofile`         |                         | string |                  | Write cpu profile to specified file                               |                 |
| `-decimationStage`    | `DECIMATION_STAGE`      | number |                  | Channel (IQ) Decimation Stage (The actual decimation will be 2^d) | 3               |
| `-demodMode`          | `DEMOD_MODE`            | string | `FM`, `AM`       | Demodulator Mode: [FM]                                            | FM              |
| `-displayPixels`      | `DISPLAY_PIXELS`        | number |                  | Width in pixels of the FFT                                        | 512             |
| `-fftDecimationStage` | `FFT_DECIMATION_STAGE`  | number |                  | FFT Decimation Stage (The actual decimation will be 2^d)          | 0               |
| `-fftFrequency`       | `FFT_FREQUENCY`         | number |                  | FFT Center Frequency in Hz                                        | 106300000       |
| `-filterBandwidth`    | `FS_BANDWIDTH`          | number |                  | First Stage Filter Bandwidth in Hert                              | 120000          |
| `-squelch`            | `SQUELCH`               | number |                  | Demodulator Squelch in dB                                         | -72             |
| `-squelchAlpha`       | `SQUELCH_ALPHA`         | number |                  | Demodulator Squelch Filter Alpha                                  | 0.001           |
| `-fmDeviation`        | `FM_DEVIATION`          | number |                  | FM Demodulator Max Deviation in Hertz                             | 75000           |
| `-fmTau`              | `FM_TAU`                | number |                  | FM Demodulator Tau in seconds (0 to disable)                      | 0.0000075       |
| `-amAudioCut`         | `AM_AUDIO_CUT`          | number |                  | AM Demodulator Audio Low Pass Cut                                 | 5000            |
| `-httpAddr`           | `HTTP_ADDRESS`          | string |                  | HTTP Service Address                                              | localhost:8080  |
| `-outputRate`         | `OUTPUT_RATE`           | number |                  | Output Rate in Hz                                                 | 48000           |
| `-record`             | `RECORD`                |  bool  | `true`, `false`  | If it should record output when not squelched                     | false           |
| `-recordMethod`       | `RECORD_METHOD`         | string | `file`           | Method to use when recording                                      | file            |
| `-radioserver`        | `RADIOSERVER`           | string |                  | radioserver Address                                                 | localhost:5555  |
| `-stationName`        | `STATION_NAME`          | string |                  | Name of the Station                                               | SegDSP          |

## Git Hooks

### pre-commit

Inside the project repository create the following file:

```bash
touch .git/hooks/pre-commit
```

And paste this code inside it:

```bash
#!/bin/bash

echo "Formatting code"
go fmt ./...
exit 0
```

And give execute permission:

```bash
chmod +x .git/hooks/pre-commit
```

