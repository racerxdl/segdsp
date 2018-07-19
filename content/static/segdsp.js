"use strict";

const margin = 50;
const spacing = 50;
const audioCtx = new (window.AudioContext || window.webkitAudioContext)();

let width, height;
let canvas;
let ctx;
let socket = null;
let fftData = [];
let url = '';
let deviceInfo = {
    ChannelCenterFrequency: 0,
    CurrentSampleRate: 0,
    DeviceName: "None",
    DisplayCenterFrequency: 0,
    DisplayOffset: 0,
    DisplayPixels: 0,
    DisplayRange: 0,
    DisplayBandwidth: 0,
    Gain: 0,
};

let audioBuffer;
let audioSource;
const buffers = [];

const node = audioCtx.createScriptProcessor(16384, 0, 1);
const source = audioCtx.createBufferSource();

node.onaudioprocess = function(event) {
    try {
        if (buffers.length > 0) {
            const buff = buffers.shift();
            event.outputBuffer.getChannelData(0).set(buff);
        } else {
           // console.log("Empty");
        }
    } catch(e) {
        console.log(e);
    }
};
source.connect(node);
node.connect(audioCtx.destination);
source.start();
/**
 * @return {number}
 */
function CalcDiv(size) {
    return (size / spacing) >> 0;
}

function DrawFFT() {
    const baseOffset = height - 100;
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

    for (let i = margin; i < baseOffset+1; i+=hDivDelta) {
        ctx.moveTo(margin - 10, i);
        ctx.lineTo(width - margin, i);
        // const dbLvl = deviceInfo.DisplayOffset + deltaDb * (baseOffset - i);
        // const dbLvlStr = dbLvl.toLocaleString();
        // ctx.save();
        // ctx.fillStyle = '#FFFFFF';
        // const dbLvlX = margin - ctx.measureText(dbLvlStr).width - 15;
        // ctx.fillText(dbLvlStr, dbLvlX, i + 5);
        // ctx.restore();
    }
    ctx.stroke();
    ctx.closePath();
    // endregion
    // region Draw FFT
    ctx.beginPath();
    ctx.lineWidth = 2;
    ctx.strokeStyle = '#AAAAAA';
    ctx.moveTo(margin, baseOffset - fftData[0]);
    for (let i = 0; i < fftWidth; i++) {
        ctx.lineTo(margin + i, baseOffset - fftData[i]);
    }
    ctx.stroke();
    ctx.closePath();
    // endregion
    // region Draw Texts
    ctx.font = "15px Arial";

    if (socket !== null) {
        ctx.fillStyle = 'green';
        ctx.fillText('Connected to ' + url + ' (' + deviceInfo.DeviceName + ')', 10, 20);
        ctx.fillStyle = 'white';
        ctx.fillText('Channel BW: ' + deviceInfo.CurrentSampleRate.toLocaleString() + ' Hz', 10, height - 46);
        ctx.fillText('FFT Center Frequency: ' + deviceInfo.DisplayCenterFrequency.toLocaleString() + ' Hz', 10, height - 28);
        ctx.fillText('Channel Center Frequency: ' + deviceInfo.ChannelCenterFrequency.toLocaleString() + ' Hz', 10, height - 10);
    } else {
        ctx.fillStyle = 'red';
        ctx.fillText('Connecting to ' + url, 10, 20);
    }
    // endregion
    // region Markers
    const channelX = (margin + (deviceInfo.ChannelCenterFrequency - startFreq) * invDelta) >> 0;
    const channelW = ((deviceInfo.CurrentSampleRate * 0.8 * 0.5) * invDelta) >> 0;
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
    fftData = data;
    DrawFFT();
}

function HandleData(data) {
    // console.log('Received buffer!');
    try {
        const buffer = data.Data;
        const audioRate = data.OutputRate;
        buffers.push(new Float32Array(buffer));
    } catch (e) {
        console.log(e);
    }
}

function HandleDevice(data) {
    console.log(data);
    deviceInfo = data;
    canvas.width = deviceInfo.DisplayPixels + margin * 2;
    width = deviceInfo.DisplayPixels + margin * 2;
    document.getElementById("contentDiv").style.maxWidth = width + "px";
    document.getElementById("contentDiv").style.width = width + "px";
    DrawFFT();
}

function Connect() {
    // Websocket
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

