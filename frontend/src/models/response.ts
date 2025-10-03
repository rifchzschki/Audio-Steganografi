export interface BaseResponse<T = unknown> {
  success: boolean;
  message: string;
  data?: T;
}

export interface SteganographyResponse {
  success: boolean;
  message: string;
  fileUrl?: string;
  fileName?: string;
  extractedData?: string;
  capacity?: number;
}