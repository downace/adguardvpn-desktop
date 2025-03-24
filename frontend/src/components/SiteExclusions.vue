<script lang="ts" setup>
import SiteExclusionsChangeMode from "@/components/SiteExclusionsChangeMode.vue";
import SiteExclusionsListItem from "@/components/SiteExclusionsListItem.vue";
import SiteExclusionsNew from "@/components/SiteExclusionsNew.vue";
import { adguard } from "@/go/models";
import { useAppStore } from "@/store";
import { shallowRef, watch } from "vue";

const store = useAppStore();

const exclusions = shallowRef([] as string[]);
const exclusionsLoading = shallowRef(false);
const error = shallowRef("");

async function loadExclusions() {
  exclusionsLoading.value = true;
  try {
    exclusions.value = await store.getExclusionsList();
  } catch (e) {
    error.value = e as string;
  } finally {
    exclusionsLoading.value = false;
  }
}

watch(
  () => store.exclusionMode,
  () => {
    loadExclusions();
  },
  { immediate: true },
);
</script>
<template>
  <q-page padding>
    <div class="row">
      <div class="col-6 column q-gutter-sm">
        <site-exclusions-change-mode #="{ onOpen }">
          <div v-if="store.exclusionMode === adguard.ExclusionMode.GENERAL">
            VPN is active
            <a @click="onOpen"> everywhere </a>
            except for the websites below
          </div>
          <div
            v-else-if="store.exclusionMode === adguard.ExclusionMode.SELECTIVE"
          >
            VPN is active <a @click="onOpen"> selectively </a>: only for the
            websites below
          </div>
        </site-exclusions-change-mode>

        <site-exclusions-new #="{ onOpen }" @save="loadExclusions">
          <q-btn
            label="Add website"
            icon="mdi-plus"
            color="green"
            @click="onOpen"
          />
        </site-exclusions-new>

        <q-card>
          <q-list dense>
            <site-exclusions-list-item
              v-for="item in exclusions"
              :item="item"
              @delete="loadExclusions"
            />
            <template v-if="exclusionsLoading">
              <q-linear-progress indeterminate />
              <q-item>
                <q-item-section class="text-center"> Loading </q-item-section>
              </q-item>
            </template>
          </q-list>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<style scoped>
a {
  text-decoration: underline;
  color: var(--q-primary);
  cursor: pointer;
}
</style>
