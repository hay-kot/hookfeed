<script setup lang="ts">
// Icons
import IconMenu from "~icons/mdi/menu";
import IconClose from "~icons/mdi/close";
import IconWeatherSunny from "~icons/mdi/weather-sunny";
import IconWeatherNight from "~icons/mdi/weather-night";
import IconCog from "~icons/mdi/cog";
import IconViewDashboard from "~icons/mdi/view-dashboard";
import IconLogout from "~icons/mdi/logout";
import IconAccount from "~icons/mdi/account-circle";

// State management
const { theme, initTheme, toggleTheme } = useTheme();
const { feeds, fetchFeeds } = useAppState();
const { user, logout } = useAuth();

// Handle logout
const handleLogout = () => {
  logout();
  navigateTo('/auth/login');
};

// Sidebar state for mobile
const isSidebarOpen = ref(false);

// Initialize theme and fetch feeds on mount
onMounted(async () => {
  initTheme();
  await fetchFeeds();
});

// Close sidebar when navigating (mobile)
const closeSidebar = () => {
  isSidebarOpen.value = false;
};

// Toggle sidebar (mobile)
const toggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value;
};
</script>

<template>
  <div class="h-screen w-screen flex flex-col overflow-hidden bg-base-100">
    <!-- Top Bar -->
    <header
      class="navbar bg-base-200 border-b border-base-300 px-4 h-16 flex-shrink-0 z-40"
    >
      <!-- Mobile menu button -->
      <div class="flex-none lg:hidden">
        <button
          class="btn btn-ghost btn-circle"
          aria-label="Toggle menu"
          @click="toggleSidebar"
        >
          <IconMenu v-if="!isSidebarOpen" class="h-6 w-6" />
          <IconClose v-else class="h-6 w-6" />
        </button>
      </div>

      <!-- Logo/Title -->
      <div class="flex-1">
        <NuxtLink
          to="/"
          class="text-xl font-bold tracking-tight hover:text-primary transition-colors"
        >
          HookFeed
        </NuxtLink>
      </div>

      <!-- Right side actions -->
      <div class="flex-none gap-2">
        <!-- Theme toggle -->
        <button
          class="btn btn-ghost btn-circle"
          aria-label="Toggle theme"
          @click="toggleTheme"
        >
          <IconWeatherNight v-if="theme === 'light'" class="h-6 w-6" />
          <IconWeatherSunny v-else class="h-6 w-6" />
        </button>
      </div>
    </header>

    <div class="flex flex-1 overflow-hidden">
      <!-- Sidebar Overlay (mobile) -->
      <div
        v-if="isSidebarOpen"
        class="fixed inset-0 bg-black/50 z-30 lg:hidden"
        @click="closeSidebar"
      />

      <!-- Sidebar -->
      <aside
        class="w-64 bg-base-200 border-r border-base-300 flex flex-col fixed lg:static inset-y-0 left-0 z-40 transition-transform duration-300 lg:translate-x-0"
        :class="isSidebarOpen ? 'translate-x-0' : '-translate-x-full'"
      >
        <!-- Sidebar content wrapper -->
        <div class="flex flex-col h-full pt-4">
          <!-- Top section: All Feeds + Feed List -->
          <nav class="flex-1 px-4 space-y-4 overflow-y-auto">
            <!-- All Feeds Link -->
            <div class="space-y-2">
              <NuxtLink
                to="/"
                class="flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-base-300 transition-colors font-medium"
                active-class="bg-primary text-primary-content hover:bg-primary"
                @click="closeSidebar"
              >
                <IconViewDashboard class="h-5 w-5" />
                <span>All Feeds</span>
              </NuxtLink>
            </div>

            <!-- Feeds Section -->
            <div class="space-y-2">
              <h2
                class="px-4 text-xs font-semibold uppercase tracking-wider text-base-content/60"
              >
                Feeds
              </h2>

              <!-- Feed List -->
              <div class="space-y-1">
                <NuxtLink
                  v-for="feed in feeds"
                  :key="feed.id"
                  :to="`/feeds/${feed.id}`"
                  class="flex items-center justify-between px-4 py-2.5 rounded-lg hover:bg-base-300 transition-colors group"
                  active-class="bg-base-300"
                  @click="closeSidebar"
                >
                  <span class="text-sm font-medium truncate">{{
                    feed.name
                  }}</span>

                  <!-- Unread badge -->
                  <span
                    v-if="feed.unreadCount > 0"
                    class="badge badge-primary badge-sm flex-shrink-0"
                  >
                    {{ feed.unreadCount }}
                  </span>
                </NuxtLink>
              </div>
            </div>
          </nav>

          <!-- Bottom section: User & Settings -->
          <div class="px-4 py-4 border-t border-base-300 space-y-2">
            <!-- User info -->
            <div v-if="user" class="flex items-center gap-3 px-4 py-2 text-sm">
              <IconAccount class="h-5 w-5 text-base-content/60" />
              <div class="flex-1 truncate">
                <div class="font-medium truncate">{{ user.username }}</div>
                <div class="text-xs text-base-content/60 truncate">{{ user.email }}</div>
              </div>
            </div>

            <NuxtLink
              to="/settings"
              class="flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-base-300 transition-colors font-medium"
              active-class="bg-base-300"
              @click="closeSidebar"
            >
              <IconCog class="h-5 w-5" />
              <span>Settings</span>
            </NuxtLink>

            <button
              class="flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-base-300 transition-colors font-medium w-full text-left text-error"
              @click="handleLogout"
            >
              <IconLogout class="h-5 w-5" />
              <span>Logout</span>
            </button>
          </div>
        </div>
      </aside>

      <!-- Main Content Area -->
      <main class="flex-1 overflow-y-auto bg-base-100">
        <div class="container mx-auto px-4 py-6 max-w-7xl">
          <!-- Page content will be rendered here -->
          <slot />
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped>
/* Ensure sidebar animation works smoothly */
aside {
  top: 64px; /* Height of navbar (h-16 = 4rem = 64px) */
  height: calc(100vh - 64px);
}

@media (min-width: 1024px) {
  aside {
    position: static;
    height: auto;
  }
}
</style>
