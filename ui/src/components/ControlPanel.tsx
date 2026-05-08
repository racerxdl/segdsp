import { useCallback } from 'preact/hooks';
import type { DeviceInfo, ControlParams } from '../services/ws';
import { ws } from '../services/ws';

interface Props {
  device: DeviceInfo | null;
}

function toMHz(hz: number): string {
  return (hz / 1e6).toFixed(3);
}

function fromMHz(mhz: string): number {
  return Math.round(parseFloat(mhz) * 1e6);
}

function tokHz(hz: number): string {
  return (hz / 1e3).toFixed(1);
}

function fromKHz(khz: string): number {
  return Math.round(parseFloat(khz) * 1e3);
}

export function ControlPanel({ device }: Props) {
  const send = useCallback((params: ControlParams) => {
    ws.sendControl(params);
  }, []);

  if (!device) return <div class="panel"><p>Connecting...</p></div>;

  const canControl = device.WebCanControl;

  return (
    <div class="panel">
      <div class="panel-section">
        <h3>Device</h3>
        <div class="info-row">
          <span class="label">Station</span>
          <span class="value">{device.StationName}</span>
        </div>
        <div class="info-row">
          <span class="label">Device</span>
          <span class="value">{device.DeviceName}</span>
        </div>
        <div class="info-row">
          <span class="label">Mode</span>
          <span class="value">
            <select
              disabled={!canControl}
              value={device.DemodulatorMode}
              onChange={(e) => send({ DemodulatorMode: (e.target as HTMLSelectElement).value })}
            >
              <option value="FM">FM</option>
              <option value="AM">AM</option>
            </select>
          </span>
        </div>
        <div class="info-row">
          <span class="label">Level</span>
          <span class="value">{device.IsMuted ? 'Muted' : 'Active'}</span>
        </div>
      </div>

      <div class="panel-section">
        <h3>Frequency</h3>
        <div class="control-row">
          <label>Channel</label>
          <div class="freq-input">
            <input
              type="number"
              step="0.001"
              value={toMHz(device.ChannelCenterFrequency)}
              disabled={!canControl}
              onChange={(e) => send({ ChannelFrequency: fromMHz((e.target as HTMLInputElement).value) })}
            />
            <span class="unit">MHz</span>
          </div>
        </div>
        <div class="control-row">
          <label>FFT</label>
          <div class="freq-input">
            <input
              type="number"
              step="0.001"
              value={toMHz(device.DisplayCenterFrequency)}
              disabled={!canControl}
              onChange={(e) => send({ FFTfrequency: fromMHz((e.target as HTMLInputElement).value) })}
            />
            <span class="unit">MHz</span>
          </div>
        </div>
      </div>

      <div class="panel-section">
        <h3>Filter</h3>
        <div class="control-row">
          <label>Bandwidth</label>
          <div class="freq-input">
            <input
              type="number"
              step="0.1"
              value={tokHz(device.FilterBandwidth)}
              disabled={!canControl}
              onChange={(e) => send({ FilterBandwidth: fromKHz((e.target as HTMLInputElement).value) })}
            />
            <span class="unit">kHz</span>
          </div>
        </div>
        <div class="control-row">
          <label>Channel BW</label>
          <span class="value static">{tokHz(device.CurrentSampleRate)} kHz</span>
        </div>
      </div>

      <div class="panel-section">
        <h3>Display</h3>
        <div class="control-row">
          <label>Offset</label>
          <input
            type="range"
            min={-120}
            max={0}
            value={device.DisplayOffset}
            disabled={!canControl}
            onInput={(e) => send({ DisplayOffset: parseInt((e.target as HTMLInputElement).value) })}
          />
          <span class="range-value">{device.DisplayOffset} dB</span>
        </div>
        <div class="control-row">
          <label>Range</label>
          <input
            type="range"
            min={20}
            max={120}
            value={device.DisplayRange}
            disabled={!canControl}
            onInput={(e) => send({ DisplayRange: parseInt((e.target as HTMLInputElement).value) })}
          />
          <span class="range-value">{device.DisplayRange} dB</span>
        </div>
      </div>

      {!canControl && (
        <div class="locked-notice">
          Controls disabled (WebCanControl=false)
        </div>
      )}
    </div>
  );
}
