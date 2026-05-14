/**
 * Checks if the token is expired or will expire within the buffer time.
 * @param expiresAt The timestamp when the token expires (in milliseconds)
 * @returns boolean
 */
export const isTokenExpired = (expiresAt: number | null): boolean => {
  if (!expiresAt) return true
  
  // Refresh 1 minute before actual expiration to be safe
  const buffer = 60 * 1000 
  const now = Date.now()
  return now + buffer >= expiresAt
}

/**
 * Calculates the expiration timestamp in milliseconds.
 * @param secondsIn The number of seconds until expiration
 * @returns number (timestamp in ms)
 */
export const calculateExpiry = (secondsIn: number): number => {
  return Date.now() + secondsIn * 1000
}
