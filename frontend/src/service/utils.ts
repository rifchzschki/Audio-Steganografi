import type { BaseResponse } from "@/models/response";
import apiClient, { handleApiError, handleApiResponse } from ".";
import type { AxiosError } from "axios";

export const checkApiHealth = async (): Promise<{
  status: string;
  timestamp: string;
}> => {
  try {
    const response = await apiClient.get<
      BaseResponse<{ status: string; timestamp: string }>
    >('/api/health');
    return handleApiResponse(response);
  } catch (error) {
    return handleApiError(error as AxiosError);
  }
};

export const createSteganographyFormData = (
  audioFile: File,
  config: {
    key?: string;
    lsbBits?: number;
    useEncryption?: boolean;
    useRandomStart?: boolean;
  },
  secretFile?: File,
  secretText?: string
): FormData => {
  const formData = new FormData();

  formData.append('audio', audioFile);
  formData.append('key', config.key || '');
  formData.append('lsb_bits', String(config.lsbBits || 2));
  formData.append('use_encryption', String(config.useEncryption || false));
  formData.append('use_random_start', String(config.useRandomStart || false));

  if (secretFile) {
    formData.append('secret_file', secretFile);
  }

  if (secretText) {
    formData.append('secret_text', secretText);
  }

  return formData;
};

export const validateFile = (
  file: File,
  allowedTypes: string[] = ['audio/mpeg', 'audio/mp3', 'audio/wav'],
  maxSizeMB: number = 50
): { valid: boolean; error?: string } => {
  if (!allowedTypes.includes(file.type)) {
    return {
      valid: false,
      error: `Unsupported file type. Allowed: ${allowedTypes.join(', ')}`,
    };
  }

  const maxSizeBytes = maxSizeMB * 1024 * 1024;
  if (file.size > maxSizeBytes) {
    return {
      valid: false,
      error: `File too large. Maximum size: ${maxSizeMB}MB`,
    };
  }

  return { valid: true };
};

export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

export const formatDuration = (seconds: number): string => {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = Math.floor(seconds % 60);
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
};