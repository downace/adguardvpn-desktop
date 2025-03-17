<script lang="ts" setup>
import LocationItem from "@/components/LocationItem.vue";
import { useAppStore } from "@/store";

const store = useAppStore();
</script>
<template>
  <div class="column no-wrap justify-end items-center">
    <div class="col flex flex-center">
      <q-img src="@/assets/images/logo.svg" width="200px" />
    </div>
    <h4>
      {{ store.status?.connected ? "Connected" : "Disconnected" }}
    </h4>
    <q-btn
      :label="store.status?.connected ? 'Disconnect' : 'Connect'"
      :loading="store.connecting"
      @click="store.toggleConnection()"
    />
    <!-- 54px is banner height, prevents size changes -->
    <q-banner
      rounded
      class="self-stretch bg-grey-3 q-ma-lg q-pa-none"
      style="height: 54px"
    >
      <location-item
        v-if="store.status?.location"
        :location="store.status.location"
      />
    </q-banner>
  </div>
</template>
