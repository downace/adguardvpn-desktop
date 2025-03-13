<script lang="ts" setup>
import AppDashboardStatus from "@/components/AppDashboardStatus.vue";
import LocationsList from "@/components/LocationsList.vue";
import { fullHeightPageStyleFn } from "@/helpers/fullHeightPageStyleFn";
</script>
<template>
  <q-page :style-fn="fullHeightPageStyleFn">
    <div class="row no-wrap full-height">
      <app-dashboard-status class="col-6" />
      <locations-list class="col-6" />
    </div>
  </q-page>
</template>
