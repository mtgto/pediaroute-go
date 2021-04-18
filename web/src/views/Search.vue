<template>
  <div>
    <header>
      <h1><span>P</span>edia <span>R</span>oute.com</h1>
      <p v-if="routes" v-text="t('message.searchResult', { wordFrom, wordTo, second: time / 1000 })" />
      <p v-else v-text="t('message.searching', { wordFrom, wordTo })" />
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
          <router-link
            :to="{ path: '/search', query: { lang: $i18n.locale, wordFrom: wordTo, wordTo: wordFrom } }"
            v-text="t('message.searchInReverse', { wordFrom, wordTo })"
          />
        </li>
        <li v-if="routes">
          <a :href="tweetFoundUrl(routes)" target="_blank" v-text="t('message.tweet')" />
        </li>
        <li v-else>
          <a :href="tweetNotFoundUrl()" target="_blank" v-text="t('message.tweet')" />
        </li>
      </ul>
    </aside>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, Ref, onMounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';

export const ErrorCode = {
  NoError: 0,
  NotFoundFrom: 1,
  NotFoundTo: 2,
  NotFoundRoute: 3,
} as const;
type ErrorCodeType = typeof ErrorCode[keyof typeof ErrorCode];

interface Result {
  readonly route: ReadonlyArray<string> | undefined;
  readonly error: ErrorCodeType;
}

type Props = {
  wordFrom: string;
  wordTo: string;
};

export default defineComponent({
  props: {
    wordFrom: {
      type: String,
      required: true,
    },
    wordTo: {
      type: String,
      required: true,
    },
  },
  setup(props: Props) {
    const i18n = useI18n();
    // data
    const routes: Ref<readonly string[] | undefined> = ref(undefined);
    const failureReason: Ref<string | undefined> = ref(undefined);
    const time = ref(0);
    // methods
    const search = async () => {
      const body = JSON.stringify({ from: props.wordFrom, to: props.wordTo });
      const headers = {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      };
      const start = new Date().getTime();
      routes.value = undefined;
      return fetch(`/api/search?lang=${encodeURI(i18n.locale.value)}`, { method: 'POST', body, headers })
        .then((response) => response.json())
        .then((result: Result) => {
          routes.value = result.route;
          if (result.error === ErrorCode.NoError) {
            failureReason.value = undefined;
          } else if (result.error === ErrorCode.NotFoundFrom) {
            const message = i18n.t('error.notFoundFrom', { from: props.wordFrom });
            if (typeof message === 'string') {
              failureReason.value = message;
            }
          } else if (result.error === ErrorCode.NotFoundTo) {
            const message = i18n.t('error.notFoundTo', { to: props.wordTo });
            if (typeof message === 'string') {
              failureReason.value = message;
            }
          } else if (result.error === ErrorCode.NotFoundRoute) {
            const message = i18n.t('error.notFoundRoute');
            if (typeof message === 'string') {
              failureReason.value = message;
            }
          }
          time.value = new Date().getTime() - start;
        });
    };
    const tweetFoundUrl = (route: string[]): string => {
      return i18n.t('message.tweetFind', {
        wordFrom: props.wordFrom,
        wordTo: props.wordTo,
        length: `${route && route.length - 1}`,
        link: `${encodeURIComponent(
          `https://pediaroute.com/search?lang=${encodeURI(i18n.locale.value)}&wordFrom=${encodeURIComponent(
            props.wordFrom,
          )}&wordTo=${encodeURIComponent(props.wordTo)}`,
        )}`,
        hashTag: encodeURIComponent('#pediaroute'),
      });
    };
    const tweetNotFoundUrl = (): string => {
      return i18n.t('message.tweetNotFound', {
        wordFrom: props.wordFrom,
        wordTo: props.wordTo,
        link: `${encodeURIComponent(
          `https://pediaroute.com/search?lang=${encodeURI(i18n.locale.value)}&wordFrom=${encodeURIComponent(
            props.wordFrom,
          )}&wordTo=${encodeURIComponent(props.wordTo)}`,
        )}`,
        hashTag: encodeURIComponent('#pediaroute'),
      });
    };
    onMounted(search);
    watch(
      () => [props.wordFrom, props.wordTo],
      () => {
        search();
      },
    );
    return { routes, failureReason, time, search, tweetFoundUrl, tweetNotFoundUrl, t: i18n.t };
  },
});
</script>
