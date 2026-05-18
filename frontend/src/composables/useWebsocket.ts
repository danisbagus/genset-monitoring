import { ref, onUnmounted } from 'vue'
import { API_CONFIG } from '@/config/api.config'

export type WsStatus = 'online' | 'reconnecting' | 'offline'

/**
 * Robust WebSocket composable with automatic reconnection and listener registry.
 */
export function useWebsocket(url?: string) {
  const socket = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const status = ref<WsStatus>('offline')
  const error = ref<Event | null>(null)
  
  // Construct WS URL: works with relative paths by using current host/protocol
  const getWsUrl = () => {
    const targetUrl = url || API_CONFIG.WS_URL
    if (targetUrl.startsWith('ws://') || targetUrl.startsWith('wss://')) {
      return targetUrl
    }
    
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    const path = targetUrl.startsWith('/') ? targetUrl : `/${targetUrl}`
    return `${protocol}//${host}${path}`
  }

  const wsUrl = getWsUrl()
  let reconnectTimer: number | null = null
  const reconnectInterval = 5000
  let isExplicitlyClosed = false

  const messageListeners = new Set<(data: any) => void>()

  function connect() {
    // If already connected or connecting, do not open a new connection
    if (socket.value && (socket.value.readyState === WebSocket.OPEN || socket.value.readyState === WebSocket.CONNECTING)) {
      return
    }

    isExplicitlyClosed = false
    status.value = reconnectTimer ? 'reconnecting' : 'offline'

    try {
      socket.value = new WebSocket(wsUrl)

      socket.value.onopen = () => {
        console.log('WS Connected to:', wsUrl)
        isConnected.value = true
        status.value = 'online'
        error.value = null
        if (reconnectTimer) {
          clearTimeout(reconnectTimer)
          reconnectTimer = null
        }
      }

      socket.value.onclose = () => {
        console.log('WS Disconnected from:', wsUrl)
        isConnected.value = false
        
        if (!isExplicitlyClosed) {
          status.value = 'reconnecting'
          scheduleReconnect()
        } else {
          status.value = 'offline'
        }
      }

      socket.value.onerror = (e) => {
        console.error('WS Error:', e)
        error.value = e
      }

      socket.value.onmessage = (event) => {
        try {
          const parsed = JSON.parse(event.data)
          messageListeners.forEach((listener) => listener(parsed))
        } catch (e) {
          messageListeners.forEach((listener) => listener(event.data))
        }
      }
    } catch (e) {
      console.error('WS Connection Error:', e)
      isConnected.value = false
      if (!isExplicitlyClosed) {
        status.value = 'reconnecting'
        scheduleReconnect()
      } else {
        status.value = 'offline'
      }
    }
  }

  function scheduleReconnect() {
    if (reconnectTimer) return
    reconnectTimer = window.setTimeout(() => {
      reconnectTimer = null
      connect()
    }, reconnectInterval)
  }

  function disconnect() {
    isExplicitlyClosed = true
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    if (socket.value) {
      socket.value.onopen = null
      socket.value.onclose = null
      socket.value.onerror = null
      socket.value.onmessage = null
      socket.value.close()
      socket.value = null
    }
    isConnected.value = false
    status.value = 'offline'
  }

  function sendMessage(data: any) {
    if (socket.value?.readyState === WebSocket.OPEN) {
      socket.value.send(JSON.stringify(data))
    } else {
      console.warn('WS not connected, cannot send message')
    }
  }

  function onMessage(callback: (data: any) => void) {
    messageListeners.add(callback)
    return () => {
      messageListeners.delete(callback)
    }
  }

  onUnmounted(() => {
    disconnect()
  })

  return {
    socket,
    isConnected,
    status,
    error,
    connect,
    disconnect,
    sendMessage,
    onMessage
  }
}
