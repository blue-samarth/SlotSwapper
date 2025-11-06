/**
 * Format a date string to a readable format
 * @param dateString - ISO date string
 * @param options - Intl.DateTimeFormatOptions
 * @returns Formatted date string
 */
export function formatDate(
  dateString: string,
  options?: Intl.DateTimeFormatOptions
): string {
  const date = new Date(dateString);
  const defaultOptions: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  };
  return date.toLocaleDateString('en-US', options || defaultOptions);
}

/**
 * Format a date string to include time
 * @param dateString - ISO date string
 * @returns Formatted date and time string
 */
export function formatDateTime(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

/**
 * Format time only from a date string
 * @param dateString - ISO date string
 * @returns Formatted time string
 */
export function formatTime(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
  });
}

/**
 * Calculate duration between two dates
 * @param startTime - ISO date string
 * @param endTime - ISO date string
 * @returns Formatted duration string (e.g., "1h 30m")
 */
export function calculateDuration(startTime: string, endTime: string): string {
  const start = new Date(startTime);
  const end = new Date(endTime);
  const durationMs = end.getTime() - start.getTime();
  const minutes = Math.round(durationMs / 60000);

  if (minutes < 60) {
    return `${minutes}m`;
  }

  const hours = Math.floor(minutes / 60);
  const remainingMinutes = minutes % 60;
  return remainingMinutes > 0 ? `${hours}h ${remainingMinutes}m` : `${hours}h`;
}

/**
 * Get relative time (e.g., "2 hours ago", "in 3 days")
 * @param dateString - ISO date string
 * @returns Relative time string
 */
export function getRelativeTime(dateString: string): string {
  const date = new Date(dateString);
  const now = new Date();
  const diffMs = date.getTime() - now.getTime();
  const diffMinutes = Math.round(diffMs / 60000);
  const diffHours = Math.round(diffMs / 3600000);
  const diffDays = Math.round(diffMs / 86400000);

  if (Math.abs(diffMinutes) < 60) {
    if (diffMinutes === 0) return 'just now';
    if (diffMinutes > 0) return `in ${diffMinutes} ${diffMinutes === 1 ? 'minute' : 'minutes'}`;
    return `${Math.abs(diffMinutes)} ${Math.abs(diffMinutes) === 1 ? 'minute' : 'minutes'} ago`;
  }

  if (Math.abs(diffHours) < 24) {
    if (diffHours > 0) return `in ${diffHours} ${diffHours === 1 ? 'hour' : 'hours'}`;
    return `${Math.abs(diffHours)} ${Math.abs(diffHours) === 1 ? 'hour' : 'hours'} ago`;
  }

  if (Math.abs(diffDays) < 7) {
    if (diffDays > 0) return `in ${diffDays} ${diffDays === 1 ? 'day' : 'days'}`;
    return `${Math.abs(diffDays)} ${Math.abs(diffDays) === 1 ? 'day' : 'days'} ago`;
  }

  return formatDate(dateString);
}

/**
 * Check if a date is in the past
 * @param dateString - ISO date string
 * @returns True if date is in the past
 */
export function isPast(dateString: string): boolean {
  return new Date(dateString) < new Date();
}

/**
 * Check if a date is today
 * @param dateString - ISO date string
 * @returns True if date is today
 */
export function isToday(dateString: string): boolean {
  const date = new Date(dateString);
  const today = new Date();
  return (
    date.getDate() === today.getDate() &&
    date.getMonth() === today.getMonth() &&
    date.getFullYear() === today.getFullYear()
  );
}

/**
 * Get the day of week
 * @param dateString - ISO date string
 * @returns Day of week (e.g., "Monday")
 */
export function getDayOfWeek(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', { weekday: 'long' });
}
