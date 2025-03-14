<script lang="ts" setup>
import AppFooter from "@/components/AppFooter.vue";
import { fullHeightPageStyleFn } from "@/helpers/fullHeightPageStyleFn";
import { useAppStore } from "@/store";

const store = useAppStore();
</script>

<template>
  <q-layout view="hhh lpR fff">
    <template v-if="store.isInitialized">
      <q-header v-if="store.account">
        <q-toolbar>
          <q-tabs align="left" no-caps shrink stretch>
            <q-route-tab to="/" label="Home" />
            <q-route-tab to="/settings" label="Settings" />
          </q-tabs>

          <q-space />

          <div>
            {{ store.cliVersion }}
          </div>
        </q-toolbar>
      </q-header>

      <q-page-container>
        <router-view />
      </q-page-container>

      <app-footer v-if="store.account" />
    </template>
    <q-page-container v-else>
      <q-page :style-fn="fullHeightPageStyleFn">
        <div class="full-height column justify-center items-center text-h4">
          <q-circular-progress indeterminate></q-circular-progress>
          Loading app
        </div>
      </q-page>
    </q-page-container>
  </q-layout>
</template>
