export type AudioCompareProps = {
  originalUrl?: string;
  originalName?: string;
  stegoUrl?: string;
  stegoName?: string;
};

export type PlayerState = {
  src: string;
  time: number;
  volume: number;
  isPlaying: boolean;
};

export const initialPlayerState: PlayerState = {
  src: '',
  time: 0,
  volume: 0,
  isPlaying: false,
};

export interface AudioMetadata {
  sampleRate: number;
  channels: number;
  bitDepth: number;
  duration: number;
  totalBytes: number;
}
