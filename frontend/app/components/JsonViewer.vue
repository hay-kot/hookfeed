<script setup lang="ts">
/**
 * Component for displaying JSON data with syntax highlighting
 */

interface Props {
  data: any
  compact?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  compact: false,
})

// Format JSON with indentation
const formattedJson = computed(() => {
  if (!props.data) return 'null'
  try {
    return JSON.stringify(props.data, null, props.compact ? 0 : 2)
  } catch {
    return String(props.data)
  }
})

// Simple syntax highlighting using HTML
const highlightedJson = computed(() => {
  let json = formattedJson.value

  // Escape HTML
  json = json
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')

  // Highlight strings (including keys)
  json = json.replace(/"([^"]+)":/g, '<span class="json-key">"$1"</span>:')
  json = json.replace(/: "([^"]*)"/g, ': <span class="json-string">"$1"</span>')

  // Highlight numbers
  json = json.replace(/: (-?\d+\.?\d*)/g, ': <span class="json-number">$1</span>')

  // Highlight booleans
  json = json.replace(/: (true|false)/g, ': <span class="json-boolean">$1</span>')

  // Highlight null
  json = json.replace(/: (null)/g, ': <span class="json-null">$1</span>')

  return json
})

// Copy to clipboard
const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(formattedJson.value)
    // Could add a toast notification here
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}
</script>

<template>
  <div class="relative group">
    <!-- Copy button -->
    <button
      class="btn btn-xs btn-ghost absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity z-10"
      @click="copyToClipboard"
    >
      Copy
    </button>

    <!-- JSON content -->
    <pre
      class="bg-base-300 rounded-lg p-4 overflow-x-auto text-sm font-mono leading-relaxed"
    ><code v-html="highlightedJson" /></pre>
  </div>
</template>

<style scoped>
/* JSON Syntax Highlighting with better contrast */
:deep(.json-key) {
  color: #0066cc;
  font-weight: 600;
}

:deep(.json-string) {
  color: #22863a;
}

:deep(.json-number) {
  color: #005cc5;
}

:deep(.json-boolean) {
  color: #d73a49;
}

:deep(.json-null) {
  color: #6f42c1;
  opacity: 0.7;
}

/* Dark theme adjustments */
[data-theme="dark"] :deep(.json-key) {
  color: #79c0ff;
}

[data-theme="dark"] :deep(.json-string) {
  color: #a5d6ff;
}

[data-theme="dark"] :deep(.json-number) {
  color: #79c0ff;
}

[data-theme="dark"] :deep(.json-boolean) {
  color: #ff7b72;
}

[data-theme="dark"] :deep(.json-null) {
  color: #d2a8ff;
}
</style>
