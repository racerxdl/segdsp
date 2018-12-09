"use strict";

const margin = 50;
const marginTop = 15;
const spacing = 60;

const waterFallLut = [];

let width, height;
let canvas;
let ctx;
let socket = null;
let fftData = [];
let url = '';
let waterFallHeight = 128;
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
let waterfallBuffers = [];

// From GQRX: https://github.com/csete/gqrx -> qtgui/plotter.cpp
for (let i = 0; i < 256; i++) {
    // level 0: black background
    if (i < 20)
        waterFallLut.push([0, 0, 0]);
    // level 1: black -> blue
    else if ((i >= 20) && (i < 70))
        waterFallLut.push([0, 0, 140*(i-20)/50]);
    // level 2: blue -> light-blue / greenish
    else if ((i >= 70) && (i < 100))
        waterFallLut.push([60*(i-70)/30, 125*(i-70)/30, 115*(i-70)/30 + 140]);
    // level 3: light blue -> yellow
    else if ((i >= 100) && (i < 150))
        waterFallLut.push([195*(i-100)/50 + 60, 130*(i-100)/50 + 125, 255-(255*(i-100)/50)]);
    // level 4: yellow -> red
    else if ((i >= 150) && (i < 250))
        waterFallLut.push([255, 255-255*(i-150)/100, 0]);
    // level 5: red -> white
    else if (i >= 250)
        waterFallLut.push([255, 255*(i-250)/5, 255*(i-250)/5]);
}

function toNotationUnit(v) {
    let unit;
    let submultiplo = ["","m","&micro;","n","p","f","a","z","y"];
    let multiplo    = ["","k","M","G","T","P","E","Z","Y"];
    let counter= 0;
    let value = v;
    if(value < 1) {
        while(value < 1) {
            counter++;
            value=value*1e3;
            if(counter === 8) break;
        }
        unit = submultiplo[counter];
    }else{
        while(value > 1000) {
            counter++;
            value=value/1e3;
            if(counter === 8) break;
        }
        unit = multiplo[counter];
    }
    value = Math.round(value*1e2)/1e2;
    return [value,unit];
}

function toHzNotation(v) {
    const z = toNotationUnit(v);
    return z[0].toLocaleString() + ' ' + z[1] + 'Hz';
}

function toBytesPerSecNotation(v) {
    const z = toNotationUnit(v);
    return z[0].toLocaleString() + ' ' + z[1] + 'b/s';
}

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
    const fftHeight = 256;
    waterFallHeight = height - 50 - fftHeight;
    // region Clear
    ctx.fillStyle = '#001111';
    ctx.fillRect(0, 0, width, height);
    // endregion
    // region Draw Grid
    const fftWidth = deviceInfo.DisplayPixels;
    const hDivs = CalcDiv(fftHeight) * 2;
    const vDivs = CalcDiv(fftWidth);
    const hDivDelta = fftHeight / hDivs;
    const vDivDelta = fftWidth / vDivs;
    const delta = deviceInfo.DisplayBandwidth / fftWidth;
    const invDelta = fftWidth / deviceInfo.DisplayBandwidth;
    const startFreq = deviceInfo.DisplayCenterFrequency - (deviceInfo.DisplayBandwidth / 2);

    // region Draw Frequency Labels
    ctx.beginPath();
    ctx.lineWidth = 1;
    ctx.strokeStyle = '#444444';
    for (let i = 0; i < fftWidth+1; i+=vDivDelta) {
        ctx.moveTo(margin + i, marginTop);
        ctx.lineTo(margin + i, marginTop + fftHeight + 5);
        const freq = (startFreq + i * delta) / 1e6;
        const freqStr = freq.toLocaleString();
        ctx.save();
        ctx.fillStyle = '#FFFFFF';
        const freqX = margin + i - ctx.measureText(freqStr).width / 2;
        ctx.fillText(freqStr, freqX, marginTop + fftHeight + 25);
        ctx.restore();
    }

    ctx.save();
    ctx.fillStyle = '#FFFFFF';
    const MHzText = 'MHz';
    ctx.fillText(MHzText, width - 50, marginTop + fftHeight + 10);
    ctx.restore();
    // endregion
    // region Draw dB Label
    const min = deviceInfo.DisplayOffset;
    const max = min - deviceInfo.DisplayRange;
    const dbPerPixel = (max - min) / fftHeight;

    for (let i = 0; i < fftHeight + 1; i+=hDivDelta) {
        ctx.moveTo(margin - 10, marginTop + i);
        ctx.lineTo(width - margin, marginTop + i);
        const dbLvl = (i * dbPerPixel) + min;
        const dbLvlStr = dbLvl.toFixed(0);
        ctx.save();
        ctx.fillStyle = '#FFFFFF';
        const dbLvlX = margin - ctx.measureText(dbLvlStr).width - 15;
        ctx.fillText(dbLvlStr, dbLvlX, marginTop + i + 5);
        ctx.restore();
    }
    ctx.stroke();
    ctx.closePath();
    // endregion
    // endregion
    // region Draw FFT
    ctx.beginPath();
    ctx.lineWidth = 2;
    ctx.strokeStyle = '#AAAAAA';
    ctx.moveTo(margin, marginTop + fftHeight - fftData[0]);
    for (let i = 1; i < fftWidth; i++) {
        ctx.lineTo(margin + i, marginTop + fftHeight - fftData[i]);
    }
    ctx.stroke();
    ctx.closePath();
    // endregion
    // region Draw Waterfall
    let waterFallOffset = marginTop + fftHeight + 30;
    for (let i = waterfallBuffers.length - 1; i >= 0; i--) {
        ctx.putImageData(waterfallBuffers[i], margin, waterFallOffset);
        waterFallOffset++;
    }
    // endregion
    // region Markers
    const channelX = (margin + (deviceInfo.ChannelCenterFrequency - startFreq) * invDelta) >> 0;
    const channelW = (deviceInfo.FilterBandwidth * invDelta) >> 0;
    ctx.fillStyle = 'rgba(127, 127, 127, 0.3)';
    ctx.fillRect(channelX - channelW / 2, marginTop, channelW, fftHeight);
    ctx.beginPath();
    ctx.lineWidth = 1;
    ctx.strokeStyle = '#AA0000';
    ctx.moveTo(channelX, marginTop);
    ctx.lineTo(channelX, marginTop + fftHeight);
    ctx.stroke();
    ctx.closePath();
    // endregion
}

function HandleFFT(data) {
    const z = atob(data);
    const buff = [];
    const wtfBuff = ctx.createImageData(data.length, 1);
    const wtfData = wtfBuff.data;

    for (let i = 0; i < z.length; i++) {
        const v = z.charCodeAt(i);
        buff.push(v);
        wtfData[i*4+0] = waterFallLut[v][0];
        wtfData[i*4+1] = waterFallLut[v][1];
        wtfData[i*4+2] = waterFallLut[v][2];
        wtfData[i*4+3] = 255;
    }

    waterfallBuffers.push(wtfBuff);

    if (waterfallBuffers.length > waterFallHeight) {
        waterfallBuffers.splice(0, 1);
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

    if (deviceInfo.connected) {
        document.title = deviceInfo.StationName;
        document.getElementById("headText").innerHTML = 'Connected to ' + deviceInfo.StationName + ' at ' + url + ' (' + deviceInfo.DeviceName + ')';
        document.getElementById("demodMode").innerHTML = deviceInfo.DemodulatorMode;
        document.getElementById("filterBw").innerHTML = toHzNotation(deviceInfo.FilterBandwidth);
        document.getElementById("channelBw").innerHTML = toHzNotation(deviceInfo.CurrentSampleRate);
        document.getElementById("fftFreq").innerHTML = toHzNotation(deviceInfo.DisplayCenterFrequency);
        document.getElementById("channelFreq").innerHTML = toHzNotation(deviceInfo.ChannelCenterFrequency);
    } else {
        // ctx.fillStyle = 'red';
        document.getElementById("headText").innerHTML = 'Connecting to ' + url;
    }

    DrawFFT();
    console.log(data);
}

function UpdateTraffic() {

    averageTraffic = trafficSum;
    trafficSum = 0;
    document.getElementById('avgTraffic').innerHTML = 'Avg. Traffic: '+ toBytesPerSecNotation(averageTraffic);

    setTimeout(UpdateTraffic, 1000);
}

function UpdateLevel(level) {
    document.getElementById('channelLevel').innerHTML = Math.round(level).toLocaleString() + ' dB';
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
                case 'fft': HandleFFT(data.FFTData); UpdateLevel(data.DemodOutputLevel); break;
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

