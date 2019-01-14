<template>
  <div>
    <header>
      <h1>
        <span>P</span>edia
        <span>R</span>oute.com
      </h1>
      <p>{{ $t('message.header') }}</p>
    </header>
    <article>
      <fieldset>
        <p>
          <label>
            <input type="text" class="word" name="wordFrom" v-model="wordFrom">
            <button
              class="random"
              v-on:click="getRandomFrom"
              v-bind:title="$t('message.buttonRandom')"
            >
              <img src="../assets/baseline-shuffle-24px.svg">
            </button>
            {{ $t('message.searchFrom') }}
          </label>
        </p>
        <p>
          <label>
            <input type="text" class="word" name="wordTo" v-model="wordTo">
            <button
              class="random"
              v-on:click="getRandomTo"
              v-bind:title="$t('message.buttonRandom')"
            >
              <img src="../assets/baseline-shuffle-24px.svg">
            </button>
            {{ $t('message.searchTo') }}
          </label>
        </p>
      </fieldset>
      <div class="center">
        <input type="button" class="submit" v-bind:value="$t('message.search')" v-on:click="search">
      </div>
      <div class="center">
        <p>
          <a v-if="this.$i18n.locale === 'ja'" href="#" v-on:click="chooseEnglish">Choose English</a>
          <a v-if="this.$i18n.locale === 'en'" href="#" v-on:click="chooseJapanese">日本語を選択</a>
        </p>
      </div>
    </article>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import { store } from '../store';

export default Vue.extend({
  computed: {
    wordFrom() {
      return store.state.wordFrom;
    },
    wordTo() {
      return store.state.wordTo;
    },
  },
  methods: {
    chooseEnglish() {
      this.$i18n.locale = 'en';
    },
    chooseJapanese() {
      this.$i18n.locale = 'ja';
    },
    getRandomFrom() {
      fetch(`/api/random?lang=${encodeURI(this.$i18n.locale)}`)
        .then(response => response.json())
        .then(word => {
          if (typeof word === 'string') {
            store.commit('setWordFrom', word);
          }
        })
        .catch(error => console.log(error));
    },
    getRandomTo() {
      fetch(`/api/random?lang=${encodeURI(this.$i18n.locale)}`)
        .then(response => response.json())
        .then(word => {
          if (typeof word === 'string') {
            store.commit('setWordTo', word);
          }
        })
        .catch(error => console.log(error));
    },
    search() {
      this.$router.push({ path: '/search', query: { wordFrom: this.wordFrom, wordTo: this.wordTo } });
    },
  },
});
</script>
