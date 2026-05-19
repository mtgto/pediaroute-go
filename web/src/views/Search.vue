<template>
  <div class="search">
    <div class="inner">
      <!-- Loading -->
      <div v-if="isLoading" class="loading">
        <div class="loading-label">
          {{ t('search.loading') }}
        </div>
        <div class="loading-pair">
          <em>{{ wordFrom }}</em>
          <span class="arrow">⤳</span>
          <em>{{ wordTo }}</em>
        </div>
      </div>

      <template v-else>
        <!-- Result header -->
        <div class="result-header">
          <div>
            <div class="result-label">
              <template v-if="errorCode === ErrorCode.NoError">
                {{ t('search.routeFound') }}
              </template>
              <template v-else>
                {{ t('search.noRoute') }}
              </template>
            </div>
            <div class="result-pair">
              <em>{{ wordFrom }}</em>
              <span class="arrow">→</span>
              <em>{{ wordTo }}</em>
            </div>
          </div>

          <Stamp
            :found="errorCode === ErrorCode.NoError"
            :meta="
              errorCode === ErrorCode.NoError
                ? `${(time / 1000).toFixed(3)} sec · ${(routes?.length ?? 1) - 1} hops`
                : `${(time / 1000).toFixed(3)} sec`
            "
          />
        </div>

        <!-- Route found: stacked catalog slips -->
        <template v-if="errorCode === ErrorCode.NoError && routes">
          <div class="list">
            <template v-for="(word, i) in routes" :key="word">
              <div :class="['slip', { endpoint: i === 0 || i === (routes?.length ?? 1) - 1 }]">
                <div class="step">
                  <template v-if="i === 0">{{ t('search.origin') }}</template>
                  <template v-else-if="i === (routes?.length ?? 1) - 1">{{ t('search.destination') }}</template>
                  <template v-else>
                    <i18n-t keypath="search.step">
                      <template #num
                        ><span class="step-n">{{ i }}</span></template
                      >
                    </i18n-t>
                  </template>
                </div>
                <div class="title">
                  <a :href="t('message.wikipediaUrl', { word })" target="_blank" rel="noopener">{{ word }}<span class="ext">↗</span></a>
                </div>
                <div class="num">{{ i + 1 }}</div>
              </div>
              <div v-if="i < (routes?.length ?? 1) - 1" class="connector">↓</div>
            </template>
          </div>

          <div class="actions">
            <Button :as="RouterLink" variant="primary" :to="{ path: '/search', query: { lang: locale, wordFrom: wordTo, wordTo: wordFrom } }">
              {{ t('search.reverseRoute') }}
            </Button>
            <Button :as="RouterLink" variant="outline" to="/">
              {{ t('search.newSearch') }}
            </Button>
            <Button as="a" variant="ghost" :href="tweetFoundUrl(routes)" target="_blank" rel="noopener">
              {{ t('message.tweet') }}
            </Button>
          </div>
        </template>

        <!-- Article not found -->
        <template v-else-if="errorCode === ErrorCode.NotFoundFrom || errorCode === ErrorCode.NotFoundTo">
          <Notice>
            <template #header-title>{{ t('search.noticeTitle') }}</template>
            <template #body>{{ failureReason }}</template>
          </Notice>
          <div class="actions">
            <Button :as="RouterLink" variant="primary" to="/">
              {{ t('search.newSearchBack') }}
            </Button>
          </div>
        </template>

        <!-- Route not found -->
        <template v-else-if="errorCode === ErrorCode.NotFoundRoute">
          <Notice :body-large="true">
            <template #header-title>{{ t('search.noticeTitle') }}</template>
            <template #body>
              <i18n-t keypath="search.notFoundBody">
                <template #em
                  ><em class="notice-em">{{ t('search.notFoundBodyEm') }}</em></template
                >
              </i18n-t>
            </template>
            <template #note>{{ t('search.notFoundNote') }}</template>
          </Notice>
          <div class="actions">
            <Button :as="RouterLink" variant="primary" to="/">
              {{ t('search.newSearchBack') }}
            </Button>
            <Button :as="RouterLink" variant="outline" :to="{ path: '/search', query: { lang: locale, wordFrom: wordTo, wordTo: wordFrom } }">
              {{ t('search.tryReverse') }}
            </Button>
            <Button as="a" variant="ghost" :href="tweetNotFoundUrl()" target="_blank" rel="noopener">
              {{ t('message.tweet') }}
            </Button>
          </div>
        </template>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { RouterLink } from 'vue-router';
import { useI18n } from 'vue-i18n';
import Button from '../components/Button.vue';
import Notice from '../components/Notice.vue';
import Stamp from '../components/Stamp.vue';

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

const isLoading = ref(true);
const routes = ref<readonly string[] | undefined>(undefined);
const failureReason = ref<string | undefined>(undefined);
const errorCode = ref<ErrorCodeType>(ErrorCode.NoError);
const time = ref(0);

const buildTweetUrl = (text: string, url: string): string => {
  const tweetUrl = new URL('https://x.com/intent/post');
  tweetUrl.searchParams.append('text', text);
  tweetUrl.searchParams.append('url', url);
  tweetUrl.searchParams.append('hashtags', 'pediaroute');
  return tweetUrl.toString();
};

const doSearch = async () => {
  isLoading.value = true;
  routes.value = undefined;
  failureReason.value = undefined;
  const body = JSON.stringify({ from: props.wordFrom, to: props.wordTo });
  const headers = { Accept: 'application/json', 'Content-Type': 'application/json' };
  const start = Date.now();
  return fetch(`/api/search?lang=${encodeURI(locale.value)}`, { method: 'POST', body, headers })
    .then((r) => r.json())
    .then((result: Result) => {
      errorCode.value = result.error;
      routes.value = result.route;
      if (result.error === ErrorCode.NotFoundFrom) {
        failureReason.value = t('error.notFoundFrom', { from: props.wordFrom });
      } else if (result.error === ErrorCode.NotFoundTo) {
        failureReason.value = t('error.notFoundTo', { to: props.wordTo });
      }
      time.value = Date.now() - start;
    })
    .finally(() => {
      isLoading.value = false;
    });
};

const searchPageUrl = (): string => {
  const url = new URL('https://pediaroute.com/search');
  url.searchParams.set('lang', locale.value);
  url.searchParams.set('wordFrom', props.wordFrom);
  url.searchParams.set('wordTo', props.wordTo);
  return url.toString();
};

const tweetFoundUrl = (route: readonly string[]): string =>
  buildTweetUrl(
    t('message.tweetFind', {
      wordFrom: props.wordFrom,
      wordTo: props.wordTo,
      length: `${route.length - 1}`,
    }),
    searchPageUrl(),
  );

const tweetNotFoundUrl = (): string => buildTweetUrl(t('message.tweetNotFound', { wordFrom: props.wordFrom, wordTo: props.wordTo }), searchPageUrl());

onMounted(doSearch);
watch(() => [props.wordFrom, props.wordTo], doSearch);
</script>

<style scoped>
.search {
  padding: 40px 64px 40px;
}

.inner {
  max-width: 880px;
  margin: 0 auto;
}

/* Loading */
.loading {
  padding: 60px 0;
  text-align: center;
}

.loading-label {
  font-family: var(--f-mono);
  font-size: 11px;
  letter-spacing: 0.22em;
  text-transform: uppercase;
  color: var(--c-accent);
  margin-bottom: 16px;
}

html.lang-ja .loading-label {
  font-family: var(--f-body);
  text-transform: none;
  letter-spacing: 0.3em;
}

.loading-pair {
  font-family: var(--f-head);
  font-size: 24px;
  color: var(--c-dim);
  display: flex;
  gap: 12px;
  justify-content: center;
  align-items: center;
  flex-wrap: wrap;
}

/* Result header */
.result-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 32px;
  padding-bottom: 18px;
  border-bottom: 1px solid var(--c-rule);
  gap: 20px;
}

.result-label {
  font-family: var(--f-mono);
  font-size: 10px;
  letter-spacing: 0.22em;
  text-transform: uppercase;
  color: var(--c-accent);
  margin-bottom: 10px;
}

html.lang-ja .result-label {
  font-family: var(--f-body);
  font-size: 11px;
  letter-spacing: 0.4em;
  text-transform: none;
}

.result-pair {
  font-family: var(--f-head);
  font-size: 26px;
  line-height: 1.3;
  font-weight: 500;
  display: flex;
  align-items: baseline;
  gap: 12px;
  flex-wrap: wrap;
}

.arrow {
  color: var(--c-dim);
  font-style: normal;
}

/* Route list */
.list {
  display: flex;
  flex-direction: column;
}

.slip {
  display: grid;
  grid-template-columns: 72px 1fr auto;
  align-items: center;
  gap: 24px;
  padding: 18px 28px;
  background: var(--c-paper);
  border: 1px solid var(--c-rule);
  border-radius: 2px;
}

.slip.endpoint {
  background: transparent;
  border-style: dashed;
}

html.lang-ja .slip {
  grid-template-columns: 88px 1fr auto;
}

.step {
  font-family: var(--f-mono);
  font-size: 11px;
  color: var(--c-dim);
  letter-spacing: 0.1em;
  text-transform: uppercase;
  border-right: 1px solid var(--c-rule);
  padding-right: 16px;
  line-height: 1.4;
}

html.lang-ja .step {
  font-family: var(--f-body);
  font-size: 12px;
  letter-spacing: 0.18em;
  text-transform: none;
}

.step-n {
  font-family: var(--f-serif);
  font-style: italic;
}

.title {
  font-family: var(--f-head);
  font-size: 22px;
  line-height: 1.2;
  font-weight: 500;
}

.title a {
  color: var(--c-ink);
}

.title a:hover {
  color: var(--c-accent);
  text-decoration: none;
}

.ext {
  margin-left: 6px;
  font-family: var(--f-mono);
  font-size: 11px;
  color: var(--c-dim);
  font-weight: 400;
}

.num {
  font-family: var(--f-serif);
  font-style: italic;
  font-size: 32px;
  color: var(--c-accent);
  opacity: 0.35;
}

.connector {
  height: 16px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-family: var(--f-mono);
  color: var(--c-accent);
  font-size: 12px;
}

/* Actions */
.actions {
  margin-top: 36px;
  display: flex;
  gap: 14px;
  flex-wrap: wrap;
}

/* Accent emphasis inside notice body slot */
.notice-em {
  color: var(--c-accent);
}

/* Mobile */
@media (max-width: 640px) {
  .search {
    padding: 24px 18px 32px;
  }

  .result-header {
    flex-wrap: wrap;
  }

  .result-pair {
    font-size: 20px;
    gap: 8px;
  }

  .slip {
    grid-template-columns: 56px 1fr auto;
    gap: 14px;
    padding: 14px 16px;
  }

  html.lang-ja .slip {
    grid-template-columns: 72px 1fr auto;
  }

  .title {
    font-size: 17px;
  }

  .num {
    font-size: 22px;
  }

  .actions {
    flex-direction: column;
  }
}
</style>
