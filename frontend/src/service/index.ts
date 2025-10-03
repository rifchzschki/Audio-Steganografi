import type { BaseResponse } from '@/models/response';
import axios, { AxiosError } from 'axios';
import type { AxiosInstance, AxiosResponse } from 'axios';

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
const API_TIMEOUT = 30000; 

const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: API_TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
});

apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    console.log(`API Response: ${response.status} ${response.config.url}`);
    return response;
  },
  (error: AxiosError) => {
    console.error('Response Error:', error);

    return Promise.reject(error);
  }
);

export const handleApiResponse = <T>(response: AxiosResponse<BaseResponse<T>>): T => {
  if (response.data.success && response.data.data) {
    return response.data.data;
  }
  throw new Error(response.data.message || 'Unknown API error');
};

export const handleApiError = (error: AxiosError): never => {
  if (error.response?.data) {
    const errorData = error.response.data as BaseResponse<null>;
    throw new Error(errorData.message || 'API request failed');
  }
  throw new Error(error.message || 'Network error');
};


export default apiClient;
