export const EventStatus = {
  BUSY: 'BUSY',
  SWAPPABLE: 'SWAPPABLE',
  SWAP_PENDING: 'SWAP_PENDING',
} as const;
export type EventStatus = typeof EventStatus[keyof typeof EventStatus];

export interface Event {
  id: number;
  title: string;
  user_id: number;
  start_time: string;
  end_time: string;
  status: EventStatus;
  created_at: string;
  updated_at: string;
}

export interface CreateEventRequest {
  title: string;
  start_time: string;
  end_time: string;
}

export interface UpdateEventStatusRequest {
  status: typeof EventStatus.BUSY | typeof EventStatus.SWAPPABLE;
}