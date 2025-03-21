import {
  AddFavoriteLocation,
  AdGuardConnect,
  AdGuardDisconnect,
  AdGuardGetLocations,
  GetAdGuardAccount,
  GetAdGuardBin,
  GetAdGuardStatus,
  GetAdGuardVersion,
  GetFavoriteLocations,
  RemoveFavoriteLocation,
  UpdateAdGuardBin,
} from "@/go/main/App";
import type { adguard } from "@/go/models";
import { EventsOn } from "@/runtime";
import { defineStore } from "pinia";
import {
  computed,
  onBeforeMount,
  onBeforeUnmount,
  readonly,
  shallowReactive,
  shallowRef,
  watch,
} from "vue";

export const useAppStore = defineStore("app", () => {
  const isInitialized = shallowRef(false);
  const cliBin = shallowRef("");
  const cliVersion = shallowRef("");
  const account = shallowRef<adguard.Account | null>(null);
  const status = shallowRef<adguard.Status | null>(null);
  const connecting = shallowRef(false);
  const locations = shallowRef<adguard.Location[]>([]);
  const favoriteLocations = shallowReactive(new Set<string>());
  const locationsLoading = shallowRef(false);

  const isPremium = computed(
    () => account.value?.subscription.type === "PREMIUM",
  );

  onBeforeMount(async () => {
    try {
      await loadCliBin();
      await updateCliVersion();
      await updateAccount();
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

  watch(status, (status) => {
    if (status) {
      // not using computed to allow instant update when user connects from app UI
      connecting.value = status.connecting;
    }
  });

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
      const [allLocations, favorites] = await Promise.all([
        AdGuardGetLocations(),
        GetFavoriteLocations(),
      ]);
      locations.value = allLocations;
      for (const loc of favorites) {
        favoriteLocations.add(loc);
      }
    } finally {
      locationsLoading.value = false;
    }
  }

  function isFavorite(location: adguard.Location) {
    return favoriteLocations.has(location.city);
  }

  async function addToFavorites(location: adguard.Location) {
    await AddFavoriteLocation(location.city);
    favoriteLocations.add(location.city);
  }

  async function removeFromFavorites(location: adguard.Location) {
    await RemoveFavoriteLocation(location.city);
    favoriteLocations.delete(location.city);
  }

  async function updateAdGuardBin(bin: string) {
    cliVersion.value = await UpdateAdGuardBin(bin);
    cliBin.value = bin;
  }

  async function toggleConnection() {
    if (status.value?.connecting) {
      return;
    }
    if (status.value?.connected) {
      return disconnect();
    } else {
      return connect();
    }
  }

  async function connect(location?: adguard.Location) {
    if (connecting.value) {
      return;
    }
    connecting.value = true;

    try {
      await AdGuardConnect(location?.city ?? "");
    } finally {
      connecting.value = false;
    }
  }

  async function disconnect() {
    if (connecting.value) {
      return;
    }

    connecting.value = true;

    try {
      await AdGuardDisconnect();
    } finally {
      connecting.value = false;
    }
  }

  const unsubscribeFns = [
    EventsOn("status-changed", (s) => {
      status.value = s;
    }),
  ];

  onBeforeUnmount(() => {
    for (const un of unsubscribeFns) {
      un();
    }
  });

  return {
    isInitialized: readonly(isInitialized),
    adGuardBin: readonly(cliBin),
    cliVersion: readonly(cliVersion),
    account: readonly(account),
    status: readonly(status),
    connecting: readonly(connecting),
    locations: readonly(locations),
    locationsLoading: readonly(locationsLoading),

    isPremium,

    updateAdGuardBin,
    updateAccount,
    connect,
    toggleConnection,
    isFavorite,
    addToFavorites,
    removeFromFavorites,
  };
});
