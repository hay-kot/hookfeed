<script setup lang="ts">
// Icons
import IconInbox from '~icons/mdi/inbox'
import IconChevronLeft from '~icons/mdi/chevron-left'
import IconDotsVertical from '~icons/mdi/dots-vertical'
import IconRefresh from '~icons/mdi/refresh'
import IconCopy from '~icons/mdi/content-copy'

// Page metadata
definePageMeta({
  layout: 'default',
})

// Get feed ID from route
const route = useRoute()
const feedId = computed(() => route.params.id as string)

// State management
const { getFeedById } = useAppState()
const { generateMessages } = useMockData()

// Get current feed
const currentFeed = computed(() => getFeedById(feedId.value))

// Mock messages data - generate messages for this specific feed
const messages = ref(generateMessages(12, currentFeed.value?.slug))

// Filters
const selectedLevel = ref<string>('all')
const selectedState = ref<string>('all')
const searchQuery = ref<string>('')

// Computed: Filtered messages
const filteredMessages = computed(() => {
  let filtered = messages.value

  if (selectedLevel.value !== 'all') {
    filtered = filtered.filter(m => m.level === selectedLevel.value)
  }

  if (selectedState.value !== 'all') {
    filtered = filtered.filter(m => m.state === selectedState.value)
  }

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(m =>
      m.title?.toLowerCase().includes(query) ||
      m.message?.toLowerCase().includes(query),
    )
  }

  return filtered
})

// Computed: Today's messages count
const todayMessagesCount = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return messages.value.filter(m => {
    const receivedAt = m.receivedAt ? new Date(m.receivedAt) : null
    return receivedAt && receivedAt >= today
  }).length
})

// Refresh messages
const refreshMessages = () => {
  messages.value = generateMessages(12, currentFeed.value?.slug)
}

// Copy webhook URL
const copyWebhookUrl = async () => {
  if (!currentFeed.value) return
  const url = `https://hookfeed.example.com/hooks/${currentFeed.value.slug}`
  try {
    await navigator.clipboard.writeText(url)
    // Could add a toast notification here
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

// Handle feed not found
onMounted(() => {
  if (!currentFeed.value) {
    // Redirect to home if feed not found
    navigateTo('/')
  }
})
</script>

<template>
  <div v-if="currentFeed" class="space-y-4">
    <!-- Compact Header -->
    <div class="flex items-center justify-between gap-4">
      <!-- Left side: Back button + Title -->
      <div class="flex items-center gap-3 min-w-0 flex-1">
        <NuxtLink
          to="/"
          class="btn btn-ghost btn-sm btn-circle flex-shrink-0"
          title="Back to All Feeds"
        >
          <IconChevronLeft class="h-5 w-5" />
        </NuxtLink>
        <div class="min-w-0 flex-1">
          <div class="flex items-center gap-2">
            <h1 class="text-2xl font-bold tracking-tight truncate">
              {{ currentFeed.name }}
            </h1>
            <span
              v-if="currentFeed.unreadCount > 0"
              class="badge badge-primary badge-sm flex-shrink-0"
            >
              {{ currentFeed.unreadCount }}
            </span>
          </div>
          <p class="text-sm text-base-content/60 truncate">
            {{ currentFeed.description }}
          </p>
        </div>
      </div>

      <!-- Right side: Stats + Actions -->
      <div class="flex items-center gap-3 flex-shrink-0">
        <!-- Compact inline stats -->
        <div class="hidden md:flex items-center gap-4 text-sm">
          <div class="text-center">
            <div class="font-bold text-lg text-primary">{{ messages.length }}</div>
            <div class="text-xs text-base-content/60">Total</div>
          </div>
          <div class="text-center">
            <div class="font-bold text-lg">{{ currentFeed.unreadCount }}</div>
            <div class="text-xs text-base-content/60">Unread</div>
          </div>
          <div class="text-center">
            <div class="font-bold text-lg">{{ todayMessagesCount }}</div>
            <div class="text-xs text-base-content/60">Today</div>
          </div>
        </div>

        <!-- Feed actions dropdown -->
        <div class="dropdown dropdown-end">
          <button tabindex="0" class="btn btn-ghost btn-sm btn-circle">
            <IconDotsVertical class="h-5 w-5" />
          </button>
          <ul tabindex="0" class="dropdown-content menu p-2 shadow-lg bg-base-200 rounded-box w-52 mt-2 z-10">
            <li><a>Mark all as read</a></li>
            <li><a>Export messages</a></li>
            <li><a>Feed settings</a></li>
            <li class="border-t border-base-300 mt-1 pt-1">
              <a class="text-error">Delete feed</a>
            </li>
          </ul>
        </div>
      </div>
    </div>

    <!-- Filters and Webhook URL combined -->
    <div class="bg-base-200 p-3 rounded-lg space-y-2">
      <!-- Filters row -->
      <div class="flex flex-wrap gap-2 items-center justify-between">
        <div class="flex flex-wrap gap-2 flex-1">
          <select v-model="selectedLevel" class="select select-bordered select-sm">
            <option value="all">
              All Levels
            </option>
            <option value="info">
              Info
            </option>
            <option value="warning">
              Warning
            </option>
            <option value="error">
              Error
            </option>
            <option value="success">
              Success
            </option>
            <option value="debug">
              Debug
            </option>
          </select>

          <select v-model="selectedState" class="select select-bordered select-sm">
            <option value="all">
              All States
            </option>
            <option value="new">
              New
            </option>
            <option value="acknowledged">
              Acknowledged
            </option>
            <option value="resolved">
              Resolved
            </option>
            <option value="archived">
              Archived
            </option>
          </select>

          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search..."
            class="input input-bordered input-sm w-32 lg:w-48"
          >

          <span class="badge badge-ghost self-center">
            {{ filteredMessages.length }}
          </span>
        </div>

        <button class="btn btn-sm btn-primary gap-1" @click="refreshMessages">
          <IconRefresh class="h-4 w-4" />
          <span class="hidden sm:inline">Refresh</span>
        </button>
      </div>

      <!-- Webhook URL row -->
      <div class="flex items-center gap-2 text-xs">
        <code class="flex-1 text-base-content/60 truncate">
          https://hookfeed.example.com/hooks/{{ currentFeed.slug }}
        </code>
        <button
          class="btn btn-xs btn-ghost gap-1 flex-shrink-0"
          @click="copyWebhookUrl"
        >
          <IconCopy class="h-3 w-3" />
          Copy
        </button>
      </div>
    </div>

    <!-- Messages Container -->
    <div class="space-y-4">
      <!-- Empty State -->
      <div
        v-if="filteredMessages.length === 0"
        class="flex flex-col items-center justify-center py-16 px-4 bg-base-200 rounded-lg"
      >
        <IconInbox class="h-16 w-16 text-base-content/20 mb-4" />
        <h3 class="text-xl font-semibold mb-2">
          {{ messages.length === 0 ? `No Messages in ${currentFeed.name}` : 'No Messages Found' }}
        </h3>
        <p class="text-base-content/60 text-center max-w-md">
          {{ messages.length === 0
            ? 'Messages sent to this feed will appear here. Use the webhook URL above to send messages.'
            : 'No messages match the selected filters. Try adjusting your filter criteria.' }}
        </p>
      </div>

      <!-- Message Cards -->
      <MessageCard
        v-for="message in filteredMessages"
        :key="message.id"
        :message="message"
      />
    </div>
  </div>
</template>
