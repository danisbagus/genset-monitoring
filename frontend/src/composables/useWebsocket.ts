import { ref, onUnmounted } from 'vue'
import { API_CONFIG } from '@/config/api.config'

export function useWebsocket(url?: string) {
  const socket = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const error = ref<Event | null>(null)
  
  // Construct WS URL: works with relative paths by using current host/protocol
  const getWsUrl = () => {
    if (url) return url
    if (API_CONFIG.WS_URL.startsWith('ws')) return API_CONFIG.WS_URL
    
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    return `${protocol}//${host}${API_CONFIG.WS_URL}`
  }

  const wsUrl = getWsUrl()
  let reconnectTimer: number | null = null
  const reconnectInterval = 5000

  function connect() {
    if (socket.value?.readyState === WebSocket.OPEN) return

    try {
      socket.value = new WebSocket(wsUrl)

      socket.value.onopen = () => {
        console.log('WS Connected')
        isConnected.value = true
        error.value = null
        if (reconnectTimer) {
          clearTimeout(reconnectTimer)
          reconnectTimer = null
        }
      }

      socket.value.onclose = () => {
        console.log('WS Disconnected')
        isConnected.value = false
        scheduleReconnect()
      }

      socket.value.onerror = (e) => {
        console.error('WS Error:', e)
        error.value = e
      }
    } catch (e) {
      console.error('WS Connection Error:', e)
      scheduleReconnect()
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
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    if (socket.value) {
      socket.value.close()
      socket.value = null
    }
  }

  function sendMessage(data: any) {
    if (socket.value?.readyState === WebSocket.OPEN) {
      socket.value.send(JSON.stringify(data))
    } else {
      console.warn('WS not connected, cannot send message')
    }
  }

  function onMessage(callback: (data: any) => void) {
    if (!socket.value) return
    
    // We wrap the existing onmessage or use addEventListener
    socket.value.addEventListener('message', (event) => {
      try {
        const data = JSON.parse(event.data)
        callback(data)
      } catch (e) {
        callback(event.data)
      }
    })
  }

  onUnmounted(() => {
    disconnect()
  })

  return {
    socket,
    isConnected,
    error,
    connect,
    disconnect,
    sendMessage,
    onMessage
  }
}
