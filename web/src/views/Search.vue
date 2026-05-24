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
          <RouteList :routes="routes" :wordFrom="wordFrom" :wordTo="wordTo" :shareUrl="shareFoundUrl(routes)" />
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

        <!-- Server error -->
        <template v-else-if="errorCode === ErrorCode.ServerError">
          <Notice :body-large="true" :error-code="ErrorCode.ServerError">
            <template #header-title>{{ t('search.noticeTitle') }}</template>
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
            <Button as="a" variant="ghost" :href="shareNotFoundUrl()" target="_blank" rel="noopener">
              {{ t('message.share') }}
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
import RouteList from '../components/RouteList.vue';
import Stamp from '../components/Stamp.vue';
import { ErrorCode, type ErrorCodeType } from '../types';

const props = defineProps<{
  wordFrom: string;
  wordTo: string;
}>();

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

const buildShareUrl = (text: string, url: string): string => {
  const shareUrl = new URL('https://x.com/intent/post');
  shareUrl.searchParams.append('text', text);
  shareUrl.searchParams.append('url', url);
  shareUrl.searchParams.append('hashtags', 'pediaroute');
  return shareUrl.toString();
};

const doSearch = async () => {
  isLoading.value = true;
  routes.value = undefined;
  failureReason.value = undefined;
  const body = JSON.stringify({ from: props.wordFrom, to: props.wordTo });
  const headers = { Accept: 'application/json', 'Content-Type': 'application/json' };
  const start = Date.now();
  return fetch(`/api/search?lang=${encodeURI(locale.value)}`, { method: 'POST', body, headers })
    .then((r) => {
      if (r.status >= 500) {
        errorCode.value = ErrorCode.ServerError;
        time.value = Date.now() - start;
        return;
      }
      return r.json().then((result: Result) => {
        errorCode.value = result.error;
        routes.value = result.route;
        if (result.error === ErrorCode.NotFoundFrom) {
          failureReason.value = t('error.notFoundFrom', { from: props.wordFrom });
        } else if (result.error === ErrorCode.NotFoundTo) {
          failureReason.value = t('error.notFoundTo', { to: props.wordTo });
        }
        time.value = Date.now() - start;
      });
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

const shareFoundUrl = (route: readonly string[]): string =>
  buildShareUrl(
    t('message.shareFound', {
      wordFrom: props.wordFrom,
      wordTo: props.wordTo,
      length: `${route.length - 1}`,
    }),
    searchPageUrl(),
  );

const shareNotFoundUrl = (): string => buildShareUrl(t('message.shareNotFound', { wordFrom: props.wordFrom, wordTo: props.wordTo }), searchPageUrl());

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

  .actions {
    flex-direction: column;
  }
}
</style>
