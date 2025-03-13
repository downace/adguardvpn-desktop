import {
  GetAdGuardBin,
  GetAdGuardStatus,
  GetAdGuardVersion,
  UpdateAdGuardBin,
} from "@/go/main/App";
import type { adguard } from "@/go/models";
import { defineStore } from "pinia";
import { onBeforeMount, readonly, shallowRef } from "vue";

export const useAppStore = defineStore("app", () => {
  const isInitialized = shallowRef(false);
  const cliBin = shallowRef("");
  const cliVersion = shallowRef("");
  const status = shallowRef<adguard.Status | null>(null);

  onBeforeMount(async () => {
    try {
      await loadCliBin();
      await updateCliVersion();
      await updateStatus();
    } finally {
      isInitialized.value = true;
    }
  });

  async function loadCliBin() {
    cliBin.value = await GetAdGuardBin();
  }

  async function updateCliVersion() {
    cliVersion.value = await GetAdGuardVersion();
  }

  async function updateStatus() {
    status.value = await GetAdGuardStatus();
  }

  async function updateAdGuardBin(bin: string) {
    cliVersion.value = await UpdateAdGuardBin(bin);
    cliBin.value = bin;
  }

  return {
    isInitialized: readonly(isInitialized),
    adGuardBin: readonly(cliBin),
    cliVersion: readonly(cliVersion),
    status: readonly(status),

    updateAdGuardBin,
  };
});
