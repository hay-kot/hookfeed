<script setup lang="ts">
/**
 * Message card component with tabbed interface for viewing message details
 */

import type { FeedMessage } from "~~/lib/api/types/data-contracts";
import IconBell from "~icons/mdi/bell";
import IconEyeCheck from "~icons/mdi/eye-check";
import IconCheckCircle from "~icons/mdi/check-circle";
import IconArchive from "~icons/mdi/archive";
import IconClock from "~icons/mdi/clock";
import IconChevronDown from "~icons/mdi/chevron-down";
import IconChevronUp from "~icons/mdi/chevron-up";

interface Props {
  message: FeedMessage;
}

const props = defineProps<Props>();
const { updateMessageState, deleteMessage } = useFeedMessages();

// Data is already JSON objects from the backend
const rawRequest = computed(() => props.message.rawRequest);
const rawHeaders = computed(() => props.message.rawHeaders);
const rawQueryParams = computed(() => props.message.rawQueryParams);

// Priority colors (1-5, where 1 is highest priority)
const priorityColors = {
  1: "badge-error", // Critical
  2: "badge-warning", // High
  3: "badge-info", // Medium
  4: "badge-success", // Low
  5: "badge-ghost", // Minimal
};

const priorityLabels = {
  1: "Critical",
  2: "High",
  3: "Medium",
  4: "Low",
  5: "Minimal",
};

// State icons
const stateIcons = {
  new: IconBell,
  acknowledged: IconEyeCheck,
  resolved: IconCheckCircle,
  archived: IconArchive,
};

// State colors
const stateColors = {
  new: "text-primary",
  acknowledged: "text-info",
  resolved: "text-success",
  archived: "text-base-content/40",
};

// State border colors
const stateBorderColors = {
  new: "border-l-primary",
  acknowledged: "border-l-info",
  resolved: "border-l-success",
  archived: "border-l-base-content/40",
};

// Expanded state for card details
const isExpanded = ref(false);
const activeTab = ref<"request" | "headers" | "queryParams" | "logs">(
  "request",
);

// Format timestamp
const formatTimestamp = (timestamp?: string) => {
  if (!timestamp) return "N/A";
  const date = new Date(timestamp);
  return date.toLocaleString();
};

// Format relative time
const formatRelativeTime = (timestamp?: string) => {
  if (!timestamp) return "";
  const date = new Date(timestamp);
  const now = new Date();
  const diff = now.getTime() - date.getTime();

  const minutes = Math.floor(diff / 60000);
  const hours = Math.floor(diff / 3600000);
  const days = Math.floor(diff / 86400000);

  if (minutes < 1) return "just now";
  if (minutes < 60) return `${minutes}m ago`;
  if (hours < 24) return `${hours}h ago`;
  return `${days}d ago`;
};

const toggleExpand = () => {
  isExpanded.value = !isExpanded.value;
};

// State management
const showStateMenu = ref(false);

// Copy message ID to clipboard
const copyId = async () => {
  if (!props.message.id) return;
  try {
    await navigator.clipboard.writeText(props.message.id);
    // Could add a toast notification here
  } catch (err) {
    console.error("Failed to copy ID:", err);
  }
};

// Mark message as acknowledged (read)
const markAsRead = async () => {
  if (!props.message.id) return;
  await updateMessageState(props.message.id, "acknowledged");
};

// Change message state
const changeState = async (
  newState: "new" | "acknowledged" | "resolved" | "archived",
) => {
  if (!props.message.id) return;
  await updateMessageState(props.message.id, newState);
  showStateMenu.value = false;
};

// Delete message
const handleDelete = async () => {
  if (!props.message.id) return;
  if (confirm("Are you sure you want to delete this message?")) {
    await deleteMessage(props.message.id);
  }
};
</script>

<template>
  <div
    class="card bg-base-200 shadow-sm border border-base-300 hover:shadow-md transition-shadow border-l-8"
    :class="stateBorderColors[message.state as keyof typeof stateBorderColors]"
  >
    <!-- Card Header - Always visible -->
    <div class="card-body p-4">
      <div class="flex items-start gap-4 cursor-pointer" @click="toggleExpand">
        <!-- State Icon -->
        <div class="flex-shrink-0 pt-1">
          <component
            :is="
              stateIcons[message.state as keyof typeof stateIcons] || IconBell
            "
            class="h-5 w-5"
            :class="stateColors[message.state as keyof typeof stateColors]"
          />
        </div>

        <!-- Main Content -->
        <div class="flex-1 min-w-0">
          <!-- Title & Priority -->
          <div class="flex items-start justify-between gap-3 mb-2">
            <h3 class="font-semibold text-base leading-tight">
              {{ message.title || "Untitled Message" }}
            </h3>
            <span
              class="badge badge-sm flex-shrink-0"
              :class="
                priorityColors[
                  message.priority as keyof typeof priorityColors
                ] || 'badge-ghost'
              "
            >
              {{
                priorityLabels[
                  message.priority as keyof typeof priorityLabels
                ] || message.priority
              }}
            </span>
          </div>

          <!-- Message -->
          <p
            v-if="message.message"
            class="text-sm text-base-content/80 mb-3 line-clamp-2"
          >
            {{ message.message }}
          </p>

          <!-- Metadata Row -->
          <div
            class="flex flex-wrap items-center gap-3 text-xs text-base-content/60"
          >
            <div class="flex items-center gap-1">
              <IconClock class="h-3.5 w-3.5" />
              <span>{{ formatRelativeTime(message.receivedAt) }}</span>
            </div>
            <span class="text-base-content/30">•</span>
            <span class="capitalize">{{ message.state }}</span>
            <span v-if="message.feedSlug" class="text-base-content/30">•</span>
            <span v-if="message.feedSlug" class="font-mono">{{
              message.feedSlug
            }}</span>
          </div>
        </div>

        <!-- Expand Button -->
        <button
          class="btn btn-ghost btn-sm btn-circle flex-shrink-0"
          @click.stop="toggleExpand"
        >
          <IconChevronDown v-if="!isExpanded" class="h-5 w-5" />
          <IconChevronUp v-else class="h-5 w-5" />
        </button>
      </div>

      <!-- Expanded Details -->
      <div class="expand" :class="isExpanded ? 'expand-open' : ''">
        <div>
          <div class="mt-4 pt-4 border-t border-base-300">
            <!-- Tabs -->
            <div role="tablist" class="tabs tabs-bordered mb-4">
              <button
                role="tab"
                class="tab"
                :class="{ 'tab-active': activeTab === 'request' }"
                @click="activeTab = 'request'"
              >
                Body
              </button>
              <button
                role="tab"
                class="tab"
                :class="{ 'tab-active': activeTab === 'headers' }"
                @click="activeTab = 'headers'"
              >
                Headers
              </button>
              <button
                role="tab"
                class="tab"
                :class="{ 'tab-active': activeTab === 'queryParams' }"
                @click="activeTab = 'queryParams'"
              >
                Query Params
              </button>
              <button
                role="tab"
                class="tab"
                :class="{ 'tab-active': activeTab === 'logs' }"
                @click="activeTab = 'logs'"
              >
                Logs
                <span
                  v-if="message.logs && message.logs.length > 0"
                  class="badge badge-xs ml-1"
                >
                  {{ message.logs.length }}
                </span>
              </button>
            </div>

            <!-- Tab Content -->
            <div class="min-h-[200px]">
              <!-- Body Tab -->
              <div v-if="activeTab === 'request'">
                <JsonViewer v-if="rawRequest" :data="rawRequest" />
                <div v-else class="text-center text-base-content/40 py-8">
                  No raw request data available
                </div>
              </div>

              <!-- Headers Tab -->
              <div v-if="activeTab === 'headers'">
                <JsonViewer v-if="rawHeaders" :data="rawHeaders" />
                <div v-else class="text-center text-base-content/40 py-8">
                  No headers available
                </div>
              </div>

              <!-- Query Params Tab -->
              <div v-if="activeTab === 'queryParams'">
                <JsonViewer v-if="rawQueryParams" :data="rawQueryParams" />
                <div v-else class="text-center text-base-content/40 py-8">
                  No query parameters available
                </div>
              </div>

              <!-- Logs Tab -->
              <div v-if="activeTab === 'logs'">
                <div
                  v-if="message.logs && message.logs.length > 0"
                  class="space-y-2"
                >
                  <div
                    v-for="(log, index) in message.logs"
                    :key="index"
                    class="bg-base-300 rounded-lg p-3 font-mono text-sm"
                  >
                    <span class="text-base-content/40 mr-2"
                      >[{{ index + 1 }}]</span
                    >
                    <span>{{ log }}</span>
                  </div>
                </div>
                <div v-else class="text-center text-base-content/40 py-8">
                  No logs available
                </div>
              </div>
            </div>

            <!-- Timestamps -->
            <div
              class="mt-4 pt-4 border-t border-base-300 grid grid-cols-2 md:grid-cols-4 gap-3 text-xs"
            >
              <div>
                <div class="text-base-content/60 mb-1">Received</div>
                <div class="font-mono">
                  {{ formatTimestamp(message.receivedAt) }}
                </div>
              </div>
              <div>
                <div class="text-base-content/60 mb-1">Processed</div>
                <div class="font-mono">
                  {{ formatTimestamp(message.processedAt) }}
                </div>
              </div>
              <div v-if="message.stateChangedAt">
                <div class="text-base-content/60 mb-1">State Changed</div>
                <div class="font-mono">
                  {{ formatTimestamp(message.stateChangedAt) }}
                </div>
              </div>
              <div>
                <div class="text-base-content/60 mb-1">Message ID</div>
                <div class="font-mono truncate" :title="message.id">
                  {{ message.id }}
                </div>
              </div>
            </div>

            <!-- Actions -->
            <div
              class="mt-4 pt-4 border-t border-base-300 flex gap-2 justify-end flex-wrap"
            >
              <button class="btn btn-sm btn-ghost" @click="copyId">
                Copy ID
              </button>
              <button
                v-if="message.state === 'new'"
                class="btn btn-sm btn-ghost"
                @click="markAsRead"
              >
                Mark as Read
              </button>

              <!-- State change dropdown -->
              <div class="dropdown dropdown-end">
                <button tabindex="0" class="btn btn-sm btn-primary">
                  Change State
                </button>
                <ul
                  tabindex="0"
                  class="dropdown-content menu p-2 shadow-lg bg-base-200 rounded-box w-48 mt-2 z-10"
                >
                  <li>
                    <a @click="changeState('new')">
                      <component
                        :is="stateIcons.new"
                        class="h-4 w-4"
                        :class="stateColors.new"
                      />
                      New
                    </a>
                  </li>
                  <li>
                    <a @click="changeState('acknowledged')">
                      <component
                        :is="stateIcons.acknowledged"
                        class="h-4 w-4"
                        :class="stateColors.acknowledged"
                      />
                      Acknowledged
                    </a>
                  </li>
                  <li>
                    <a @click="changeState('resolved')">
                      <component
                        :is="stateIcons.resolved"
                        class="h-4 w-4"
                        :class="stateColors.resolved"
                      />
                      Resolved
                    </a>
                  </li>
                  <li>
                    <a @click="changeState('archived')">
                      <component
                        :is="stateIcons.archived"
                        class="h-4 w-4"
                        :class="stateColors.archived"
                      />
                      Archived
                    </a>
                  </li>
                </ul>
              </div>

              <button
                class="btn btn-sm btn-error btn-outline"
                @click="handleDelete"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.expand {
  display: grid;
  transition: grid-template-rows 0.1s ease-in-out;
  grid-template-rows: 0fr;
}

.expand > div {
  overflow: hidden;
}

.expand-open {
  grid-template-rows: 1fr;
}
</style>
