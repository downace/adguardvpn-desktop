<script setup lang="ts">
import { adguard } from "@/go/models";
import { useAppStore } from "@/store";
import * as flags from "country-flag-icons/string/1x1";
import { computed, shallowRef } from "vue";

const store = useAppStore();

const { location, showFavoriteIcon } = defineProps<{
  location: adguard.Location;
  showFavoriteIcon?: boolean;
}>();

const flag = computed(() =>
  location.iso ? (flags as Record<string, string>)[location.iso] : null,
);

const isFavorite = computed(() => location && store.isFavorite(location));

const togglingFavorite = shallowRef(false);

async function toggleFavorite() {
  togglingFavorite.value = true;
  try {
    if (isFavorite.value) {
      await store.removeFromFavorites(location!);
    } else {
      await store.addToFavorites(location!);
    }
  } finally {
    togglingFavorite.value = false;
  }
}
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
    <q-item-section v-if="location && showFavoriteIcon" side>
      <q-btn
        flat
        round
        :icon="isFavorite ? 'mdi-bookmark' : 'mdi-bookmark-outline'"
        :title="isFavorite ? 'Remove from favorites' : 'Add to favorites'"
        :loading="togglingFavorite"
        @click.stop="toggleFavorite"
      ></q-btn>
    </q-item-section>
  </q-item>
</template>
