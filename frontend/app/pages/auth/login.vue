<script setup lang="ts">
import IconLogin from "~icons/mdi/login";
import IconEmail from "~icons/mdi/email";
import IconLock from "~icons/mdi/lock";
import IconWeatherSunny from "~icons/mdi/weather-sunny";
import IconWeatherNight from "~icons/mdi/weather-night";

// Disable default layout for auth pages
definePageMeta({
  layout: false,
});

const router = useRouter();
const { login } = useAuth();
const { theme, initTheme, toggleTheme } = useTheme();

// Initialize theme on mount
onMounted(() => {
  initTheme();
});

// Form state
const email = ref("");
const password = ref("");
const error = ref("");
const isLoading = ref(false);

// Handle login
const handleLogin = async () => {
  error.value = "";
  isLoading.value = true;

  try {
    const result = await login(email.value, password.value);

    if (result.success) {
      // Redirect to home page
      await router.push("/");
    } else {
      error.value = result.error || "Login failed";
    }
  } catch (e) {
    error.value = "An unexpected error occurred";
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div
    class="min-h-screen flex items-center justify-center bg-base-200 p-4 relative"
  >
    <!-- Theme Toggle Button -->
    <div class="absolute top-4 right-4">
      <button
        class="btn btn-ghost btn-circle"
        aria-label="Toggle theme"
        @click="toggleTheme"
      >
        <IconWeatherNight v-if="theme === 'light'" class="h-6 w-6" />
        <IconWeatherSunny v-else class="h-6 w-6" />
      </button>
    </div>

    <div class="card w-full max-w-md bg-base-100 shadow-2xl">
      <div class="card-body p-8">
        <!-- Header -->
        <div class="text-center mb-8">
          <h1 class="text-3xl font-bold mb-2">Welcome Back</h1>
          <p class="text-base-content/60">Sign in to your HookFeed account</p>
        </div>

        <!-- Error Alert -->
        <div v-if="error" class="alert alert-error mb-6">
          <span>{{ error }}</span>
        </div>

        <!-- Login Form -->
        <form class="space-y-5" @submit.prevent="handleLogin">
          <!-- Email Field -->
          <div class="form-control">
            <label class="label pb-2">
              <span class="label-text font-medium">Email</span>
            </label>
            <div class="relative">
              <div
                class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-base-content/50 z-10"
              >
                <IconEmail class="h-5 w-5" />
              </div>
              <input
                v-model="email"
                type="email"
                class="input input-bordered w-full pl-10 relative"
                placeholder="you@example.com"
                required
                :disabled="isLoading"
              />
            </div>
          </div>

          <!-- Password Field -->
          <div class="form-control">
            <label class="label pb-2">
              <span class="label-text font-medium">Password</span>
            </label>
            <div class="relative">
              <div
                class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-base-content/50 z-10"
              >
                <IconLock class="h-5 w-5" />
              </div>
              <input
                v-model="password"
                type="password"
                class="input input-bordered w-full pl-10 relative"
                placeholder="Enter your password"
                required
                :disabled="isLoading"
              />
            </div>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            class="btn btn-primary w-full mt-6 gap-2"
            :disabled="isLoading"
          >
            <IconLogin v-if="!isLoading" class="h-5 w-5" />
            <span v-if="isLoading" class="loading loading-spinner loading-sm" />
            {{ isLoading ? "Signing in..." : "Sign In" }}
          </button>
        </form>

        <!-- Footer Links -->
        <div class="divider my-6 text-xs text-base-content/50">OR</div>
        <div class="text-center space-y-3">
          <p class="text-sm text-base-content/70">
            Don't have an account?
            <NuxtLink to="/auth/register" class="link link-primary font-medium">
              Sign up
            </NuxtLink>
          </p>
          <p class="text-sm">
            <NuxtLink
              to="/auth/forgot-password"
              class="link link-hover text-base-content/60"
            >
              Forgot password?
            </NuxtLink>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
