import { defineStore } from 'pinia';

export interface State {
  wordFrom: string;
  wordTo: string;
}

export const useMainStore = defineStore('main', {
  //strict: process.env.NODE_ENV !== 'production',
  state: () => ({
    wordFrom: '',
    wordTo: '',
  }),
  actions: {
    setWordFrom(wordFrom: string) {
      this.wordFrom = wordFrom;
    },
    setWordTo(wordTo: string) {
      this.wordTo = wordTo;
    },
  },
});
