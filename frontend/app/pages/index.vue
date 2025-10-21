<script setup lang="ts">
// Icons
import IconInbox from "~icons/mdi/inbox";
import IconRefresh from "~icons/mdi/refresh";

// Page metadata
definePageMeta({
  layout: "default",
});

// State management
const { feeds } = useAppState();
const { generateMessages } = useMockData();

// Mock messages data - generate 15 sample messages
const messages = ref(generateMessages(15));

// Filters
const selectedLevel = ref<string>("all");
const selectedState = ref<string>("all");

// Computed: Total message count across all feeds
const totalFeeds = computed(() => feeds.value.length);

// Computed: Filtered messages
const filteredMessages = computed(() => {
  let filtered = messages.value;

  if (selectedLevel.value !== "all") {
    filtered = filtered.filter((m) => m.level === selectedLevel.value);
  }

  if (selectedState.value !== "all") {
    filtered = filtered.filter((m) => m.state === selectedState.value);
  }

  return filtered;
});

// Refresh messages
const refreshMessages = () => {
  messages.value = generateMessages(15);
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
      <div class="flex flex-wrap gap-2">
        <select
          v-model="selectedLevel"
          class="select select-bordered select-sm"
        >
          <option value="all">All Levels</option>
          <option value="info">Info</option>
          <option value="warning">Warning</option>
          <option value="error">Error</option>
          <option value="success">Success</option>
          <option value="debug">Debug</option>
        </select>

        <select
          v-model="selectedState"
          class="select select-bordered select-sm"
        >
          <option value="all">All States</option>
          <option value="new">New</option>
          <option value="acknowledged">Acknowledged</option>
          <option value="resolved">Resolved</option>
          <option value="archived">Archived</option>
        </select>

        <span class="badge badge-ghost self-center">
          {{ filteredMessages.length }} messages
        </span>
      </div>

      <div class="flex gap-2">
        <button class="btn btn-sm btn-ghost">Mark All as Read</button>
        <button class="btn btn-sm btn-primary gap-2" @click="refreshMessages">
          <IconRefresh class="h-4 w-4" />
          Refresh
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
        <h3 class="text-xl font-semibold mb-2">No Messages Found</h3>
        <p class="text-base-content/60 text-center max-w-md">
          {{
            messages.length === 0
              ? "Messages from all your feeds will appear here. Send a webhook to get started!"
              : "No messages match the selected filters. Try adjusting your filter criteria."
          }}
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
