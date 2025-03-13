<script setup lang="ts">
import { useAppStore } from "@/store";
import { computed } from "vue";

const store = useAppStore();

const statusColor = computed(() => {
  if (!store.status) {
    return "amber";
  } else if (store.status.connected) {
    return "green";
  } else {
    return "grey";
  }
});

const statusMessage = computed(() =>
  store.status?.connected ? "Connected" : "Disconnected",
);
</script>

<template>
  <q-footer>
    <q-toolbar>
      <div>
        <q-badge :color="statusColor" rounded />
        {{ statusMessage }}
      </div>
    </q-toolbar>
  </q-footer>
</template>
