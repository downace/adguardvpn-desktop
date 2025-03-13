<script setup lang="ts">
import { PickFilePath } from "@/go/main/App";
import { useAppStore } from "@/store";
import { shallowRef, watch } from "vue";

const store = useAppStore();

const adGuardBin = shallowRef("");

watch(
  () => store.adGuardBin,
  (bin) => {
    adGuardBin.value = bin;
  },
  { immediate: true },
);

async function pickCliPath() {
  const newPath = await PickFilePath();
  if (!newPath) {
    // No file selected
    return;
  }
  adGuardBin.value = newPath;
}

const saveResult = shallowRef<{
  success: boolean;
  message: string;
} | null>(null);

async function save() {
  saveResult.value = null;
  try {
    await store.updateAdGuardBin(adGuardBin.value);
    saveResult.value = {
      success: true,
      message: store.cliVersion,
    };
  } catch (e) {
    saveResult.value = {
      success: false,
      message: e as string,
    };
  }
}
</script>

<template>
  <q-card>
    <q-card-section>
      <q-input
        v-model.lazy="adGuardBin"
        label="AdGuard CLI binary"
        hint="Usually located in /opt/adguardvpn_cli"
        placeholder="adguardvpn-cli"
        class="q-mb-sm"
        bottom-slots
      >
        <template #append>
          <q-btn
            flat
            round
            icon="mdi-attachment"
            @click="pickCliPath"
            title="Pick path"
          />
        </template>
      </q-input>

      <q-btn icon="mdi-content-save" label="Save" @click="save" />
    </q-card-section>
    <q-banner
      v-if="saveResult"
      class="text-white full-width overflow-auto"
      :class="saveResult.success ? 'bg-green ' : 'bg-red'"
    >
      <pre>{{ saveResult.message }}</pre>
    </q-banner>
  </q-card>
</template>
