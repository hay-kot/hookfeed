import { Requests } from '~~/lib/requests/requests'
import { route } from '~~/lib/api/base/urls'

/**
 * Composable that provides an authenticated API client
 * Automatically includes bearer token from auth state in all requests
 */
export const useApiClient = () => {
  const { token } = useAuth()

  // Create a requests instance with bearer token
  const client = new Requests()
    .withBearer(() => token.value)
    .withResponseInterceptor((response) => {
      // If we get a 401, the token is invalid - logout
      if (response.status === 401) {
        const { logout } = useAuth()
        logout()
        navigateTo('/auth/login')
      }
    })

  // Return wrapper object with convenient methods
  return {
    apiClient: {
      get: <T>(url: string) => client.get<T>({ url: route(url as any) }),
      post: <T, U = T>(url: string, body?: T) => client.post<T, U>({ url: route(url as any), body }),
      put: <T, U = T>(url: string, body?: T) => client.put<T, U>({ url: route(url as any), body }),
      patch: <T, U = T>(url: string, body?: T) => client.patch<T, U>({ url: route(url as any), body }),
      delete: <T>(url: string) => client.delete<T>({ url: route(url as any) }),
    }
  }
}
