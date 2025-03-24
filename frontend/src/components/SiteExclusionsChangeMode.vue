<script setup lang="ts">
import { adguard } from "@/go/models";
import { useAppStore } from "@/store";
import { shallowRef, watch } from "vue";

const store = useAppStore();

const isOpen = shallowRef(false);

const currentMode = shallowRef(adguard.ExclusionMode.GENERAL);
const error = shallowRef("");

watch(
  [() => store.exclusionMode, isOpen],
  ([mode]) => {
    currentMode.value = mode;
  },
  { immediate: true },
);

const saving = shallowRef(false);

async function save() {
  saving.value = true;

  try {
    await store.setExclusionMode(currentMode.value);
    isOpen.value = false;
  } catch (e) {
    error.value = e as string;
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <slot @open="isOpen = true" />
  <q-dialog v-model="isOpen">
    <q-card>
      <q-card-section class="row">
        <div class="text-h6">Where VPN is active</div>
      </q-card-section>
      <q-card-section>
        <div>
          <q-radio
            v-model="currentMode"
            :val="adguard.ExclusionMode.GENERAL"
            label="Everywhere except for excluded websites"
          />
        </div>
        <div>
          <q-radio
            v-model="currentMode"
            :val="adguard.ExclusionMode.SELECTIVE"
            label="Only on selected websites"
          />
        </div>
      </q-card-section>
      <q-card-actions>
        <q-btn label="Save" color="primary" :loading="saving" @click="save" />
        <q-btn label="Cancel" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
