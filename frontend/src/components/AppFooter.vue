<script setup lang="ts">
import { BrowserOpenURL } from "@/runtime";
import { useAppStore } from "@/store";
import dayjs from "dayjs";

const store = useAppStore();
</script>

<template>
  <q-footer>
    <q-toolbar>
      <template v-if="store.account">
        <q-chip
          clickable
          @click="BrowserOpenURL('https://auth.adguardaccount.net/')"
        >
          {{ store.account.username }}
          <q-tooltip class="text-body2">
            Click to open account settings
          </q-tooltip>
        </q-chip>

        <q-space />

        <q-chip
          :icon="store.isPremium ? 'mdi-check-decagram' : undefined"
          :color="store.isPremium ? 'green' : undefined"
          :text-color="store.isPremium ? 'white' : undefined"
        >
          {{ store.account.subscription.type }}
          <q-tooltip v-if="store.isPremium" class="text-body2">
            Subscription valid until
            {{ dayjs(store.account.subscription.validUntil).format("LL") }}
          </q-tooltip>
        </q-chip>
      </template>
    </q-toolbar>
  </q-footer>
</template>
