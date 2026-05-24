import { ref, computed } from 'vue'
import { monitoringApi } from '@/services/api/monitoring.api'
import type {
  MonitoringDeviceFilters,
  MonitoringDeviceItem,
  MonitoringPagination
} from '@/types/monitoring'

const DEFAULT_FILTERS: MonitoringDeviceFilters = {
  search: '',
  online: '',
  engine_running: '',
  status: '',
  sort_by: 'last_seen_at',
  sort_order: 'desc',
  limit: 20,
  offset: 0
}

/**
 * useMonitoringDevices
 *
 * Composable that owns all state for the monitoring device list page.
 * Ready for realtime WebSocket partial updates (upsertDevice).
 */
export function useMonitoringDevices() {
  const devices = ref<MonitoringDeviceItem[]>([])
  const pagination = ref<MonitoringPagination>({ limit: 20, page: 1, total: 0 })
  const filters = ref<MonitoringDeviceFilters>({ ...DEFAULT_FILTERS })

  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Derived current page (offset-based)
  const currentPage = computed(() =>
    Math.floor(filters.value.offset / filters.value.limit) + 1
  )

  const totalPages = computed(() =>
    Math.ceil(pagination.value.total / filters.value.limit)
  )

  // ── Fetch ────────────────────────────────────────────────────────────────────

  const fetchDevices = async () => {
    isLoading.value = true
    error.value = null

    try {
      const result = await monitoringApi.getDevices(filters.value)
      devices.value = result.devices ?? []
      pagination.value = result.pagination ?? { limit: filters.value.limit, page: 1, total: 0 }
    } catch (err: unknown) {
      console.error('[monitoring] fetchDevices error:', err)
      error.value = 'Failed to load device list. Please check your connection and try again.'
    } finally {
      isLoading.value = false
    }
  }

  // ── Filter helpers ────────────────────────────────────────────────────────────

  const applyFilters = () => {
    filters.value.offset = 0 // Reset to first page on filter change
    fetchDevices()
  }

  const resetFilters = () => {
    filters.value = { ...DEFAULT_FILTERS }
    fetchDevices()
  }

  const goToPage = (page: number) => {
    filters.value.offset = (page - 1) * filters.value.limit
    fetchDevices()
  }

  // ── Realtime upsert (WebSocket ready) ─────────────────────────────────────────

  const upsertDevice = (updated: Partial<MonitoringDeviceItem> & { device_id: string }) => {
    const idx = devices.value.findIndex(d => d.device_id === updated.device_id)
    if (idx !== -1) {
      devices.value.splice(idx, 1, { ...devices.value[idx], ...updated })
    }
  }

  return {
    // State
    devices,
    pagination,
    filters,
    isLoading,
    error,
    currentPage,
    totalPages,

    // Actions
    fetchDevices,
    applyFilters,
    resetFilters,
    goToPage,
    upsertDevice
  }
}
