import { router } from "@/router";
import dayjs from "dayjs";
import localizedFormat from "dayjs/plugin/localizedFormat";
import { createPinia } from "pinia";
import { Quasar } from "quasar";
import quasarIconSet from "quasar/icon-set/svg-mdi-v7";
import "quasar/dist/quasar.css";
import "@quasar/extras/roboto-font/roboto-font.css";
import { createApp } from "vue";
import App from "./App.vue";
import "@quasar/extras/mdi-v7/mdi-v7.css";

dayjs.extend(localizedFormat);

const app = createApp(App);

app.use(Quasar, {
  iconSet: quasarIconSet,
});
app.use(createPinia());
app.use(router);

app.mount("#app");
