<script lang="ts" setup>
import { useAppStore } from "@/store";
import { useConfirmDialog } from "@vueuse/core";
import { shallowRef } from "vue";

const store = useAppStore();

const { item } = defineProps<{
  item: string;
}>();

const emit = defineEmits<{
  delete: [];
}>();

const deleting = shallowRef(false);
const error = shallowRef("");

const { isRevealed, reveal, confirm, onConfirm, cancel } = useConfirmDialog();

onConfirm(async () => {
  deleting.value = true;
  try {
    await store.deleteExclusion(item);
    emit("delete");
  } catch (e) {
    error.value = e as string;
  } finally {
    deleting.value = false;
  }
});
</script>
<template>
  <q-item>
    <q-item-section>
      <q-item-label>{{ item }}</q-item-label>
    </q-item-section>
    <q-item-section side>
      <q-btn flat round icon="mdi-delete" color="red" @click="reveal" />
    </q-item-section>
  </q-item>

  <q-dialog
    :model-value="isRevealed"
    @update:model-value="cancel"
    :persistent="deleting"
  >
    <q-card>
      <q-card-section class="row">
        <div class="text-h6">Delete website?</div>
      </q-card-section>
      <q-card-section>
        {{ item }}
      </q-card-section>
      <q-card-actions>
        <q-btn
          @click="confirm"
          label="Delete"
          color="red"
          :loading="deleting"
        />
        <q-btn @click="cancel" label="Cancel" :disable="deleting" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
