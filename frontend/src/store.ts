import {
  GetAdGuardAccount,
  GetAdGuardBin,
  GetAdGuardStatus,
  GetAdGuardVersion,
  UpdateAdGuardBin,
} from "@/go/main/App";
import type { adguard } from "@/go/models";
import { defineStore } from "pinia";
import { computed, onBeforeMount, readonly, shallowRef } from "vue";

export const useAppStore = defineStore("app", () => {
  const isInitialized = shallowRef(false);
  const cliBin = shallowRef("");
  const cliVersion = shallowRef("");
  const account = shallowRef<adguard.Account | null>(null);
  const status = shallowRef<adguard.Status | null>(null);

  const isPremium = computed(
    () => account.value?.subscription.type === "PREMIUM",
  );

  onBeforeMount(async () => {
    try {
      await loadCliBin();
      await updateCliVersion();
      await updateStatus();
      await updateAccount();
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

  async function updateAccount() {
    account.value = await GetAdGuardAccount();
  }

  async function updateAdGuardBin(bin: string) {
    cliVersion.value = await UpdateAdGuardBin(bin);
    cliBin.value = bin;
  }

  return {
    isInitialized: readonly(isInitialized),
    adGuardBin: readonly(cliBin),
    cliVersion: readonly(cliVersion),
    account: readonly(account),
    status: readonly(status),

    isPremium,

    updateAdGuardBin,
    updateAccount,
  };
});
