import AppDashboard from "@/components/AppDashboard.vue";
import SiteExclusions from "@/components/SiteExclusions.vue";
import AppSettings from "@/components/AppSettings.vue";
import AppUnauthorized from "@/components/AppUnauthorized.vue";
import { useAppStore } from "@/store";
import { until } from "@vueuse/core";
import {
  createRouter,
  createWebHashHistory,
  NavigationGuard,
} from "vue-router";

const authenticatedGuard: NavigationGuard = function (to, from, next) {
  const store = useAppStore();

  if (store.account) {
    next();
  } else {
    next("/unauthorized");
  }
};

const guestGuard: NavigationGuard = async function (to, from, next) {
  const store = useAppStore();

  await until(() => store.isInitialized).toBeTruthy();

  if (store.account) {
    next("/");
  } else {
    next();
  }
};

export const router = createRouter({
  // do not use `createWebHistory`: https://github.com/wailsapp/wails/issues/2262
  history: createWebHashHistory(),
  routes: [
    {
      path: "/",
      component: AppDashboard,
      beforeEnter: authenticatedGuard,
    },
    {
      path: "/unauthorized",
      component: AppUnauthorized,
      beforeEnter: guestGuard,
    },
    {
      path: "/exclusions",
      component: SiteExclusions,
      beforeEnter: authenticatedGuard,
    },
    {
      path: "/settings",
      beforeEnter: authenticatedGuard,
      component: AppSettings,
    },
  ],
});
