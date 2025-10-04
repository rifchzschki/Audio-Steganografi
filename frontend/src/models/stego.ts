export type LsbBits = 1 | 2 | 3 | 4;

export interface EncodeRequestDto {
  audioFile: File;
  secretFile: File;
  key: string;
  lsbBits: LsbBits;
  useEncryption: boolean;
  useRandomStart: boolean;
}

export interface EncodeResponseDto {
  success: boolean;
  message: string;
  psnr?: number;
  stegoFileUrl?: string;
}

export interface DecodeRequestDto {
  stegoFile: File;
  key: string;
  useRandomStart: boolean;
  outputFileName?: string;
}

export interface DecodeResponseDto {
  success: boolean;
  message: string;
  secretFileUrl?: string;
  secretFilename?: string;
}

export interface EncodeFormState {
  coverFile: File | null;
  secretFile: File | null;
  key: string;
  useEncryption: boolean;
  useRandomStart: boolean;
  lsbBits: LsbBits;
}

export interface EncodeResultState {
  psnr?: number;
  stegoFileUrl?: string;
  stegoFilename?: string;
}

export interface DecodeFormState {
  stegoFile: File | null;
  key: string;
  useRandomStart: boolean;
  outputFileName?: string;
}

export interface DecodeResultState {
  secretFileUrl?: string;
  secretFilename?: string;
}
