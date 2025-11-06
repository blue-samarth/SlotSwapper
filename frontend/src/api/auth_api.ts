import apiClient from './axios';
import type { LoginRequest, SignupRequest, AuthResponse } from '@/types/auth_types';

export const login = (credentials: LoginRequest) => {
  return apiClient.post<AuthResponse>('/auth/login', credentials);
};

export const signup = (data: SignupRequest) => {
  return apiClient.post<AuthResponse>('/auth/signup', data);
};

export const getCurrentUser = () => {
  return apiClient.get('/users/me');
};