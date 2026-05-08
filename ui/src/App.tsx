import { useState, useEffect, useCallback } from 'preact/hooks';
import { ws } from './services/ws';
import type { DeviceInfo, FFTMessage } from './services/ws';
import { Spectrum } from './components/Spectrum';
import { ControlPanel } from './components/ControlPanel';

export function App() {
  const [device, setDevice] = useState<DeviceInfo | null>(null);
  const [fft, setFFT] = useState<FFTMessage | null>(null);
  const [connected, setConnected] = useState(false);
  const [showPanel, setShowPanel] = useState(true);

  useEffect(() => {
    ws.onDeviceMessage((msg) => {
      setDevice(msg);
      setConnected(true);
      document.title = msg.StationName;
    });
    ws.onFFTMessage((msg) => {
      setFFT(msg);
    });
    ws.onDisconnected(() => {
      setConnected(false);
    });
    ws.connect();
    return () => ws.disconnect();
  }, []);

  const handleTune = useCallback((freqHz: number) => {
    if (device?.WebCanControl) {
      ws.sendControl({ ChannelFrequency: freqHz });
    }
  }, [device?.WebCanControl]);

  const level = fft?.DemodOutputLevel ?? 0;

  return (
    <div class="app">
      <header class="header">
        <div class="header-left">
          <h1>SegDSP</h1>
          <span class={`status ${connected ? 'connected' : 'disconnected'}`}>
            {connected ? 'Connected' : 'Connecting...'}
          </span>
          {device && <span class="station">{device.StationName} ({device.DeviceName})</span>}
        </div>
        <div class="header-right">
          <span class="level">{level.toFixed(1)} dB</span>
          <button class="toggle-panel" onClick={() => setShowPanel(!showPanel)}>
            {showPanel ? 'Hide Panel' : 'Show Panel'}
          </button>
        </div>
      </header>
      <div class="main">
        <div class="spectrum-area">
          <Spectrum fft={fft} device={device} onTune={handleTune} />
        </div>
        {showPanel && (
          <div class="side-panel">
            <ControlPanel device={device} />
          </div>
        )}
      </div>
    </div>
  );
}
