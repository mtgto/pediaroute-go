<template>
  <div>
    <header>
      <h1>
        <span>P</span>edia
        <span>R</span>oute.com
      </h1>
      <p>「{{wordFrom}}」から「{{wordTo}}」へのリンクの検索結果 (実行時間 {{time / 1000}} sec)</p>
    </header>
    <article v-if="route">
      <ol start="0">
        <li v-for="word in route">
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
            :to="{path: '/search', query: {wordFrom: wordTo, wordTo: wordFrom}}"
          >「{{wordTo}}」から「{{wordFrom}}」を検索する</router-link>
        </li>
        <li v-if="route">
          <a
            v-bind:href="`https://twitter.com/home?status=「${wordFrom}」から「${wordTo}」へはWikipediaで${route.length-1}リンクで行けるよ！ ${encodeURIComponent(`https://pediaroute.com/search?wordFrom=${encodeURIComponent(wordFrom)}&wordTo=${encodeURIComponent(wordTo)}`)} ${encodeURIComponent('#pediaroute')}`"
            target="_blank"
          >結果をTwitterにつぶやく (別ウィンドウで開きます)</a>
        </li>
        <li v-else>
          <a
            v-bind:href="`https://twitter.com/home?status=「${wordFrom}」から「${wordTo}」へはWikipediaで6回のリンクじゃいけないみたい… ${encodeURIComponent('#pediaroute')}`"
            target="_blank"
          >結果をTwitterにつぶやく (別ウィンドウで開きます)</a>
        </li>
      </ul>
    </aside>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';

interface Result {
    readonly route: ReadonlyArray<string> | undefined;
    readonly error: string;
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
            fetch('/api/search', { method: 'POST', body, headers })
                .then(response => response.json())
                .then((result: Result) => {
                    this.route = result.route;
                    this.failureReason = result.error;
                    this.time = new Date().getTime() - start;
                });
        },
    },
});
</script>
