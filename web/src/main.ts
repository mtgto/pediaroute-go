import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import { i18n } from './i18n';
import { router } from './router';
import { useMainStore } from './store';

export const rootI18n = i18n.global;
createApp(App).use(i18n).use(router).use(createPinia()).mount('#app');
useMainStore();
