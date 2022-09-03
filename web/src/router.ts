import { createRouter, createWebHistory, RouteLocationNormalized } from 'vue-router';
import { rootI18n } from './main';
import About from './views/About.vue';
import Home from './views/Home.vue';
import Search from './views/Search.vue';

declare module 'vue-router' {
  interface RouteMeta {
    title: string | ((route: RouteLocationNormalized) => string);
  }
}

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
      meta: {
        title: 'PediaRoute',
      },
    },
    {
      path: '/about',
      name: 'about',
      component: About,
      meta: {
        title: 'PediaRouteについて',
      },
    },
    {
      path: '/search',
      name: 'search',
      component: Search,
      // TODO
      props: (route: RouteLocationNormalized) => ({
        wordFrom: route.query.wordFrom as string,
        wordTo: route.query.wordTo as string,
      }),
      meta: {
        title: (route: RouteLocationNormalized) => `${route.query.wordFrom} から ${route.query.wordTo} の検索結果 - PediaRoute`,
      },
    },
  ],
});

router.beforeEach((to: RouteLocationNormalized, from: RouteLocationNormalized, next) => {
  if (to.query.lang === 'ja' || to.query.lang === 'en') {
    rootI18n.locale.value = to.query.lang;
  }
  // TODO
  if (typeof to.meta.title === 'string') {
    document.title = to.meta.title;
  } else {
    document.title = to.meta.title(to);
  }
  next();
});
