import type { BaseResponse, SteganographyResponse } from "@/models/response";
import apiClient, { handleApiError, handleApiResponse } from ".";
import type { AxiosError } from "axios";

export class SteganographyService {
  static async embedData(formData: FormData): Promise<SteganographyResponse> {
    try {
      const response = await apiClient.post<
        BaseResponse<SteganographyResponse>
      >('/api/steganography/embed', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      return handleApiResponse(response);
    } catch (error) {
      return handleApiError(error as AxiosError);
    }
  }

  static async extractData(formData: FormData): Promise<SteganographyResponse> {
    try {
      const response = await apiClient.post<
        BaseResponse<SteganographyResponse>
      >('/api/steganography/extract', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      return handleApiResponse(response);
    } catch (error) {
      return handleApiError(error as AxiosError);
    }
  }

  static async getCapacity(formData: FormData): Promise<{ capacity: number }> {
    try {
      const response = await apiClient.post<BaseResponse<{ capacity: number }>>(
        '/api/steganography/capacity',
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
}