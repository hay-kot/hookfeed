/**
 * Auth middleware to protect routes
 * Redirects unauthenticated users to login page
 */
export default defineNuxtRouteMiddleware(async (to) => {
  const { isAuthenticated, initAuth } = useAuth()

  // Initialize auth state if not already done
  await initAuth()

  // Allow access to auth pages without being logged in
  if (to.path.startsWith('/auth/')) {
    // If already authenticated, redirect to home
    if (isAuthenticated.value) {
      return navigateTo('/')
    }
    return
  }

  // Require authentication for all other pages
  if (!isAuthenticated.value) {
    return navigateTo('/auth/login')
  }
})
