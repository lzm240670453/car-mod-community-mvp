import "element-plus/dist/index.css";
import "./styles.css";

import { createApp, type Plugin } from "vue";
import { createPinia } from "pinia";
import {
  ElAlert,
  ElAside,
  ElButton,
  ElContainer,
  ElDialog,
  ElForm,
  ElFormItem,
  ElHeader,
  ElInput,
  ElInputNumber,
  ElLoading,
  ElMain,
  ElOption,
  ElPagination,
  ElSelect,
  ElSwitch,
  ElTabPane,
  ElTable,
  ElTableColumn,
  ElTabs,
  ElTag
} from "element-plus";
import * as ElementPlusIconsVue from "@element-plus/icons-vue";
import App from "./App.vue";
import { router } from "./router";

const app = createApp(App);
const elementPlusPlugins: Plugin[] = [
  ElAlert,
  ElAside,
  ElButton,
  ElContainer,
  ElDialog,
  ElForm,
  ElFormItem,
  ElHeader,
  ElInput,
  ElInputNumber,
  ElLoading,
  ElMain,
  ElOption,
  ElPagination,
  ElSelect,
  ElSwitch,
  ElTabPane,
  ElTable,
  ElTableColumn,
  ElTabs,
  ElTag
];

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component);
}

for (const plugin of elementPlusPlugins) {
  app.use(plugin);
}

app.use(createPinia());
app.use(router);
app.mount("#app");
