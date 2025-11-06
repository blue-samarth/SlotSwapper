import apiClient from './axios';
import type { Event, CreateEventRequest, UpdateEventStatusRequest } from '@/types/event_types';

export const getEvents = () => {
  return apiClient.get<Event[]>('/events');
};

export const getEventById = (id: number) => {
  return apiClient.get<Event>(`/events/${id}`);
};

export const createEvent = (data: CreateEventRequest) => {
  return apiClient.post<Event>('/events', data);
};

export const updateEventStatus = (id: number, data: UpdateEventStatusRequest) => {
  return apiClient.patch<Event>(`/events/${id}/status`, data);
};

export const deleteEvent = (id: number) => {
  return apiClient.delete(`/events/${id}`);
};