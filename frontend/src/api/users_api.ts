import apiClient from './axios';
import type { User } from '@/types/user_types';

export const updateUsername = (username: string) => {
  return apiClient.patch<User>('/users/me', { username });
};

export const getProfile = () => {
  return apiClient.get<User>('/users/me');
};