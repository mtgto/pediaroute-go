<template>
  <div>
    <header>
      <h1>
        <span>P</span>edia
        <span>R</span>oute.com
      </h1>
      <p>Wikipediaで6回リンク辿ればいけるか調べる(仮)</p>
    </header>
    <article>
      <form action="/search" method="get">
        <fieldset>
          <p>
            <label>
              <input type="text" class="word" name="wordFrom" v-model="wordFrom" />から
            </label>
          </p>
          <p>
            <label>
              <input type="text" class="word" name="wordTo" v-model="wordTo" />へのルートを
            </label>
          </p>
        </fieldset>
        <input type="button" class="submit" value="検索" v-on:click="search" />
        <input type="button" id="getRandom" value="ランダムに２つのページを選択" v-on:click="getRandom" />
      </form>
    </article>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';

export default Vue.extend({
  data() {
    return {
      wordFrom: '',
      wordTo: '',
    };
  },
  methods: {
    getRandom() {
      fetch('/api/random')
        .then((response) => response.json())
        .then((pair) => {
          if (pair.hasOwnProperty('from') && pair.hasOwnProperty('to')) {
            this.wordFrom = pair.from;
            this.wordTo = pair.to;
          }
        })
        .catch((error) => console.log(error));
    },
    search() {
      this.$router.push({path: '/search', query: {wordFrom: this.wordFrom, wordTo: this.wordTo}});
    },
  },
});
</script>
