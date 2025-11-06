import type { Event } from './event_types';
import type { User } from './auth_types';

export const SwapStatus = {
  PENDING: 'PENDING',
  ACCEPTED: 'ACCEPTED',
  REJECTED: 'REJECTED',
  CANCELLED: 'CANCELLED',
} as const;
export type SwapStatus = typeof SwapStatus[keyof typeof SwapStatus];


export interface SwapRequest {
  id: number;
  requester_event_id: number;
  receiver_event_id: number;
  requester_id: number;
  receiver_id: number;
  status: SwapStatus;
  requester_event: Event;
  receiver_event: Event;
  created_at: string;
  updated_at: string;
}

export interface InitiateSwapRequest {
  requester_event_id: number;
  receiver_event_id: number;
}

export interface SwappableSlot extends Event {
  user: User;
}