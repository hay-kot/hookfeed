<script setup lang="ts">
import IconEmailFast from "~icons/mdi/email-fast";
import IconEmail from "~icons/mdi/email";
import IconArrowLeft from "~icons/mdi/arrow-left";
import IconWeatherSunny from "~icons/mdi/weather-sunny";
import IconWeatherNight from "~icons/mdi/weather-night";

// Disable default layout for auth pages
definePageMeta({
  layout: false,
});

const { theme, initTheme, toggleTheme } = useTheme();

// Initialize theme on mount
onMounted(() => {
  initTheme();
});

// Form state
const email = ref("");
const error = ref("");
const success = ref(false);
const isLoading = ref(false);

// Handle forgot password
const handleForgotPassword = async () => {
  error.value = "";
  isLoading.value = true;

  // TODO: Implement password reset
  await new Promise((resolve) => setTimeout(resolve, 1000));
  success.value = true;
  isLoading.value = false;
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
          <h1 class="text-3xl font-bold mb-2">Forgot Password</h1>
          <p class="text-base-content/60">
            Enter your email to receive a password reset link
          </p>
        </div>

        <!-- Success Alert -->
        <div v-if="success" class="alert alert-success mb-6">
          <span>Password reset link sent! Check your email.</span>
        </div>

        <!-- Error Alert -->
        <div v-if="error" class="alert alert-error mb-6">
          <span>{{ error }}</span>
        </div>

        <!-- Forgot Password Form -->
        <form v-if="!success" class="space-y-5" @submit.prevent="handleForgotPassword">
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

          <!-- Submit Button -->
          <button
            type="submit"
            class="btn btn-primary w-full mt-6 gap-2"
            :disabled="isLoading"
          >
            <IconEmailFast v-if="!isLoading" class="h-5 w-5" />
            <span v-if="isLoading" class="loading loading-spinner loading-sm" />
            {{ isLoading ? "Sending..." : "Send Reset Link" }}
          </button>
        </form>

        <!-- Back to Login Button (shown after success) -->
        <NuxtLink
          v-if="success"
          to="/auth/login"
          class="btn btn-primary w-full gap-2"
        >
          <IconArrowLeft class="h-5 w-5" />
          Back to Login
        </NuxtLink>

        <!-- Footer Links -->
        <div class="divider my-6 text-xs text-base-content/50">OR</div>
        <div class="text-center">
          <p class="text-sm text-base-content/70">
            Remember your password?
            <NuxtLink to="/auth/login" class="link link-primary font-medium">
              Sign in
            </NuxtLink>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
