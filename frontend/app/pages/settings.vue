<script setup lang="ts">
// Icons
import IconPalette from '~icons/mdi/palette'
import IconBell from '~icons/mdi/bell'
import IconDatabase from '~icons/mdi/database'
import IconInformation from '~icons/mdi/information'

// Page metadata
definePageMeta({
  layout: 'default',
})

// State management
const { theme, toggleTheme } = useTheme()

// Settings state (placeholder - will be connected to real settings later)
const notificationSettings = ref({
  enabled: true,
  sound: false,
  desktop: true,
})

const retentionSettings = ref({
  maxMessages: 10000,
  maxAgeDays: 90,
})

const version = ref('1.0.0')
</script>

<template>
  <div class="space-y-6 max-w-4xl">
    <!-- Page Header -->
    <div>
      <h1 class="text-3xl font-bold tracking-tight">
        Settings
      </h1>
      <p class="text-base-content/60 mt-1">
        Manage your HookFeed preferences and configuration
      </p>
    </div>

    <!-- Settings Sections -->
    <div class="space-y-6">
      <!-- Appearance Settings -->
      <div class="card bg-base-200 shadow-sm">
        <div class="card-body">
          <div class="flex items-center gap-3 mb-4">
            <IconPalette class="h-6 w-6 text-primary" />
            <h2 class="card-title">Appearance</h2>
          </div>

          <div class="form-control">
            <label class="label cursor-pointer justify-start gap-4">
              <input
                type="checkbox"
                class="toggle toggle-primary"
                :checked="theme === 'dark'"
                @change="toggleTheme"
              >
              <div>
                <span class="label-text font-semibold">Dark Mode</span>
                <p class="text-sm text-base-content/60">
                  Use dark theme across the application
                </p>
              </div>
            </label>
          </div>
        </div>
      </div>

      <!-- Notification Settings -->
      <div class="card bg-base-200 shadow-sm">
        <div class="card-body">
          <div class="flex items-center gap-3 mb-4">
            <IconBell class="h-6 w-6 text-primary" />
            <h2 class="card-title">Notifications</h2>
          </div>

          <div class="space-y-4">
            <div class="form-control">
              <label class="label cursor-pointer justify-start gap-4">
                <input
                  v-model="notificationSettings.enabled"
                  type="checkbox"
                  class="toggle toggle-primary"
                >
                <div>
                  <span class="label-text font-semibold">Enable Notifications</span>
                  <p class="text-sm text-base-content/60">
                    Show notifications for new messages
                  </p>
                </div>
              </label>
            </div>

            <div class="form-control">
              <label class="label cursor-pointer justify-start gap-4">
                <input
                  v-model="notificationSettings.sound"
                  type="checkbox"
                  class="toggle toggle-primary"
                  :disabled="!notificationSettings.enabled"
                >
                <div>
                  <span class="label-text font-semibold">Notification Sound</span>
                  <p class="text-sm text-base-content/60">
                    Play sound when receiving new messages
                  </p>
                </div>
              </label>
            </div>

            <div class="form-control">
              <label class="label cursor-pointer justify-start gap-4">
                <input
                  v-model="notificationSettings.desktop"
                  type="checkbox"
                  class="toggle toggle-primary"
                  :disabled="!notificationSettings.enabled"
                >
                <div>
                  <span class="label-text font-semibold">Desktop Notifications</span>
                  <p class="text-sm text-base-content/60">
                    Show browser notifications for new messages
                  </p>
                </div>
              </label>
            </div>
          </div>
        </div>
      </div>

      <!-- Data & Storage Settings -->
      <div class="card bg-base-200 shadow-sm">
        <div class="card-body">
          <div class="flex items-center gap-3 mb-4">
            <IconDatabase class="h-6 w-6 text-primary" />
            <h2 class="card-title">Data & Storage</h2>
          </div>

          <div class="space-y-4">
            <div class="form-control">
              <label class="label">
                <span class="label-text font-semibold">Maximum Messages per Feed</span>
              </label>
              <input
                v-model.number="retentionSettings.maxMessages"
                type="number"
                class="input input-bordered w-full max-w-xs"
                min="100"
                step="100"
              >
              <label class="label">
                <span class="label-text-alt text-base-content/60">
                  Older messages will be automatically deleted
                </span>
              </label>
            </div>

            <div class="form-control">
              <label class="label">
                <span class="label-text font-semibold">Message Retention (Days)</span>
              </label>
              <input
                v-model.number="retentionSettings.maxAgeDays"
                type="number"
                class="input input-bordered w-full max-w-xs"
                min="1"
                step="1"
              >
              <label class="label">
                <span class="label-text-alt text-base-content/60">
                  Messages older than this will be deleted
                </span>
              </label>
            </div>

            <div class="divider" />

            <div class="flex gap-2">
              <button class="btn btn-error btn-outline btn-sm">
                Clear All Messages
              </button>
              <button class="btn btn-warning btn-outline btn-sm">
                Clear Archived Messages
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- About Section -->
      <div class="card bg-base-200 shadow-sm">
        <div class="card-body">
          <div class="flex items-center gap-3 mb-4">
            <IconInformation class="h-6 w-6 text-primary" />
            <h2 class="card-title">About</h2>
          </div>

          <div class="space-y-2">
            <div class="flex justify-between">
              <span class="text-base-content/60">Version</span>
              <span class="font-mono font-semibold">{{ version }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-base-content/60">API Status</span>
              <span class="badge badge-success">Connected</span>
            </div>
            <div class="flex justify-between">
              <span class="text-base-content/60">WebSocket Status</span>
              <span class="badge badge-success">Connected</span>
            </div>
          </div>

          <div class="divider" />

          <div class="flex gap-2">
            <a
              href="https://github.com/hookfeed/hookfeed"
              target="_blank"
              class="btn btn-sm btn-ghost"
            >
              Documentation
            </a>
            <a
              href="https://github.com/hookfeed/hookfeed/issues"
              target="_blank"
              class="btn btn-sm btn-ghost"
            >
              Report Issue
            </a>
          </div>
        </div>
      </div>

      <!-- Save Button -->
      <div class="flex justify-end gap-2 pt-4">
        <button class="btn btn-ghost">
          Reset to Defaults
        </button>
        <button class="btn btn-primary">
          Save Settings
        </button>
      </div>
    </div>
  </div>
</template>
