export const STORAGE_KEYS = {
  AUTH_DATA: 'genset_auth_data',
  USER_DATA: 'genset_user_data'
} as const

export const getStorage = <T>(key: string): T | null => {
  const item = localStorage.getItem(key)
  if (!item) return null
  try {
    return JSON.parse(item) as T
  } catch (error) {
    console.error(`Error parsing storage key "${key}":`, error)
    return null
  }
}

export const setStorage = <T>(key: string, value: T): void => {
  try {
    localStorage.setItem(key, JSON.stringify(value))
  } catch (error) {
    console.error(`Error setting storage key "${key}":`, error)
  }
}

export const removeStorage = (key: string): void => {
  localStorage.removeItem(key)
}

export const clearStorage = (): void => {
  localStorage.clear()
}
