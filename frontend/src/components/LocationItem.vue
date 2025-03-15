<script setup lang="ts">
import { adguard } from "@/go/models";
import * as flags from "country-flag-icons/string/1x1";
import { computed } from "vue";

const { location } = defineProps<{
  location: adguard.Location;
}>();

const flag = computed(() =>
  location.iso ? (flags as Record<string, string>)[location.iso] : null,
);
</script>

<template>
  <q-item>
    <q-item-section avatar>
      <q-avatar class="shadow-1">
        <q-img v-if="flag" v-html="flag" />
        <q-icon v-else name="mdi-help"></q-icon>
      </q-avatar>
    </q-item-section>
    <q-item-section>
      <q-item-label>{{ location.city }}</q-item-label>
      <q-item-label caption>
        {{ location.country }}
      </q-item-label>
    </q-item-section>
    <q-item-section v-if="location.ping >= 0" side>
      {{ location.ping }}ms
    </q-item-section>
  </q-item>
</template>
