<template>
  <div>
    <header>
      <h1><span>P</span>edia <span>R</span>oute.com</h1>
      <p v-if="routes" v-t="{ path: 'message.searchResult', args: { wordFrom, wordTo, second: time / 1000 } }" />
      <p v-else v-t="{ path: 'message.searching', args: { wordFrom, wordTo } }" />
    </header>
    <article v-if="routes">
      <ol start="0">
        <li v-for="(word, index) in routes" :key="`${index}`">
          <a :href="t('message.wikipediaUrl', { word })" v-text="word" />
        </li>
      </ol>
    </article>
    <article v-if="failureReason" class="error">
      <p v-text="failureReason" />
    </article>
    <aside id="links">
      <ul>
        <li>
          <router-link :to="{ path: '/search', query: { lang: locale, wordFrom: wordTo, wordTo: wordFrom } }">{{
            t('message.searchInReverse', { wordFrom, wordTo })
          }}</router-link>
        </li>
        <li v-if="routes">
          <a v-t="{ path: 'message.tweet' }" :href="tweetFoundUrl(routes)" target="_blank" />
        </li>
        <li v-else>
          <a v-t="{ path: 'message.tweet' }" :href="tweetNotFoundUrl()" target="_blank" />
        </li>
      </ul>
    </aside>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  wordFrom: string;
  wordTo: string;
}>();

const ErrorCode = {
  NoError: 0,
  NotFoundFrom: 1,
  NotFoundTo: 2,
  NotFoundRoute: 3,
} as const;
type ErrorCodeType = (typeof ErrorCode)[keyof typeof ErrorCode];

interface Result {
  readonly route: ReadonlyArray<string> | undefined;
  readonly error: ErrorCodeType;
}

const { t, locale } = useI18n();

const routes = ref<readonly string[] | undefined>(undefined);
const failureReason = ref<string | undefined>(undefined);
const time = ref(0);

const buildTweetUrl = (text: string, url: string): string => {
  const tweetUrl = new URL('https://twitter.com/intent/tweet');
  tweetUrl.searchParams.append('text', text);
  tweetUrl.searchParams.append('url', url);
  tweetUrl.searchParams.append('hashtags', 'pediaroute');
  return tweetUrl.toString();
};

const search = async () => {
  const body = JSON.stringify({ from: props.wordFrom, to: props.wordTo });
  const headers = { Accept: 'application/json', 'Content-Type': 'application/json' };
  const start = new Date().getTime();
  routes.value = undefined;
  return fetch(`/api/search?lang=${encodeURI(locale.value)}`, { method: 'POST', body, headers })
    .then((response) => response.json())
    .then((result: Result) => {
      routes.value = result.route;
      if (result.error === ErrorCode.NoError) {
        failureReason.value = undefined;
      } else if (result.error === ErrorCode.NotFoundFrom) {
        failureReason.value = t('error.notFoundFrom', { from: props.wordFrom });
      } else if (result.error === ErrorCode.NotFoundTo) {
        failureReason.value = t('error.notFoundTo', { to: props.wordTo });
      } else if (result.error === ErrorCode.NotFoundRoute) {
        failureReason.value = t('error.notFoundRoute');
      }
      time.value = new Date().getTime() - start;
    });
};

const searchPageUrl = (): string =>
  `https://pediaroute.com/search?lang=${encodeURI(locale.value)}&wordFrom=${encodeURIComponent(props.wordFrom)}&wordTo=${encodeURIComponent(props.wordTo)}`;

const tweetFoundUrl = (route: readonly string[]): string =>
  buildTweetUrl(t('message.tweetFind', { wordFrom: props.wordFrom, wordTo: props.wordTo, length: `${route.length - 1}` }), searchPageUrl());

const tweetNotFoundUrl = (): string =>
  buildTweetUrl(t('message.tweetNotFound', { wordFrom: props.wordFrom, wordTo: props.wordTo }), searchPageUrl());

onMounted(search);
watch(() => [props.wordFrom, props.wordTo], search);
</script>
