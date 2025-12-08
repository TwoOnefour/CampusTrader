<script setup lang="ts">
import { NMessageProvider } from 'naive-ui'
import Dashboard from './components/Dashboard.vue'
</script>

<template>
  <n-message-provider>
    <Dashboard />
  </n-message-provider>
</template>