import type { AudioMetadata } from "@/models/audio";
import apiClient, { handleApiError, handleApiResponse } from ".";
import type { BaseResponse } from "@/models/response";
import type { AxiosError } from "axios";

export class AudioService {
  static async getMetadata(file: File): Promise<AudioMetadata> {
    try {
      const formData = new FormData();
      formData.append('audio', file);

      const response = await apiClient.post<BaseResponse<AudioMetadata>>(
        '/api/audio/metadata',
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        }
      );
      return handleApiResponse(response);
    } catch (error) {
      return handleApiError(error as AxiosError);
    }
  }

  static async validateFormat(
    file: File
  ): Promise<{ valid: boolean; format: string }> {
    try {
      const formData = new FormData();
      formData.append('audio', file);

      const response = await apiClient.post<
        BaseResponse<{ valid: boolean; format: string }>
      >('/api/audio/validate', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      return handleApiResponse(response);
    } catch (error) {
      return handleApiError(error as AxiosError);
    }
  }

  static async convertFormat(file: File, targetFormat: string): Promise<Blob> {
    try {
      const formData = new FormData();
      formData.append('audio', file);
      formData.append('format', targetFormat);

      const response = await apiClient.post('/api/audio/convert', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        responseType: 'blob',
      });
      return response.data;
    } catch (error) {
      return handleApiError(error as AxiosError);
    }
  }
}