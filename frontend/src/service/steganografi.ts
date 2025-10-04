import type {
  DecodeRequestDto,
  DecodeResponseDto,
  EncodeRequestDto,
  EncodeResponseDto,
} from '@/models/stego';
import type { AxiosError } from 'axios';
import apiClient, { handleApiError } from '.';

const MULTIPART_HEADERS = { 'Content-Type': 'multipart/form-data' } as const;

interface EncodeResponseRaw {
  success: boolean;
  message: string;
  psnr?: number;
  stego_file_url?: string;
}

interface DecodeResponseRaw {
  success: boolean;
  message: string;
  secret_file_url?: string;
  secret_filename?: string;
}

const toBoolString = (value: boolean) => (value ? 'true' : 'false');

const mapEncodeResponse = (data: EncodeResponseRaw): EncodeResponseDto => ({
  success: data.success,
  message: data.message,
  psnr: data.psnr,
  stegoFileUrl: data.stego_file_url,
});

const mapDecodeResponse = (data: DecodeResponseRaw): DecodeResponseDto => ({
  success: data.success,
  message: data.message,
  secretFileUrl: data.secret_file_url,
  secretFilename: data.secret_filename,
});

const buildBaseUrl = () => apiClient.defaults.baseURL ?? '';

export class SteganographyService {
  static async encode(payload: EncodeRequestDto): Promise<EncodeResponseDto> {
    try {
      const formData = new FormData();
      formData.append('audioFile', payload.audioFile);
      formData.append('secretFile', payload.secretFile);
      formData.append('key', payload.key);
      formData.append('lsbBits', String(payload.lsbBits));
      formData.append('useEncryption', toBoolString(payload.useEncryption));
      formData.append('useRandomStart', toBoolString(payload.useRandomStart));

      const { data } = await apiClient.post<EncodeResponseRaw>(
        '/api/encode',
        formData,
        { headers: MULTIPART_HEADERS }
      );

      return mapEncodeResponse(data);
    } catch (error) {
      return handleApiError(error as AxiosError);
    }
  }

  static async decode(payload: DecodeRequestDto): Promise<DecodeResponseDto> {
    try {
      const formData = new FormData();
      formData.append('stegoFile', payload.stegoFile);
      formData.append('key', payload.key);
      formData.append('useRandomStart', toBoolString(payload.useRandomStart));
      if (payload.outputFileName) {
        formData.append('outputFileName', payload.outputFileName);
      }

      const { data } = await apiClient.post<DecodeResponseRaw>(
        '/api/decode',
        formData,
        { headers: MULTIPART_HEADERS }
      );

      return mapDecodeResponse(data);
    } catch (error) {
      return handleApiError(error as AxiosError);
    }
  }

  static async downloadStego(filename: string): Promise<Blob> {
    try {
      const { data } = await apiClient.get<Blob>(
        `/api/download/stego/${encodeURIComponent(filename)}`,
        { responseType: 'blob' }
      );
      return data;
    } catch (error) {
      return handleApiError(error as AxiosError);
    }
  }

  static getStegoStreamUrl(filename: string): string {
    return `${buildBaseUrl()}/api/play/stego/${encodeURIComponent(filename)}`;
  }

  static async downloadExtracted(filename: string): Promise<Blob> {
    try {
      const { data } = await apiClient.get<Blob>(
        `/api/download/extracted/${encodeURIComponent(filename)}`,
        { responseType: 'blob' }
      );
      return data;
    } catch (error) {
      return handleApiError(error as AxiosError);
    }
  }

  static getStegoDownloadUrl(filename: string): string {
    return `${buildBaseUrl()}/api/download/stego/${encodeURIComponent(
      filename
    )}`;
  }

  static getExtractedDownloadUrl(filename: string): string {
    return `${buildBaseUrl()}/api/download/extracted/${encodeURIComponent(
      filename
    )}`;
  }
}
