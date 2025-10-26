<script setup lang="ts">
// Icons
import IconInbox from "~icons/mdi/inbox";
import IconRefresh from "~icons/mdi/refresh";
import IconDelete from "~icons/mdi/delete";
import IconCheckCircle from "~icons/mdi/check-circle";

// Page metadata
definePageMeta({
  layout: "default",
  middleware: "auth",
});

// Get feed ID from route
const route = useRoute();
const feedId = computed(() => route.params.id as string);

// State management
const { getFeedById, fetchFeeds } = useAppState();
const { messages, totalMessages, isLoadingMessages, fetchMessages } =
  useFeedMessages();
const { setBreadcrumbs, clearBreadcrumbs } = useBreadcrumbs();

// Get current feed
const currentFeed = computed(() => getFeedById(feedId.value));

// Update breadcrumbs when feed changes
watch(
  currentFeed,
  (feed) => {
    if (feed) {
      setBreadcrumbs([{ label: "All Feeds", to: "/" }, { label: feed.name }]);
    }
  },
  { immediate: true },
);

// Clear breadcrumbs on unmount
onUnmounted(() => {
  clearBreadcrumbs();
});

// Computed: Today's messages count
const todayMessagesCount = computed(() => {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  return messages.value.filter((m) => {
    const receivedAt = m.receivedAt ? new Date(m.receivedAt) : null;
    return receivedAt && receivedAt >= today;
  }).length;
});

// Computed: Separate new and seen messages
const newMessages = computed(() =>
  messages.value.filter((m) => m.state === "new"),
);
const seenMessages = computed(() =>
  messages.value.filter((m) => m.state !== "new"),
);

// Fetch data on mount and when feedId changes
onMounted(async () => {
  await fetchFeeds();

  if (!currentFeed.value) {
    // Redirect to home if feed not found
    navigateTo("/");
    return;
  }

  await fetchMessages({
    feedSlug: currentFeed.value.slug,
  });
});

// Watch feedId changes
watch(feedId, async (newId) => {
  const feed = getFeedById(newId);
  if (!feed) {
    navigateTo("/");
    return;
  }

  await fetchMessages({
    feedSlug: feed.slug,
  });
});

// Refresh messages
const refreshMessages = async () => {
  if (!currentFeed.value) return;

  await fetchMessages({
    feedSlug: currentFeed.value.slug,
    skip: 0,
  });
};

// Bulk delete messages by state
const bulkDeleteByState = async (
  state: "new" | "acknowledged" | "resolved" | "archived" | "all",
) => {
  if (!currentFeed.value) return;

  const confirmMessage =
    state === "all"
      ? `Delete all messages in ${currentFeed.value.name}?`
      : `Delete all "${state}" messages in ${currentFeed.value.name}?`;

  if (!confirm(confirmMessage)) return;

  try {
    const { apiClient } = useApiClient();

    // Fetch all message IDs for the given state
    const stateParam = state === "all" ? "" : `&state=${state}`;
    const response = await apiClient.get<{ items: any[]; total: number }>(
      `/feed-messages?feedSlug=${currentFeed.value.slug}${stateParam}&limit=10000`,
    );

    if (response.data.items.length === 0) {
      alert("No messages to delete");
      return;
    }

    const messageIds = response.data.items.map((m) => m.id);

    // Bulk delete
    await apiClient.post(
      `/feeds/${currentFeed.value.slug}/messages/bulk-delete`,
      {
        messageIds,
      },
    );

    // Refresh messages
    await refreshMessages();
  } catch (err) {
    console.error("Failed to delete messages:", err);
    alert("Failed to delete messages");
  }
};

// Bulk update message states
const bulkUpdateState = async (
  fromState: "new" | "acknowledged" | "resolved" | "archived" | "all",
  toState: "new" | "acknowledged" | "resolved" | "archived",
) => {
  if (!currentFeed.value) return;

  const fromLabel =
    fromState === "all" ? "all messages" : `all "${fromState}" messages`;
  if (!confirm(`Mark ${fromLabel} as "${toState}"?`)) return;

  try {
    const { apiClient } = useApiClient();

    // Fetch all message IDs for the given state
    const stateParam = fromState === "all" ? "" : `&state=${fromState}`;
    const response = await apiClient.get<{ items: any[]; total: number }>(
      `/feed-messages?feedSlug=${currentFeed.value.slug}${stateParam}&limit=10000`,
    );

    if (response.data.items.length === 0) {
      alert("No messages to update");
      return;
    }

    const messageIds = response.data.items.map((m) => m.id);

    // Bulk update state
    await apiClient.post(
      `/feeds/${currentFeed.value.slug}/messages/bulk-state`,
      {
        messageIds,
        state: toState,
      },
    );

    // Refresh messages
    await refreshMessages();
  } catch (err) {
    console.error("Failed to update messages:", err);
    alert("Failed to update messages");
  }
};
</script>

<template>
  <div v-if="currentFeed" class="space-y-6">
    <!-- Header Section -->
    <div class="space-y-2">
      <div class="flex items-center justify-between gap-3">
        <div class="flex items-center gap-3">
          <h1 class="text-3xl font-bold tracking-tight">
            {{ currentFeed.name }}
          </h1>
          <span v-if="currentFeed.unreadCount > 0" class="badge badge-primary">
            {{ currentFeed.unreadCount }} unread
          </span>
        </div>

        <!-- Action Icons -->
        <div class="flex gap-1">
          <button
            class="btn btn-ghost btn-sm btn-square"
            :disabled="isLoadingMessages"
            title="Refresh messages"
            @click="refreshMessages"
          >
            <IconRefresh
              class="h-5 w-5"
              :class="{ 'animate-spin': isLoadingMessages }"
            />
          </button>

          <!-- Mark as dropdown -->
          <div class="dropdown dropdown-end">
            <button
              tabindex="0"
              class="btn btn-ghost btn-sm btn-square"
              title="Mark messages"
            >
              <IconCheckCircle class="h-5 w-5" />
            </button>
            <ul
              tabindex="0"
              class="dropdown-content menu p-2 shadow-lg bg-base-200 rounded-box w-64 mt-2 z-10"
            >
              <li class="menu-title"><span>Mark new as:</span></li>
              <li>
                <a @click="bulkUpdateState('new', 'acknowledged')"
                  >Acknowledged</a
                >
              </li>
              <li>
                <a @click="bulkUpdateState('new', 'resolved')">Resolved</a>
              </li>
              <li>
                <a @click="bulkUpdateState('new', 'archived')">Archived</a>
              </li>

              <li class="menu-title"><span>Mark acknowledged as:</span></li>
              <li>
                <a @click="bulkUpdateState('acknowledged', 'resolved')"
                  >Resolved</a
                >
              </li>
              <li>
                <a @click="bulkUpdateState('acknowledged', 'archived')"
                  >Archived</a
                >
              </li>

              <li class="menu-title"><span>Mark resolved as:</span></li>
              <li>
                <a @click="bulkUpdateState('resolved', 'archived')">Archived</a>
              </li>

              <li class="border-t border-base-300 mt-1 pt-1">
                <a
                  class="font-semibold"
                  @click="bulkUpdateState('all', 'archived')"
                  >Archive all</a
                >
              </li>
            </ul>
          </div>

          <!-- Delete dropdown -->
          <div class="dropdown dropdown-end">
            <button
              tabindex="0"
              class="btn btn-ghost btn-sm btn-square"
              title="Delete messages"
            >
              <IconDelete class="h-5 w-5" />
            </button>
            <ul
              tabindex="0"
              class="dropdown-content menu p-2 shadow-lg bg-base-200 rounded-box w-56 mt-2 z-10"
            >
              <li><a @click="bulkDeleteByState('new')">Delete all new</a></li>
              <li>
                <a @click="bulkDeleteByState('acknowledged')"
                  >Delete all acknowledged</a
                >
              </li>
              <li>
                <a @click="bulkDeleteByState('resolved')"
                  >Delete all resolved</a
                >
              </li>
              <li>
                <a @click="bulkDeleteByState('archived')"
                  >Delete all archived</a
                >
              </li>
              <li class="border-t border-base-300 mt-1 pt-1">
                <a
                  class="text-error font-semibold"
                  @click="bulkDeleteByState('all')"
                  >Delete all</a
                >
              </li>
            </ul>
          </div>
        </div>
      </div>

      <p class="text-base-content/60">
        {{ currentFeed.description }}
      </p>

      <!-- Stats Row -->
      <div class="flex items-center gap-4 text-sm">
        <div class="flex items-center gap-1.5">
          <span class="font-semibold text-primary">{{ totalMessages }}</span>
          <span class="text-base-content/60">total messages</span>
        </div>
        <div class="flex items-center gap-1.5">
          <span class="font-semibold">{{ todayMessagesCount }}</span>
          <span class="text-base-content/60">today</span>
        </div>
      </div>
    </div>

    <!-- Messages Container -->
    <div class="space-y-4">
      <!-- Loading State -->
      <div
        v-if="isLoadingMessages"
        class="flex flex-col items-center justify-center py-16 px-4 bg-base-200 rounded-lg"
      >
        <span class="loading loading-spinner loading-lg" />
        <p class="text-base-content/60 mt-4">Loading messages...</p>
      </div>

      <!-- Empty State -->
      <div
        v-else-if="messages.length === 0"
        class="flex flex-col items-center justify-center py-16 px-4 bg-base-200 rounded-lg"
      >
        <IconInbox class="h-16 w-16 text-base-content/20 mb-4" />
        <h3 class="text-xl font-semibold mb-2">
          No Messages in {{ currentFeed.name }}
        </h3>
        <p class="text-base-content/60 text-center max-w-md">
          Messages sent to this feed will appear here. Use the webhook URL above
          to send messages.
        </p>
      </div>

      <!-- Message Cards -->
      <template v-else>
        <!-- New Messages -->
        <MessageCard
          v-for="message in newMessages"
          :key="message.id"
          :message="message"
        />

        <!-- Separator (only show if there are both new and seen messages) -->
        <div
          v-if="newMessages.length > 0 && seenMessages.length > 0"
          class="flex items-center gap-4 py-4"
        >
          <div class="flex-1 border-t border-base-300" />
          <span
            class="text-sm font-medium text-base-content/50 uppercase tracking-wide"
            >Seen</span
          >
          <div class="flex-1 border-t border-base-300" />
        </div>

        <!-- Seen Messages -->
        <MessageCard
          v-for="message in seenMessages"
          :key="message.id"
          :message="message"
        />
      </template>
    </div>
  </div>
</template>
