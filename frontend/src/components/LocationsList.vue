<script setup lang="ts">
import LocationItem from "@/components/LocationItem.vue";
import { adguard } from "@/go/models";
import { useAppStore } from "@/store";
import { MaybeRefOrGetter } from "@vueuse/core";
import { useFuse } from "@vueuse/integrations/useFuse";
import { computed, shallowRef } from "vue";

const store = useAppStore();

const searchString = shallowRef("");

const { results: allLocations } = useFuse(
  searchString,
  (() => store.locations) as MaybeRefOrGetter<adguard.Location[]>,
  {
    matchAllWhenSearchEmpty: true,
    fuseOptions: {
      keys: ["city", "country"] as (keyof adguard.Location)[],
    },
  },
);

const favoriteLocations = computed(() =>
  allLocations.value.filter(({ item }) => store.isFavorite(item)),
);

const tab = shallowRef("all");
</script>

<template>
  <div class="column no-wrap">
    <q-input v-model="searchString" label="Search" dense>
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
    <template v-else>
      <q-tabs v-model="tab" dense>
        <q-tab name="all" label="All" />
        <q-tab name="favorite" label="Favorite" />
      </q-tabs>
      <q-tab-panels v-model="tab" class="full-height">
        <q-tab-panel name="all">
          <q-list>
            <location-item
              v-for="res in allLocations"
              :location="res.item"
              show-favorite-icon
              clickable
              @click="store.connect(res.item)"
            />
          </q-list>
        </q-tab-panel>
        <q-tab-panel name="favorite">
          <q-list v-if="favoriteLocations.length > 0">
            <location-item
              v-for="res in favoriteLocations"
              :location="res.item"
              show-favorite-icon
              clickable
              @click="store.connect(res.item)"
            />
          </q-list>
          <div v-else class="full-height flex flex-center text-h6">
            <p class="text-center">
              Use <q-icon name="mdi-bookmark-outline" /> icon
              <br />
              to make location favorite
            </p>
          </div>
        </q-tab-panel>
      </q-tab-panels>
    </template>
  </div>
</template>
