import { useColorMode } from '@vueuse/core'

export function useTheme() {
  const mode = useColorMode({
    attribute: 'class',
    modes: {
      dark: 'dark',
      light: '',
    },
  })

  return {
    mode,
    isDark: mode.value === 'dark',
    toggleTheme: () => {
      mode.value = mode.value === 'dark' ? 'light' : 'dark'
    }
  }
}
