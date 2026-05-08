export interface DeviceInfo {
  MessageType: string;
  DeviceName: string;
  DisplayBandwidth: number;
  DisplayCenterFrequency: number;
  DisplayOffset: number;
  DisplayRange: number;
  DisplayPixels: number;
  ChannelCenterFrequency: number;
  CurrentSampleRate: number;
  Gain: number;
  OutputRate: number;
  FilterBandwidth: number;
  DemodulatorMode: string;
  DemodulatorParams: Record<string, unknown>;
  StationName: string;
  WebCanControl: boolean;
  TCPCanControl: boolean;
  IsMuted: boolean;
}

export interface FFTMessage {
  MessageType: string;
  DemodOutputLevel: number;
  FFTData: string;
}

export type MessageHandler = (data: FFTMessage | DeviceInfo) => void;

export interface ControlParams {
  ChannelFrequency?: number;
  FFTfrequency?: number;
  DemodulatorMode?: string;
  FilterBandwidth?: number;
  Squelch?: number;
  DisplayOffset?: number;
  DisplayRange?: number;
}

export class WSService {
  private socket: WebSocket | null = null;
  private onFFT: ((msg: FFTMessage) => void) | null = null;
  private onDevice: ((msg: DeviceInfo) => void) | null = null;
  private onDisconnect: (() => void) | null = null;
  private audioCtx: AudioContext | null = null;
  private audioNode: ScriptProcessorNode | null = null;
  private audioBuffers: Float32Array[] = [];
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private avgTraffic = 0;
  private trafficSum = 0;

  onFFTMessage(handler: (msg: FFTMessage) => void) {
    this.onFFT = handler;
  }

  onDeviceMessage(handler: (msg: DeviceInfo) => void) {
    this.onDevice = handler;
  }

  onDisconnected(handler: () => void) {
    this.onDisconnect = handler;
  }

  getAvgTraffic() {
    return this.avgTraffic;
  }

  private initAudio(sampleRate: number) {
    if (this.audioCtx) return;
    this.audioCtx = new AudioContext({ sampleRate });
    this.audioNode = this.audioCtx.createScriptProcessor(16384, 0, 1);
    const source = this.audioCtx.createBufferSource();
    source.connect(this.audioNode);
    this.audioNode.connect(this.audioCtx.destination);

    this.audioNode.onaudioprocess = (e) => {
      if (this.audioBuffers.length > 0) {
        e.outputBuffer.getChannelData(0).set(this.audioBuffers.shift()!);
      }
    };

    source.start();
  }

  connect() {
    const proto = location.protocol === 'https:' ? 'wss://' : 'ws://';
    const url = proto + location.host + '/ws';
    this.socket = new WebSocket(url);

    this.startTrafficMonitor();

    this.socket.onopen = () => {
      console.log('WS connected');
    };

    this.socket.onclose = () => {
      console.log('WS disconnected');
      this.socket = null;
      this.onDisconnect?.();
      this.reconnectTimer = setTimeout(() => this.connect(), 1500);
    };

    this.socket.onmessage = (evt) => {
      if (typeof evt.data !== 'string') {
        this.handleBinaryAudio(evt.data);
        return;
      }
      this.trafficSum += evt.data.length;
      try {
        const data = JSON.parse(evt.data);
        switch (data.MessageType) {
          case 'fft':
            this.onFFT?.(data);
            break;
          case 'device':
            this.onDevice?.(data);
            this.initAudio(data.OutputRate);
            break;
          case 'controlAck':
            break;
          default:
            console.log('Unknown message type:', data.MessageType);
        }
      } catch (e) {
        console.error('JSON parse error:', e);
      }
    };

    this.socket.onerror = () => {
      console.error('WS error');
    };
  }

  private handleBinaryAudio(data: Blob | ArrayBuffer) {
    if (!this.audioCtx) return;
    if (data instanceof Blob) {
      const reader = new FileReader();
      reader.onload = (e) => {
        if (e.target?.result) {
          this.trafficSum += (e.target.result as ArrayBuffer).byteLength;
          this.audioBuffers.push(new Float32Array(e.target.result as ArrayBuffer));
        }
      };
      reader.readAsArrayBuffer(data);
    } else {
      this.trafficSum += data.byteLength;
      this.audioBuffers.push(new Float32Array(data));
    }
  }

  sendControl(params: ControlParams) {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) return;
    const msg = { MessageType: 'control', ...params };
    this.socket.send(JSON.stringify(msg));
  }

  disconnect() {
    if (this.reconnectTimer) clearTimeout(this.reconnectTimer);
    this.socket?.close();
    this.audioCtx?.close();
    this.audioCtx = null;
  }

  private startTrafficMonitor() {
    setInterval(() => {
      this.avgTraffic = this.trafficSum;
      this.trafficSum = 0;
    }, 1000);
  }
}

export const ws = new WSService();
