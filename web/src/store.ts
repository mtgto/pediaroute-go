import { createStore } from 'vuex';

export interface State {
  wordFrom: string;
  wordTo: string;
}

export const store = createStore({
  //strict: process.env.NODE_ENV !== 'production',
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
