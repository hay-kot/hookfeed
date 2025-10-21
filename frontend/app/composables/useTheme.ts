/**
 * Composable for managing dark/light theme using DaisyUI
 * Persists theme preference to localStorage
 * Uses 'silk' theme for light mode and 'dark' theme for dark mode
 */
export const useTheme = () => {
  const theme = useState<'light' | 'dark'>('theme', () => 'light')

  // Map internal theme names to DaisyUI theme names
  const themeMap = {
    light: 'lofi',
    dark: 'dark',
  }

  // Initialize theme from localStorage on mount
  const initTheme = () => {
    if (import.meta.client) {
      const savedTheme = localStorage.getItem('theme') as 'light' | 'dark' | null
      if (savedTheme) {
        theme.value = savedTheme
      } else {
        // Check system preference
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
        theme.value = prefersDark ? 'dark' : 'light'
      }
      applyTheme(theme.value)
    }
  }

  // Apply theme to HTML element
  const applyTheme = (newTheme: 'light' | 'dark') => {
    if (import.meta.client) {
      const daisyTheme = themeMap[newTheme]
      document.documentElement.setAttribute('data-theme', daisyTheme)
    }
  }

  // Toggle between light and dark
  const toggleTheme = () => {
    const newTheme = theme.value === 'light' ? 'dark' : 'light'
    theme.value = newTheme
    applyTheme(newTheme)

    if (import.meta.client) {
      localStorage.setItem('theme', newTheme)
    }
  }

  // Set specific theme
  const setTheme = (newTheme: 'light' | 'dark') => {
    theme.value = newTheme
    applyTheme(newTheme)

    if (import.meta.client) {
      localStorage.setItem('theme', newTheme)
    }
  }

  return {
    theme: readonly(theme),
    initTheme,
    toggleTheme,
    setTheme,
  }
}
