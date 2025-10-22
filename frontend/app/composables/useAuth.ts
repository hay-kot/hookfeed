import type { DtosUser, DtosUserAuthenticate, DtosUserRegister, DtosUserSession } from '~~/lib/api/types/data-contracts'
import { Requests } from '~~/lib/requests/requests'
import { route } from '~~/lib/api/base/urls'

const AUTH_TOKEN_KEY = 'auth-token'

/**
 * Composable for managing user authentication
 * Persists auth token to localStorage
 * Provides login, logout, register, and user profile methods
 */
export const useAuth = () => {
  const user = useState<DtosUser | null>('auth-user', () => null)
  const token = useState<string | null>('auth-token', () => null)
  const isInitialized = useState<boolean>('auth-initialized', () => false)

  const requests = new Requests()

  // Initialize auth state from localStorage on mount
  const initAuth = async () => {
    if (import.meta.client && !isInitialized.value) {
      const savedToken = localStorage.getItem(AUTH_TOKEN_KEY)
      if (savedToken) {
        token.value = savedToken
        // Attempt to fetch user profile to validate token
        await fetchUser()
      }
      isInitialized.value = true
    }
  }

  // Login with email and password
  const login = async (email: string, password: string): Promise<{ success: boolean; error?: string }> => {
    const credentials: DtosUserAuthenticate = { email, password }

    const response = await requests.post<DtosUserAuthenticate, DtosUserSession>({
      url: route('/users/login/'),
      body: credentials,
    })

    if (response.error || !response.data.token) {
      return {
        success: false,
        error: 'Invalid email or password',
      }
    }

    token.value = response.data.token
    if (import.meta.client) {
      localStorage.setItem(AUTH_TOKEN_KEY, response.data.token)
    }

    // Fetch user profile after successful login
    await fetchUser()

    return { success: true }
  }

  // Register new user
  const register = async (
    username: string,
    email: string,
    password: string,
  ): Promise<{ success: boolean; error?: string }> => {
    const userData: DtosUserRegister = { username, email, password }

    const response = await requests.post<DtosUserRegister, DtosUserSession>({
      url: route('/users/register/'),
      body: userData,
    })

    if (response.error || !response.data.token) {
      return {
        success: false,
        error: 'Registration failed',
      }
    }

    token.value = response.data.token
    if (import.meta.client) {
      localStorage.setItem(AUTH_TOKEN_KEY, response.data.token)
    }

    // Fetch user profile after successful registration
    await fetchUser()

    return { success: true }
  }

  // Fetch current user profile
  const fetchUser = async (): Promise<boolean> => {
    if (!token.value) {
      return false
    }

    const response = await requests
      .withBearer(() => token.value)
      .get<DtosUser>({
        url: route('/users/self/'),
      })

    if (response.error) {
      // Token is invalid, clear auth state
      logout()
      return false
    }

    user.value = response.data
    return true
  }

  // Logout user
  const logout = () => {
    user.value = null
    token.value = null
    if (import.meta.client) {
      localStorage.removeItem(AUTH_TOKEN_KEY)
    }
  }

  const isAuthenticated = computed(() => !!token.value && !!user.value)

  return {
    user: readonly(user),
    token: readonly(token),
    isAuthenticated,
    initAuth,
    login,
    register,
    logout,
    fetchUser,
  }
}
