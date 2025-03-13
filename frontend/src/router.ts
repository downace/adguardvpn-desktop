import AppDashboard from "@/components/AppDashboard.vue";
import AppSettings from "@/components/AppSettings.vue";
import { createRouter, createWebHashHistory } from "vue-router";

export const router = createRouter({
  // do not use `createWebHistory`: https://github.com/wailsapp/wails/issues/2262
  history: createWebHashHistory(),
  routes: [
    {
      path: "/",
      component: AppDashboard,
    },
    {
      path: "/settings",
      component: AppSettings,
    },
  ],
});
