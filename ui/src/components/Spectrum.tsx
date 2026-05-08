import { useRef, useEffect, useCallback } from 'preact/hooks';
import { waterfallLut } from '../services/colormap';
import type { DeviceInfo, FFTMessage } from '../services/ws';

interface Props {
  fft: FFTMessage | null;
  device: DeviceInfo | null;
  onTune?: (freqHz: number) => void;
}

const MARGIN = 50;
const MARGIN_TOP = 15;
const FFT_HEIGHT = 256;
const SPACING = 60;

function toHz(v: number): string {
  const units = ['', 'k', 'M', 'G', 'T'];
  let val = v;
  let idx = 0;
  while (val > 1000 && idx < units.length - 1) {
    val /= 1000;
    idx++;
  }
  return val.toFixed(idx === 0 ? 0 : 1) + ' ' + units[idx] + 'Hz';
}

export function Spectrum({ fft, device, onTune }: Props) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const wfBuffers = useRef<ImageData[]>([]);
  const peakHold = useRef<Float32Array | null>(null);
  const peakDecay = useRef<Float32Array | null>(null);

  const draw = useCallback(() => {
    const canvas = canvasRef.current;
    if (!canvas || !device) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const fftWidth = device.DisplayPixels;
    const wfHeight = canvas.height - 50 - FFT_HEIGHT;
    const width = fftWidth + MARGIN * 2;

    ctx.fillStyle = '#0a0e14';
    ctx.fillRect(0, 0, width, canvas.height);

    const hDivs = Math.floor(FFT_HEIGHT / SPACING) * 2;
    const vDivs = Math.floor(fftWidth / SPACING);
    const hDivDelta = FFT_HEIGHT / hDivs;
    const vDivDelta = fftWidth / vDivs;
    const delta = device.DisplayBandwidth / fftWidth;
    const invDelta = fftWidth / device.DisplayBandwidth;
    const startFreq = device.DisplayCenterFrequency - device.DisplayBandwidth / 2;

    ctx.beginPath();
    ctx.lineWidth = 1;
    ctx.strokeStyle = '#1a2332';
    for (let i = 0; i <= fftWidth; i += vDivDelta) {
      ctx.moveTo(MARGIN + i, MARGIN_TOP);
      ctx.lineTo(MARGIN + i, MARGIN_TOP + FFT_HEIGHT + 5);
      const freq = (startFreq + i * delta) / 1e6;
      ctx.fillStyle = '#8892a0';
      ctx.font = '11px monospace';
      ctx.fillText(freq.toFixed(1), MARGIN + i - 15, MARGIN_TOP + FFT_HEIGHT + 20);
    }

    const min = device.DisplayOffset;
    const max = min - device.DisplayRange;
    const dbPerPixel = (max - min) / FFT_HEIGHT;

    for (let i = 0; i <= FFT_HEIGHT; i += hDivDelta) {
      ctx.moveTo(MARGIN - 10, MARGIN_TOP + i);
      ctx.lineTo(width - MARGIN, MARGIN_TOP + i);
      const dbLvl = i * dbPerPixel + min;
      ctx.fillStyle = '#8892a0';
      ctx.font = '11px monospace';
      ctx.fillText(dbLvl.toFixed(0), 2, MARGIN_TOP + i + 4);
    }
    ctx.stroke();

    ctx.fillStyle = '#8892a0';
    ctx.font = '11px monospace';
    ctx.fillText('MHz', width - 40, MARGIN_TOP + FFT_HEIGHT + 10);

    if (fft) {
      const raw = atob(fft.FFTData);
      const len = raw.length;
      const fftData = new Float32Array(len);
      for (let i = 0; i < len; i++) {
        fftData[i] = raw.charCodeAt(i);
      }

      if (!peakHold.current || peakHold.current.length !== len) {
        peakHold.current = new Float32Array(len);
        peakDecay.current = new Float32Array(len).fill(0);
      }

      const peak = peakHold.current;
      const decay = peakDecay.current;
      for (let i = 0; i < len; i++) {
        if (fftData[i] >= peak[i]) {
          peak[i] = fftData[i];
          decay[i] = 0;
        } else {
          decay[i] += 0.3;
          peak[i] = Math.max(fftData[i], peak[i] - decay[i]);
        }
      }

      ctx.beginPath();
      ctx.lineWidth = 1.5;
      ctx.strokeStyle = '#2d5aa0';
      for (let i = 0; i < len; i++) {
        const x = MARGIN + i;
        const y = MARGIN_TOP + FFT_HEIGHT - peak[i];
        if (i === 0) ctx.moveTo(x, y);
        else ctx.lineTo(x, y);
      }
      ctx.stroke();

      ctx.beginPath();
      ctx.lineWidth = 1.5;
      ctx.strokeStyle = '#4fc3f7';
      for (let i = 0; i < len; i++) {
        const x = MARGIN + i;
        const y = MARGIN_TOP + FFT_HEIGHT - fftData[i];
        if (i === 0) ctx.moveTo(x, y);
        else ctx.lineTo(x, y);
      }
      ctx.stroke();

      const wfBuf = ctx.createImageData(len, 1);
      const wfData = wfBuf.data;
      for (let i = 0; i < len; i++) {
        const v = Math.round(fftData[i]);
        const ci = Math.max(0, Math.min(255, v));
        wfData[i * 4 + 0] = waterfallLut[ci][0];
        wfData[i * 4 + 1] = waterfallLut[ci][1];
        wfData[i * 4 + 2] = waterfallLut[ci][2];
        wfData[i * 4 + 3] = 255;
      }
      wfBuffers.current.push(wfBuf);
      if (wfBuffers.current.length > wfHeight) {
        wfBuffers.current.splice(0, wfBuffers.current.length - wfHeight);
      }
    }

    let wfOffset = MARGIN_TOP + FFT_HEIGHT + 30;
    for (let i = wfBuffers.current.length - 1; i >= 0; i--) {
      ctx.putImageData(wfBuffers.current[i], MARGIN, wfOffset);
      wfOffset++;
    }

    if (device.ChannelCenterFrequency) {
      const chX = Math.round(MARGIN + (device.ChannelCenterFrequency - startFreq) * invDelta);
      const chW = Math.round(device.FilterBandwidth * invDelta);
      ctx.fillStyle = 'rgba(79, 195, 247, 0.1)';
      ctx.fillRect(chX - chW / 2, MARGIN_TOP, chW, FFT_HEIGHT);
      ctx.beginPath();
      ctx.lineWidth = 1.5;
      ctx.strokeStyle = '#ff5252';
      ctx.moveTo(chX, MARGIN_TOP);
      ctx.lineTo(chX, MARGIN_TOP + FFT_HEIGHT);
      ctx.stroke();
    }
  }, [fft, device]);

  useEffect(() => {
    if (device) {
      const canvas = canvasRef.current;
      if (canvas) {
        const w = device.DisplayPixels + MARGIN * 2;
        canvas.width = w;
        canvas.height = 600;
      }
    }
  }, [device?.DisplayPixels]);

  useEffect(() => {
    draw();
  }, [draw]);

  const handleClick = useCallback((e: MouseEvent) => {
    if (!device || !onTune) return;
    const canvas = canvasRef.current;
    if (!canvas) return;
    const rect = canvas.getBoundingClientRect();
    const scaleX = canvas.width / rect.width;
    const x = (e.clientX - rect.left) * scaleX;
    const pixelOffset = x - MARGIN;
    if (pixelOffset < 0 || pixelOffset > device.DisplayPixels) return;
    const startFreq = device.DisplayCenterFrequency - device.DisplayBandwidth / 2;
    const freqHz = Math.round(startFreq + pixelOffset * (device.DisplayBandwidth / device.DisplayPixels));
    onTune(freqHz);
  }, [device, onTune]);

  return (
    <canvas
      ref={canvasRef}
      onClick={handleClick}
      style={{ width: '100%', cursor: onTune ? 'crosshair' : 'default' }}
    />
  );
}
