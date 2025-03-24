<script setup lang="ts">
import { useAppStore } from "@/store";
import { QInput } from "quasar";
import { computed, shallowRef, watch } from "vue";

const emit = defineEmits<{
  save: [];
}>();

const isOpen = shallowRef(false);

const exclusions = shallowRef("");

watch(isOpen, () => {
  exclusions.value = "";
});

const saving = shallowRef(false);
const error = shallowRef("");

const store = useAppStore();

async function saveExclusions() {
  saving.value = true;
  try {
    await store.addExclusions(exclusions.value.split(/\s+/).filter((v) => !!v));
    isOpen.value = false;
    emit("save");
  } catch (e) {
    error.value = e as string;
  } finally {
    saving.value = false;
  }
}

const invalid = computed(() => !exclusions.value);
</script>

<template>
  <slot @open="isOpen = true" />

  <q-dialog v-model="isOpen" :persistent="saving">
    <q-card style="min-width: 400px">
      <q-card-section class="row">
        <div class="text-h6">Add websites</div>
      </q-card-section>
      <q-card-section>
        <q-input
          v-model="exclusions"
          autofocus
          type="textarea"
          :placeholder="'example.com\n192.168.0.3'"
          label="Website or IP address"
          bottom-slots
          hint="Separate multiple websites using spaces or newlines"
        />
      </q-card-section>
      <q-card-actions>
        <q-btn
          label="Add"
          color="primary"
          :loading="saving"
          :disable="invalid"
          @click="saveExclusions"
        />
        <q-btn label="Cancel" v-close-popup :disable="saving" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
