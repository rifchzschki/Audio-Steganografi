export interface BaseResponse<T = unknown> {
  success: boolean;
  message: string;
  data?: T;
}

export interface EncodeResponse {
  success: boolean;
  message: string;
  psnr?: number;
  stegoFileUrl?: string;
}

export interface DecodeResponse {
  success: boolean;
  message: string;
  secretFileUrl?: string;
  secretFilename?: string;
}
