export interface ApiError {
  error: string;
  message: string;
}

export interface ApiResponse<T> {
  data: T;
  message?: string;
}