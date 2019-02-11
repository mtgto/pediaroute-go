<template>
  <div>
    <header>
      <h1>
        <span>P</span>edia
        <span>R</span>oute.com
      </h1>
      <p v-if="route">{{ $t('message.searchResult', { wordFrom, wordTo, second: time / 1000 }) }}</p>
      <p v-else>{{ $t('message.searching', { wordFrom, wordTo }) }}</p>
    </header>
    <article v-if="route">
      <ol start="0">
        <li v-for="word in route" v-bind:key="word">
          <a v-bind:href="`http://ja.wikipedia.org/wiki/${word}`">{{word}}</a>
        </li>
      </ol>
    </article>
    <article class="error" v-if="failureReason">
      <p>{{failureReason}}</p>
    </article>
    <aside id="links">
      <ul>
        <li>
          <router-link
            :to="{path: '/search', query: {lang: this.$i18n.locale, wordFrom: wordTo, wordTo: wordFrom}}"
          >{{ $t('message.searchInReverse', { wordFrom, wordTo })}}</router-link>
        </li>
        <li v-if="route">
          <a
            v-bind:href="`https://twitter.com/home?status=「${wordFrom}」から「${wordTo}」へはWikipediaで${route.length-1}リンクで行けるよ！ ${encodeURIComponent(`https://pediaroute.com/search?lang=${encodeURI(this.$i18n.locale)}&wordFrom=${encodeURIComponent(wordFrom)}&wordTo=${encodeURIComponent(wordTo)}`)} ${encodeURIComponent('#pediaroute')}`"
            target="_blank"
          >{{ $t('message.tweet') }}</a>
        </li>
        <li v-else>
          <a
            v-bind:href="`https://twitter.com/home?status=「${wordFrom}」から「${wordTo}」へはWikipediaで6回のリンクじゃいけないみたい… ${encodeURIComponent('#pediaroute')}`"
            target="_blank"
          >{{ $t('message.tweet') }}</a>
        </li>
      </ul>
    </aside>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';

enum ErrorCode {
  NoError = 0,
  NotFoundFrom = 1,
  NotFoundTo = 2,
  NotFoundRoute = 3,
}

interface Result {
  readonly route: ReadonlyArray<string> | undefined;
  readonly error: ErrorCode;
}

// Data structure
interface Data {
  route: ReadonlyArray<string> | undefined;
  failureReason: string | undefined;
  time: number;
}

export default Vue.extend({
  data(): Data {
    return {
      route: undefined, // Array of string
      failureReason: undefined, // string of faulure reason
      time: 0, // milliseconds of searching
    };
  },
  props: ['wordFrom', 'wordTo'],
  created() {
    this.search();
  },
  watch: {
    $route: 'search',
  },
  methods: {
    search() {
      const body = JSON.stringify({ from: this.wordFrom, to: this.wordTo });
      const headers = {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      };
      const start = new Date().getTime();
      this.route = undefined;
      fetch(`/api/search?lang=${encodeURI(this.$i18n.locale)}`, { method: 'POST', body, headers })
        .then(response => response.json())
        .then((result: Result) => {
          this.route = result.route;
          if (result.error === ErrorCode.NoError) {
            this.failureReason = undefined;
          } else if (result.error === ErrorCode.NotFoundFrom) {
            const message = this.$i18n.t('error.notFoundFrom', this.$i18n.locale, { from: this.wordFrom });
            if (typeof message === 'string') {
              this.failureReason = message;
            }
          } else if (result.error === ErrorCode.NotFoundTo) {
            const message = this.$i18n.t('error.notFoundTo', this.$i18n.locale, { to: this.wordTo });
            if (typeof message === 'string') {
              this.failureReason = message;
            }
          } else if (result.error === ErrorCode.NotFoundRoute) {
            const message = this.$i18n.t('error.notFoundRoute', this.$i18n.locale);
            if (typeof message === 'string') {
              this.failureReason = message;
            }
          }
          this.time = new Date().getTime() - start;
        });
    },
  },
});
</script>
