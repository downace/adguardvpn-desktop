<script setup lang="ts">
import AppSettingsAdguardCli from "@/components/AppSettings/AppSettingsAdguardCli.vue";
import { useAppStore } from "@/store";
import { whenever } from "@vueuse/core";
import { copyToClipboard } from "quasar";
import { shallowRef } from "vue";
import { useRouter } from "vue-router";

const store = useAppStore();

const error = shallowRef("");

const router = useRouter();

const updating = shallowRef(false);

async function updateAccount() {
  updating.value = true;
  error.value = "";
  try {
    await store.updateAccount();
    if (store.account) {
      await router.push("/");
    }
  } catch (e) {
    error.value = e as string;
  } finally {
    updating.value = false;
  }
}

const showSettings = shallowRef(false);
const cliWarning = shallowRef(false);

whenever(
  () => store.isInitialized,
  () => {
    if (!store.cliVersion) {
      cliWarning.value = true;
    }
  },
  { immediate: true },
);

function handleCliUpdateSuccess() {
  cliWarning.value = false;
  showSettings.value = false;
  void updateAccount();
}
</script>

<template>
  <q-page class="column justify-center items-center q-gutter-sm">
    <div>
      <q-img src="@/assets/images/logo.svg" width="200px" />
    </div>
    <div class="text-h6">You are not logged in</div>
    <div class="text-h6">You can log in by running</div>
    <div class="text-h6">
      <q-chip
        square
        dense
        outline
        clickable
        icon-right="mdi-content-copy"
        @click="copyToClipboard('adguardvpn-cli login')"
      >
        <code>adguardvpn-cli login</code>
      </q-chip>
    </div>
    <div class="column no-wrap q-gutter-sm items-center">
      <q-btn
        label="Refresh"
        icon="mdi-refresh"
        :loading="updating"
        @click="updateAccount"
      />
      <q-btn label="Settings" icon="mdi-cog" @click="showSettings = true" />
      <q-banner class="bg-warning" v-if="cliWarning">
        Could not determine AdGuard CLI version. Open settings and check the CLI
        path
      </q-banner>
      <q-dialog v-model="showSettings">
        <app-settings-adguard-cli
          style="min-width: 400px"
          @ok="handleCliUpdateSuccess"
        />
      </q-dialog>
    </div>
    <q-card v-if="error" class="bg-red text-white">
      <q-card-section>
        {{ error }}
      </q-card-section>
    </q-card>
  </q-page>
</template>
