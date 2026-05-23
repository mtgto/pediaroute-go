import { createApp } from 'vue';
import { createPinia } from 'pinia';

import './styles/tokens.css';
import './styles/base.css';

import App from './App.vue';
import { i18n } from './i18n';
import { router } from './router';

createApp(App).use(i18n).use(router).use(createPinia()).mount('#app');
