import Vue from 'vue';
import Vuex from 'vuex';
import Router, { Route } from 'vue-router';
import Home from './views/Home.vue';
import Search from './views/Search.vue';

Vue.use(Router);

const router = new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
      meta: {
        title: (route: Route) => 'PediaRoute',
      },
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import(/* webpackChunkName: "about" */ './views/About.vue'),
      meta: {
        title: (route: Route) => 'PediaRouteについて',
      },
    },
    {
      path: '/search',
      name: 'search',
      component: Search,
      props: (route: Route): { [key: string]: string | string[] } => ({
        wordFrom: route.query.wordFrom,
        wordTo: route.query.wordTo,
      }),
      meta: {
        title: (route: Route) => `${route.query.wordFrom} から ${route.query.wordTo} の検索結果 - PediaRoute`,
      },
    },
  ],
});

router.beforeEach((to: Route, from: Route, next) => {
  document.title = to.meta.title(to);
  next();
});

export default router;
