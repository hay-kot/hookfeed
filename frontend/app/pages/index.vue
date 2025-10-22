<script setup lang="ts">
// Icons
import IconInbox from "~icons/mdi/inbox";
import IconRefresh from "~icons/mdi/refresh";

// Page metadata
definePageMeta({
  layout: 'default',
  middleware: 'auth',
});

// State management
const { feeds, fetchFeeds } = useAppState();
const { messages, isLoadingMessages, fetchMessages } = useFeedMessages();

// Filters
const selectedState = ref<string | undefined>(undefined);
const searchQuery = ref<string>("");

// Computed: Total message count across all feeds
const totalFeeds = computed(() => feeds.value.length);

// Fetch data on mount
onMounted(async () => {
  await fetchFeeds();
  await fetchMessages();
});

// Watch filters and refetch messages
watch([selectedState, searchQuery], async () => {
  await fetchMessages({
    state: selectedState.value as any,
    q: searchQuery.value || undefined,
    skip: 0,
  });
});

// Refresh messages
const refreshMessages = async () => {
  await fetchMessages({
    state: selectedState.value as any,
    q: searchQuery.value || undefined,
    skip: 0,
  });
};
</script>

<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold tracking-tight">All Feeds</h1>
        <p class="text-base-content/60 mt-1">
          Viewing messages from all {{ totalFeeds }} feeds
        </p>
      </div>
    </div>

    <!-- Filters/Actions Bar -->
    <div
      class="flex flex-wrap gap-3 items-center justify-between bg-base-200 p-4 rounded-lg"
    >
      <div class="flex flex-wrap gap-2 flex-1">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search messages..."
          class="input input-bordered input-sm flex-1 min-w-[200px]"
        />

        <select
          v-model="selectedState"
          class="select select-bordered select-sm"
        >
          <option :value="undefined">All States</option>
          <option value="new">New</option>
          <option value="acknowledged">Acknowledged</option>
          <option value="resolved">Resolved</option>
          <option value="archived">Archived</option>
        </select>

        <span class="badge badge-ghost self-center">
          {{ messages.length }} messages
        </span>
      </div>

      <div class="flex gap-2">
        <button
          class="btn btn-sm btn-primary gap-2"
          @click="refreshMessages"
          :disabled="isLoadingMessages"
        >
          <IconRefresh class="h-4 w-4" :class="{ 'animate-spin': isLoadingMessages }" />
          Refresh
        </button>
      </div>
    </div>

    <!-- Messages Container -->
    <div class="space-y-4">
      <!-- Loading State -->
      <div
        v-if="isLoadingMessages"
        class="flex flex-col items-center justify-center py-16 px-4 bg-base-200 rounded-lg"
      >
        <span class="loading loading-spinner loading-lg"></span>
        <p class="text-base-content/60 mt-4">Loading messages...</p>
      </div>

      <!-- Empty State -->
      <div
        v-else-if="messages.length === 0"
        class="flex flex-col items-center justify-center py-16 px-4 bg-base-200 rounded-lg"
      >
        <IconInbox class="h-16 w-16 text-base-content/20 mb-4" />
        <h3 class="text-xl font-semibold mb-2">No Messages Found</h3>
        <p class="text-base-content/60 text-center max-w-md">
          Messages from all your feeds will appear here. Send a webhook to get started!
        </p>
      </div>

      <!-- Message Cards -->
      <MessageCard
        v-else
        v-for="message in messages"
        :key="message.id"
        :message="message"
      />
    </div>
  </div>
</template>
