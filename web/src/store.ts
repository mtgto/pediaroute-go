import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

export const store = new Vuex.Store({
  strict: process.env.NODE_ENV !== 'production',
  state: {
    wordFrom: '',
    wordTo: '',
  },
  mutations: {
    setWordFrom(state, wordFrom: string) {
      state.wordFrom = wordFrom;
    },
    setWordTo(state, wordTo: string) {
      state.wordTo = wordTo;
    },
  },
  actions: {},
});
