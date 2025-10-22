<script setup lang="ts">
// Initialize auth state and theme on app mount
const { initAuth } = useAuth()
const { initTheme } = useTheme()

onMounted(() => {
  initTheme()
  initAuth()
})
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
