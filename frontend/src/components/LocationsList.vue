<script setup lang="ts">
import LocationItem from "@/components/LocationItem.vue";
import { adguard } from "@/go/models";
import { useAppStore } from "@/store";
import { MaybeRefOrGetter } from "@vueuse/core";
import { useFuse } from "@vueuse/integrations/useFuse";
import { shallowRef } from "vue";

const store = useAppStore();

const searchString = shallowRef("");

const { results } = useFuse(
  searchString,
  (() => store.locations) as MaybeRefOrGetter<adguard.Location[]>,
  {
    matchAllWhenSearchEmpty: true,
    fuseOptions: {
      keys: ["city", "country"] as (keyof adguard.Location)[],
    },
  },
);
</script>

<template>
  <div class="column no-wrap">
    <q-input v-model="searchString" label="Search" dense class="q-ma-md">
      <template #prepend>
        <q-icon name="mdi-magnify"></q-icon>
      </template>
    </q-input>
    <div
      v-if="store.locationsLoading"
      class="col column no-wrap justify-center items-center"
    >
      <q-circular-progress indeterminate size="xl" />
      Loading Locations...
    </div>
    <q-list v-else class="scroll-y">
      <location-item
        v-for="res in results"
        :location="res.item"
        clickable
        @click="store.connect(res.item)"
      />
    </q-list>
  </div>
</template>
