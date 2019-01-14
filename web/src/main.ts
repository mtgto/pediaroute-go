import Vue from 'vue';
import App from './App.vue';
import router from './router';
import { i18n } from './i18n';
import { store } from './store';

Vue.config.productionTip = false;

new Vue({
  i18n,
  router,
  store,
  render: h => h(App),
}).$mount('#app');
