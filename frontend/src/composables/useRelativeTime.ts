import { ref, onUnmounted, watch, type Ref } from 'vue'

/**
 * Formats a given Date object to a meaningful relative time string.
 *
 * @param date The date to compare with the current time
 * @returns Formatted relative time string
 */
export function formatRelativeTime(date: Date): string {
  const now = new Date()
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000)

  // Handle future dates or extremely small differences
  if (diffInSeconds < 30) {
    return 'Updated just now'
  }

  // Under 1 minute
  if (diffInSeconds < 60) {
    return 'Updated 30 seconds ago'
  }

  // Under 1 hour
  if (diffInSeconds < 3600) {
    const minutes = Math.floor(diffInSeconds / 60)
    return `Updated ${minutes} ${minutes === 1 ? 'minute' : 'minutes'} ago`
  }

  // Under 24 hours
  if (diffInSeconds < 86400) {
    const hours = Math.floor(diffInSeconds / 3600)
    return `Updated ${hours} ${hours === 1 ? 'hour' : 'hours'} ago`
  }

  // More than 24 hours
  const days = Math.floor(diffInSeconds / 86400)
  return `Updated ${days} ${days === 1 ? 'day' : 'days'} ago`
}

/**
 * A lightweight composable to manage a real-time relative update indicator.
 * It sets up a lightweight timer to recalculate the time display without API calls,
 * and immediately responds to reactive date updates.
 *
 * @param dateRef A reactive reference containing the last update date
 */
export function useRelativeTime(dateRef: Ref<Date | string | null | undefined>) {
  const relativeText = ref('Updated just now')

  const update = () => {
    if (!dateRef.value) {
      relativeText.value = 'Never updated'
      return
    }

    const date = typeof dateRef.value === 'string' ? new Date(dateRef.value) : dateRef.value
    relativeText.value = formatRelativeTime(date)
  }

  // Watch the reactive date source for external updates (e.g., API, WebSocket)
  // to trigger an immediate, responsive UI refresh.
  watch(
    () => dateRef.value,
    () => {
      update()
    },
    { immediate: true, deep: true }
  )

  // Lightweight 30-second recalculation timer (local to UI, does not trigger API requests)
  const timer = setInterval(() => {
    update()
  }, 30000)

  // Clean up the timer when the component using this composable is unmounted
  onUnmounted(() => {
    clearInterval(timer)
  })

  return {
    relativeText,
    refresh: update
  }
}
