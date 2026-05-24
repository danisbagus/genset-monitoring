import { ref } from 'vue'
import { monitoringApi } from '@/services/api/monitoring.api'
import type { MonitoringDeviceDetailOutput } from '@/types/monitoring'

/**
 * useMonitoringDeviceDetail
 *
 * Composable that owns all state for the device detail page.
 * Architecture is ready for realtime partial updates via WebSocket.
 */
export function useMonitoringDeviceDetail() {
  const detail = ref<MonitoringDeviceDetailOutput | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const notFound = ref(false)

  // ── Fetch ─────────────────────────────────────────────────────────────────────

  const fetchDetail = async (deviceId: string) => {
    isLoading.value = true
    error.value = null
    notFound.value = false
    detail.value = null

    try {
      detail.value = await monitoringApi.getDeviceDetail(deviceId)
    } catch (err: any) {
      console.error('[monitoring] fetchDetail error:', err)
      if (err?.response?.status === 404) {
        notFound.value = true
      } else {
        error.value = 'Failed to load device details. Please try again.'
      }
    } finally {
      isLoading.value = false
    }
  }

  // ── Realtime patch (WebSocket ready) ──────────────────────────────────────────

  const patchLatestState = (update: Partial<MonitoringDeviceDetailOutput['latest_state']>) => {
    if (!detail.value) return
    detail.value = {
      ...detail.value,
      latest_state: { ...detail.value.latest_state, ...update }
    }
  }

  const patchConnectivity = (update: Partial<MonitoringDeviceDetailOutput['connectivity']>) => {
    if (!detail.value) return
    detail.value = {
      ...detail.value,
      connectivity: { ...detail.value.connectivity, ...update }
    }
  }

  const patchEngineTelemetry = (update: Partial<NonNullable<MonitoringDeviceDetailOutput['engine_telemetry']>>) => {
    if (!detail.value) return
    detail.value = {
      ...detail.value,
      engine_telemetry: { ...(detail.value.engine_telemetry ?? {}), ...update } as any
    }
  }

  const patchElectricalTelemetry = (update: Partial<NonNullable<MonitoringDeviceDetailOutput['electrical_telemetry']>>) => {
    if (!detail.value) return
    detail.value = {
      ...detail.value,
      electrical_telemetry: { ...(detail.value.electrical_telemetry ?? {}), ...update } as any
    }
  }

  return {
    // State
    detail,
    isLoading,
    error,
    notFound,

    // Actions
    fetchDetail,

    // Realtime patch helpers
    patchLatestState,
    patchConnectivity,
    patchEngineTelemetry,
    patchElectricalTelemetry
  }
}
