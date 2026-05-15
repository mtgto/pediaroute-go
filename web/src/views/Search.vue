<template>
  <div class="search">
    <div class="search__inner">
      <!-- Loading -->
      <div v-if="isLoading" class="search__loading">
        <div class="search__loading-label">
          {{ isJa ? '経路を探しています...' : 'Tracing route…' }}
        </div>
        <div class="search__loading-pair">
          <em>{{ wordFrom }}</em>
          <span class="search__arrow">⤳</span>
          <em>{{ wordTo }}</em>
        </div>
      </div>

      <template v-else>
        <!-- Result header -->
        <div class="search__result-header">
          <div>
            <div class="search__result-label">
              <template v-if="errorCode === ErrorCode.NoError">
                {{ isJa ? '蔵書照会 ・ 経路あり' : 'Catalog Entry · Route Found' }}
              </template>
              <template v-else>
                {{ isJa ? '蔵書照会 ・ 経路なし' : 'Catalog Entry · No Route' }}
              </template>
            </div>
            <div class="search__result-pair">
              <em>{{ wordFrom }}</em>
              <span class="search__arrow">→</span>
              <em>{{ wordTo }}</em>
            </div>
          </div>

          <!-- Stamp -->
          <div v-if="errorCode === ErrorCode.NoError" class="lib-stamp lib-stamp--found">
            <div class="lib-stamp__word">{{ isJa ? '発 見' : 'FOUND' }}</div>
            <div class="lib-stamp__meta">{{ (time / 1000).toFixed(3) }} sec · {{ (routes?.length ?? 1) - 1 }} hops</div>
          </div>
          <div v-else class="lib-stamp lib-stamp--notfound">
            <div class="lib-stamp__word">{{ isJa ? '未 到 達' : 'NO ROUTE' }}</div>
            <div class="lib-stamp__meta">{{ (time / 1000).toFixed(3) }} sec</div>
          </div>
        </div>

        <!-- Route found: stacked catalog slips -->
        <template v-if="errorCode === ErrorCode.NoError && routes">
          <div class="route-list">
            <template v-for="(word, i) in routes" :key="word">
              <div :class="['route-slip', i === 0 || i === (routes?.length ?? 1) - 1 ? 'route-slip--endpoint' : '']">
                <div class="route-slip__step">
                  <template v-if="i === 0">{{ isJa ? '出 発 点' : 'Origin' }}</template>
                  <template v-else-if="i === (routes?.length ?? 1) - 1">{{ isJa ? '到 着 点' : 'Destination' }}</template>
                  <template v-else>
                    <template v-if="isJa"
                      >第 <span class="route-slip__step-n">{{ i }}</span> 歩</template
                    >
                    <template v-else>Step {{ String(i).padStart(2, '0') }}</template>
                  </template>
                </div>
                <div class="route-slip__title">
                  <a :href="t('message.wikipediaUrl', { word })" target="_blank" rel="noopener">{{ word }}<span class="route-slip__ext">↗</span></a>
                </div>
                <div class="route-slip__num">{{ i + 1 }}</div>
              </div>
              <div v-if="i < (routes?.length ?? 1) - 1" class="route-connector">↓</div>
            </template>
          </div>

          <!-- Actions for found state -->
          <div class="search__actions">
            <RouterLink class="lib-btn lib-btn--primary" :to="{ path: '/search', query: { lang: locale, wordFrom: wordTo, wordTo: wordFrom } }">
              {{ isJa ? '← 逆 経 路' : '← Reverse Route' }}
            </RouterLink>
            <RouterLink class="lib-btn lib-btn--outline" to="/">
              {{ isJa ? '新しく探す' : 'New Search' }}
            </RouterLink>
            <a class="lib-btn lib-btn--ghost" :href="tweetFoundUrl(routes)" target="_blank" rel="noopener">
              {{ t('message.tweet') }}
            </a>
          </div>
        </template>

        <!-- Article not found -->
        <template v-else-if="errorCode === ErrorCode.NotFoundFrom || errorCode === ErrorCode.NotFoundTo">
          <div class="lib-notice">
            <div class="lib-notice__header">
              <span>{{ isJa ? '蔵 書 通 知' : 'Catalog Notice' }}</span>
              <span class="lib-notice__num">№ 0001 · v</span>
            </div>
            <p class="lib-notice__body">{{ failureReason }}</p>
          </div>
          <div class="search__actions">
            <RouterLink class="lib-btn lib-btn--primary" to="/">
              {{ isJa ? '← 新しく探す' : '← New Search' }}
            </RouterLink>
          </div>
        </template>

        <!-- Route not found -->
        <template v-else-if="errorCode === ErrorCode.NotFoundRoute">
          <div class="lib-notice">
            <div class="lib-notice__header">
              <span>{{ isJa ? '蔵 書 通 知' : 'Catalog Notice' }}</span>
              <span class="lib-notice__num">№ 0001 · v</span>
            </div>
            <p class="lib-notice__body lib-notice__body--large">
              <template v-if="isJa"> この二つの記事を<em class="lib-notice__em">６リンク以内</em>で結ぶ経路は見つかりませんでした。 </template>
              <template v-else> No chain of <em class="lib-notice__em">six links or fewer</em> connects these two articles. </template>
            </p>
            <p class="lib-notice__note">
              <template v-if="isJa">
                記事名の表記揺れか、もしくは到着点が他から孤立した記事である可能性があります。PediaRoute は信念をもって探索を６歩までに留めています —
                それ以上長い経路は、百科事典を当てもなくさまようことになるため。
              </template>
              <template v-else>
                This usually means one of the titles is misspelled, or the goal article is a very isolated entry. PediaRoute caps the search at 6 hops
                on principle — longer chains tend to wander aimlessly across the encyclopedia.
              </template>
            </p>
          </div>
          <div class="search__actions">
            <RouterLink class="lib-btn lib-btn--primary" to="/">
              {{ isJa ? '← 新しく探す' : '← New Search' }}
            </RouterLink>
            <RouterLink class="lib-btn lib-btn--outline" :to="{ path: '/search', query: { lang: locale, wordFrom: wordTo, wordTo: wordFrom } }">
              {{ isJa ? '逆方向を試す' : '↻ Try reverse' }}
            </RouterLink>
            <a class="lib-btn lib-btn--ghost" :href="tweetNotFoundUrl()" target="_blank" rel="noopener">
              {{ t('message.tweet') }}
            </a>
          </div>
        </template>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, watch } from 'vue';
import { RouterLink } from 'vue-router';
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
const isJa = computed(() => locale.value === 'ja');

const isLoading = ref(true);
const routes = ref<readonly string[] | undefined>(undefined);
const failureReason = ref<string | undefined>(undefined);
const errorCode = ref<ErrorCodeType>(ErrorCode.NoError);
const time = ref(0);

const buildTweetUrl = (text: string, url: string): string => {
  const tweetUrl = new URL('https://twitter.com/intent/tweet');
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

.search__inner {
  max-width: 880px;
  margin: 0 auto;
}

/* Loading */
.search__loading {
  padding: 60px 0;
  text-align: center;
}

.search__loading-label {
  font-family: var(--f-mono);
  font-size: 11px;
  letter-spacing: 0.22em;
  text-transform: uppercase;
  color: var(--c-accent);
  margin-bottom: 16px;
}

html.lang-ja .search__loading-label {
  font-family: var(--f-body);
  text-transform: none;
  letter-spacing: 0.3em;
}

.search__loading-pair {
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
.search__result-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 32px;
  padding-bottom: 18px;
  border-bottom: 1px solid var(--c-rule);
  gap: 20px;
}

.search__result-label {
  font-family: var(--f-mono);
  font-size: 10px;
  letter-spacing: 0.22em;
  text-transform: uppercase;
  color: var(--c-accent);
  margin-bottom: 10px;
}

html.lang-ja .search__result-label {
  font-family: var(--f-body);
  font-size: 11px;
  letter-spacing: 0.4em;
  text-transform: none;
}

.search__result-pair {
  font-family: var(--f-head);
  font-size: 26px;
  line-height: 1.3;
  font-weight: 500;
  display: flex;
  align-items: baseline;
  gap: 12px;
  flex-wrap: wrap;
}

.search__arrow {
  color: var(--c-dim);
  font-style: normal;
}

/* Stamp */
.lib-stamp {
  border: 1.5px solid var(--c-accent);
  padding: 8px 14px;
  color: var(--c-accent);
  font-family: var(--f-mono);
  font-size: 11px;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  border-radius: 3px;
  background: rgba(138, 51, 36, 0.04);
  text-align: center;
  flex-shrink: 0;
}

html.lang-ja .lib-stamp {
  font-family: var(--f-head);
  font-size: 14px;
  letter-spacing: 0.4em;
  text-transform: none;
}

.lib-stamp--found {
  transform: rotate(-2deg);
}

.lib-stamp--notfound {
  transform: rotate(-3deg);
}

.lib-stamp__word {
  font-weight: 600;
  font-size: 13px;
  letter-spacing: 0.16em;
}

html.lang-ja .lib-stamp__word {
  font-size: 14px;
  letter-spacing: 0.5em;
}

.lib-stamp__meta {
  font-size: 9px;
  opacity: 0.8;
  margin-top: 1px;
  font-family: var(--f-mono);
  letter-spacing: 0.08em;
}

/* Route list */
.route-list {
  display: flex;
  flex-direction: column;
}

.route-slip {
  display: grid;
  grid-template-columns: 72px 1fr auto;
  align-items: center;
  gap: 24px;
  padding: 18px 28px;
  background: var(--c-paper);
  border: 1px solid var(--c-rule);
  border-radius: 2px;
}

.route-slip--endpoint {
  background: transparent;
  border-style: dashed;
}

html.lang-ja .route-slip {
  grid-template-columns: 88px 1fr auto;
}

.route-slip__step {
  font-family: var(--f-mono);
  font-size: 11px;
  color: var(--c-dim);
  letter-spacing: 0.1em;
  text-transform: uppercase;
  border-right: 1px solid var(--c-rule);
  padding-right: 16px;
  line-height: 1.4;
}

html.lang-ja .route-slip__step {
  font-family: var(--f-body);
  font-size: 12px;
  letter-spacing: 0.18em;
  text-transform: none;
}

.route-slip__step-n {
  font-family: var(--f-serif);
  font-style: italic;
}

.route-slip__title {
  font-family: var(--f-head);
  font-size: 22px;
  line-height: 1.2;
  font-weight: 500;
}

.route-slip__title a {
  color: var(--c-ink);
}

.route-slip__title a:hover {
  color: var(--c-accent);
  text-decoration: none;
}

.route-slip__ext {
  margin-left: 6px;
  font-family: var(--f-mono);
  font-size: 11px;
  color: var(--c-dim);
  font-weight: 400;
}

.route-slip__num {
  font-family: var(--f-serif);
  font-style: italic;
  font-size: 32px;
  color: var(--c-accent);
  opacity: 0.35;
}

.route-connector {
  height: 16px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-family: var(--f-mono);
  color: var(--c-accent);
  font-size: 12px;
}

/* Catalog notice */
.lib-notice {
  background: var(--c-paper);
  border: 1px solid var(--c-rule);
  padding: 28px 36px;
  border-radius: 2px;
  margin-bottom: 36px;
  box-shadow:
    0 1px 0 rgba(28, 27, 24, 0.04),
    0 12px 24px -18px rgba(28, 27, 24, 0.18);
}

.lib-notice__header {
  font-family: var(--f-mono);
  font-size: 10px;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--c-dim);
  margin-bottom: 16px;
  display: flex;
  justify-content: space-between;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--c-rule);
}

html.lang-ja .lib-notice__header {
  font-family: var(--f-body);
  font-size: 11px;
  letter-spacing: 0.32em;
  text-transform: none;
}

.lib-notice__num {
  font-family: var(--f-mono);
  letter-spacing: 0.12em;
}

.lib-notice__body {
  font-family: var(--f-body);
  font-size: 17px;
  line-height: 1.6;
  color: var(--c-ink);
}

.lib-notice__body--large {
  font-size: 21px;
  line-height: 1.45;
  margin-bottom: 18px;
}

html.lang-ja .lib-notice__body--large {
  font-family: var(--f-head);
  font-size: 22px;
  line-height: 1.7;
  font-weight: 500;
}

.lib-notice__em {
  color: var(--c-accent);
}

.lib-notice__note {
  font-family: var(--f-body);
  font-size: 16px;
  line-height: 1.65;
  color: var(--c-dim);
  max-width: 640px;
}

html.lang-ja .lib-notice__note {
  font-size: 14px;
  line-height: 1.9;
}

/* Actions */
.search__actions {
  margin-top: 36px;
  display: flex;
  gap: 14px;
  flex-wrap: wrap;
}

.lib-btn {
  display: inline-block;
  padding: 12px 20px;
  font-family: var(--f-mono);
  font-size: 11px;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  cursor: pointer;
  text-decoration: none;
}

html.lang-ja .lib-btn {
  font-family: var(--f-head);
  font-size: 13px;
  letter-spacing: 0.32em;
  text-transform: none;
}

.lib-btn--primary {
  background: var(--c-ink);
  color: var(--c-paper);
}

.lib-btn--primary:hover {
  background: var(--c-accent);
  text-decoration: none;
  color: var(--c-paper);
}

.lib-btn--outline {
  border: 1px solid var(--c-ink);
  color: var(--c-ink);
}

.lib-btn--outline:hover {
  border-color: var(--c-accent);
  color: var(--c-accent);
  text-decoration: none;
}

.lib-btn--ghost {
  color: var(--c-dim);
  border-bottom: 1px solid var(--c-rule);
}

.lib-btn--ghost:hover {
  color: var(--c-ink);
  text-decoration: none;
}

/* Mobile */
@media (max-width: 640px) {
  .search {
    padding: 24px 18px 32px;
  }

  .search__result-header {
    flex-wrap: wrap;
  }

  .search__result-pair {
    font-size: 20px;
    gap: 8px;
  }

  .route-slip {
    grid-template-columns: 56px 1fr auto;
    gap: 14px;
    padding: 14px 16px;
  }

  html.lang-ja .route-slip {
    grid-template-columns: 72px 1fr auto;
  }

  .route-slip__title {
    font-size: 17px;
  }

  .route-slip__num {
    font-size: 22px;
  }

  .lib-notice {
    padding: 20px 18px;
  }

  .search__actions {
    flex-direction: column;
  }

  .lib-btn {
    text-align: center;
  }
}
</style>
