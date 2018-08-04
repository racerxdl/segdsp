"use strict";

const margin = 50;
const spacing = 60;
const bottomSize = 150;

let width, height;
let canvas;
let ctx;
let socket = null;
let fftData = [];
let url = '';
let deviceInfo = {
    OutputRate: 48000,
    ChannelCenterFrequency: 0,
    CurrentSampleRate: 0,
    DeviceName: "None",
    DisplayCenterFrequency: 0,
    DisplayOffset: 0,
    DisplayPixels: 0,
    DisplayRange: 0,
    DisplayBandwidth: 0,
    FilterBandwidth: 0,
    DemodulatorMode: '',
    DemodulatorParams: null,
    Gain: 0,
    connected: false,
    StationName: "SegDSP",
    WebCanControl: true,
};

let averageTraffic = 0;
let trafficSum = 0;
let buffers = null;

function InitWebAudio(sampleRate) {
    console.log('Initializing WebAudio with SR: ' + sampleRate);
    const audioCtx = new (window.AudioContext || window.webkitAudioContext)({
        sampleRate: sampleRate,
    });

    const node = audioCtx.createScriptProcessor(16384, 0, 1);
    const source = audioCtx.createBufferSource();

    node.onaudioprocess = function(event) {
        try {
            if (buffers.length > 0) {
                const buff = buffers.shift();
                event.outputBuffer.getChannelData(0).set(buff);
            }
        } catch(e) {
            console.log(e);
        }
    };
    source.sampleRate = 48000;
    source.connect(node);
    node.connect(audioCtx.destination);
    source.start();
    buffers = [];
}

/**
 * @return {number}
 */
function CalcDiv(size) {
    return (size / spacing) >> 0;
}

function DrawFFT() {
    const baseOffset = height - bottomSize;
    // region Clear
    ctx.fillStyle = '#001111';
    ctx.fillRect(0, 0, width, height);
    // endregion
    // region Draw Grid
    const fftWidth = deviceInfo.DisplayPixels;
    const hDivs = CalcDiv(baseOffset) * 2;
    const vDivs = CalcDiv(fftWidth);
    const hDivDelta = baseOffset / hDivs;
    const vDivDelta = fftWidth / vDivs;
    // const deltaDb = baseOffset / (deviceInfo.DisplayRange / 256);
    const delta = deviceInfo.DisplayBandwidth / fftWidth;
    const invDelta = fftWidth / deviceInfo.DisplayBandwidth;
    const startFreq = deviceInfo.DisplayCenterFrequency - (deviceInfo.DisplayBandwidth / 2);
    const endFreq = deviceInfo.DisplayCenterFrequency + (deviceInfo.DisplayBandwidth / 2);

    ctx.beginPath();
    ctx.lineWidth = 1;
    ctx.strokeStyle = '#444444';
    for (let i = 0; i < fftWidth+1; i+=vDivDelta) {
        ctx.moveTo(margin + i, margin);
        ctx.lineTo(margin + i, baseOffset + 5);
        const freq = (startFreq + i * delta) / 1e6;
        const freqStr = freq.toLocaleString();
        ctx.save();
        ctx.fillStyle = '#FFFFFF';
        const freqX = margin + i - ctx.measureText(freqStr).width / 2;
        ctx.fillText(freqStr, freqX, baseOffset + 25);
        ctx.restore();
    }

    ctx.save();
    ctx.fillStyle = '#FFFFFF';
    const MHzText = 'MHz';
    ctx.fillText(MHzText, width - 50, baseOffset + 40);
    ctx.restore();

    const min = -deviceInfo.DisplayOffset;
    const max = min - deviceInfo.DisplayRange;
    const range = (baseOffset) - margin;
    const dbPerPixel = (max - min) / range;

    for (let i = margin; i < baseOffset+1; i+=hDivDelta) {
        ctx.moveTo(margin - 10, i);
        ctx.lineTo(width - margin, i);
        const z = i - margin;
        const dbLvl = (z * dbPerPixel) + min;
        const dbLvlStr = dbLvl.toFixed(0);
        ctx.save();
        ctx.fillStyle = '#FFFFFF';
        const dbLvlX = margin - ctx.measureText(dbLvlStr).width - 15;
        ctx.fillText(dbLvlStr, dbLvlX, i + 5);
        ctx.restore();
    }
    ctx.stroke();
    ctx.closePath();
    // endregion
    // region Draw FFT
    ctx.beginPath();
    ctx.lineWidth = 2;
    ctx.strokeStyle = '#AAAAAA';
    ctx.moveTo(margin, baseOffset - fftData[0]);
    for (let i = 1; i < fftWidth; i++) {
        ctx.lineTo(margin + i, baseOffset - fftData[i]);
    }
    ctx.stroke();
    ctx.closePath();
    // endregion
    // region Draw Texts
    ctx.font = "15px Arial";

    if (deviceInfo.connected) {
        ctx.fillStyle = 'green';
        ctx.fillText('Connected to ' + deviceInfo.StationName + ' at ' + url + ' (' + deviceInfo.DeviceName + ')', 10, 20);
        ctx.fillStyle = 'white';
        ctx.fillText('Demodulator Mode: ' + deviceInfo.DemodulatorMode, 10, height - 64);
        ctx.fillText('Channel BW: ' + deviceInfo.CurrentSampleRate.toLocaleString() + ' Hz', 10, height - 46);
        ctx.fillText('FFT Center Frequency: ' + deviceInfo.DisplayCenterFrequency.toLocaleString() + ' Hz', 10, height - 28);
        ctx.fillText('Channel Center Frequency: ' + deviceInfo.ChannelCenterFrequency.toLocaleString() + ' Hz', 10, height - 10);
        ctx.fillText('Avg. Traffic: '+ (averageTraffic / 1024).toFixed(2) + " kb/s", width - 170, height - 10);
    } else {
        ctx.fillStyle = 'red';
        ctx.fillText('Connecting to ' + url, 10, 20);
    }
    // endregion
    // region Markers
    const channelX = (margin + (deviceInfo.ChannelCenterFrequency - startFreq) * invDelta) >> 0;
    const channelW = (deviceInfo.FilterBandwidth * invDelta) >> 0;
    ctx.fillStyle = 'rgba(127, 127, 127, 0.3)';
    ctx.fillRect(channelX - channelW / 2, margin, channelW, baseOffset - margin);
    ctx.beginPath();
    ctx.lineWidth = 1;
    ctx.strokeStyle = '#AA0000';
    ctx.moveTo(channelX, margin);
    ctx.lineTo(channelX, baseOffset);
    ctx.stroke();
    ctx.closePath();
    // endregion
}

function HandleFFT(data) {
    const z = atob(data);
    const buff = [];

    for (let i = 0; i < z.length; i++) {
        buff.push(z.charCodeAt(i));
    }

    fftData = buff;
    DrawFFT();
}

function HandleData(data) {
    // console.log('Received buffer!');
    try {
        const buffer = data.Data;
        const audioRate = data.OutputRate;
        if (buffers === null) {
            InitWebAudio(audioRate);
        }
        // audioCtx.sampleRate = audioRate;
        buffers.push(new Float32Array(buffer));
    } catch (e) {
        console.log(e);
    }
}

function HandleRawData(data) {
    if (buffers !== null) {
        let fileReader = new FileReader();
        let buff;
        fileReader.onload = function(event) {
            buff = event.target.result;
            trafficSum += buff.byteLength;
            buffers.push(new Float32Array(buff));
        };
        fileReader.readAsArrayBuffer(data);
    }
}

function addClass(el, cls) {
    let arr = el.className.split(" ");
    if (arr.indexOf(cls) === -1) {
        el.className += " " + cls;
    }
}

function removeClass(el, cls) {
    let arr = el.className.split(" ");
    if (arr.indexOf(cls) !== -1) {
        arr.splice(arr.indexOf(cls), 1);
        el.className = arr.join(' ');
    }
}

function HandleDevice(data) {
    deviceInfo = data;
    deviceInfo.connected = true;

    if (!deviceInfo.WebCanControl) {
        addClass(document.getElementById("contentDiv"), 'lockedBorder');
        addClass(document.getElementById("lockedImg"), 'locked');
        removeClass(document.getElementById("lockedImg"), 'unlocked');
    } else {
        removeClass(document.getElementById("contentDiv"), 'lockedBorder');
        removeClass(document.getElementById("lockedImg"), 'locked');
        addClass(document.getElementById("lockedImg"), 'unlocked');
    }

    InitWebAudio(data.OutputRate);
    canvas.width = deviceInfo.DisplayPixels + margin * 2;
    width = deviceInfo.DisplayPixels + margin * 2;
    document.getElementById("contentDiv").style.maxWidth = width + "px";
    document.getElementById("contentDiv").style.width = width + "px";
    DrawFFT();
    console.log(data);
}

function UpdateTraffic() {

    averageTraffic = trafficSum;
    trafficSum = 0;

    setTimeout(UpdateTraffic, 1000);
}

function Connect() {
    // Websocket
    UpdateTraffic();
    const proto = location.protocol === 'https:' ? 'wss://' : 'ws://';
    url = proto + location.host + "/ws";
    socket = new WebSocket(url);
    socket.onopen = (evt) => {
        console.log('Connected!')
    };
    socket.onclose = (evt) => {
        console.log('Connection closed!');
        socket = null;
        DrawFFT();
        setTimeout(Connect, 1500);
    };
    socket.onmessage = (evt) => {
        // console.log('Received message: ' + evt.data);
        if (typeof evt.data !== 'string') {
            HandleRawData(evt.data);
            return;
        }
        trafficSum += evt.data.length;
        try {
            const data = JSON.parse(evt.data);
            switch (data.MessageType) {
                case 'fft': HandleFFT(data.FFTData); break;
                case 'data': HandleData(data.Data); break;
                case 'device': HandleDevice(data); break;
                default: console.log('Unknown Type: ' + data.MessageType);
            }
        } catch (e) {
            console.log('Error parsing json: ' + e);
        }
    };
    socket.onerror = (evt) => {
        console.log('Error on socket: ' + evt.data);
    };
    DrawFFT();
}

function Init() {
    canvas = document.getElementById("fft");
    ctx = canvas.getContext("2d");
    width = canvas.width;
    height = canvas.height;

    Connect();
}

