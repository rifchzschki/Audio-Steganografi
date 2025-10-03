import type { BaseResponse } from "@/models/response";
import apiClient, { handleApiError, handleApiResponse } from ".";
import type { AxiosError } from "axios";

export class FileService {
  static async uploadFile(
    file: File,
    category: string = 'general'
  ): Promise<{ url: string; id: string }> {
    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('category', category);

      const response = await apiClient.post<
        BaseResponse<{ url: string; id: string }>
      >('/api/file/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      return handleApiResponse(response);
    } catch (error) {
      return handleApiError(error as AxiosError);
    }
  }

  static async downloadFile(fileId: string, filename?: string): Promise<Blob> {
    try {
      const response = await apiClient.get(`/api/file/download/${fileId}`, {
        responseType: 'blob',
      });

      const blob = new Blob([response.data]);
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename || `file_${fileId}`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);

      return blob;
    } catch (error) {
      return handleApiError(error as AxiosError);
    }
  }

  static async deleteFile(fileId: string): Promise<{ success: boolean }> {
    try {
      const response = await apiClient.delete<
        BaseResponse<{ success: boolean }>
      >(`/api/file/${fileId}`);
      return handleApiResponse(response);
    } catch (error) {
      return handleApiError(error as AxiosError);
    }
  }
}