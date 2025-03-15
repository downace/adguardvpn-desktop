import {
  AdGuardListLocations,
  GetAdGuardAccount,
  GetAdGuardBin,
  GetAdGuardStatus,
  GetAdGuardVersion,
  UpdateAdGuardBin,
} from "@/go/main/App";
import type { adguard } from "@/go/models";
import { defineStore } from "pinia";
import { computed, onBeforeMount, readonly, shallowRef, watch } from "vue";

export const useAppStore = defineStore("app", () => {
  const isInitialized = shallowRef(false);
  const cliBin = shallowRef("");
  const cliVersion = shallowRef("");
  const account = shallowRef<adguard.Account | null>(null);
  const status = shallowRef<adguard.Status | null>(null);
  const locations = shallowRef<adguard.Location[]>([]);
  const locationsLoading = shallowRef(false);

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

  watch(account, (acc) => {
    if (acc) {
      void loadLocations();
    }
  });

  async function loadLocations() {
    locationsLoading.value = true;
    try {
      locations.value = await AdGuardListLocations();
    } finally {
      locationsLoading.value = false;
    }
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
    locations: readonly(locations),
    locationsLoading: readonly(locationsLoading),

    isPremium,

    updateAdGuardBin,
    updateAccount,
  };
});
