<template>
  <div class="home">
    <div class="inner">
      <p class="tagline">{{ t('home.tagline') }}</p>

      <i18n-t keypath="home.hero" tag="h1" class="hero">
        <template #break><br /></template>
        <template #em
          ><em class="hero-em">{{ t('home.heroEm') }}</em></template
        >
      </i18n-t>

      <p class="subtitle">{{ t('home.subtitle') }}</p>

      <!-- Catalog card -->
      <Card>
        <template #header-title>{{ t('home.catalogTitle') }}</template>

        <!-- From field -->
        <div class="field-row">
          <Field
            :label="t('home.fieldFrom')"
            type="text"
            v-model="wordFrom"
            :placeholder="t('home.placeholder')"
            :suggestions="suggestionsFrom"
            @keydown.ctrl.enter.prevent="search"
            @keydown.meta.enter.prevent="search"
          />
          <button class="random-btn" :title="t('message.buttonRandom')" @click="getRandomFrom">↻ {{ t('home.random') }}</button>
        </div>

        <!-- To field -->
        <div class="field-row">
          <Field
            :label="t('home.fieldTo')"
            type="text"
            v-model="wordTo"
            :placeholder="t('home.placeholder')"
            :suggestions="suggestionsTo"
            @keydown.ctrl.enter.prevent="search"
            @keydown.meta.enter.prevent="search"
          />
          <button class="random-btn" :title="t('message.buttonRandom')" @click="getRandomTo">↻ {{ t('home.random') }}</button>
        </div>

        <!-- Actions row -->
        <div class="card-actions">
          <button class="submit" @click="search">
            {{ t('home.cta') }}
          </button>
        </div>
      </Card>

      <!-- Stats bar -->
      <div class="stats">
        <span>{{ t('home.statsArticles', { count: pageCount }) }}</span>
        <span v-if="version" class="stats-sep">·</span>
        <span v-if="version">{{ version }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { storeToRefs } from 'pinia';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useMainStore } from '../store';
import { useRelated } from '../composables/useRelated';
import Card from '../components/Card.vue';
import Field from '../components/Field.vue';

interface LangInfo {
  page_count: number;
  link_count: number;
  version: string;
}

const { t, locale } = useI18n();
const router = useRouter();
const store = useMainStore();
const { wordFrom, wordTo } = storeToRefs(store);

const { suggestions: suggestionsFrom, fetch: fetchFrom, clear: clearFrom } = useRelated(locale);
const { suggestions: suggestionsTo, fetch: fetchTo, clear: clearTo } = useRelated(locale);

watch(wordFrom, (val) => fetchFrom(val));
watch(wordTo, (val) => fetchTo(val));

const langInfoMap = ref<Record<string, LangInfo>>({});

onMounted(() => {
  fetch('/api/info')
    .then((r) => r.json())
    .then((data: Record<string, LangInfo>) => {
      langInfoMap.value = data;
    })
    .catch(console.error);
});

const pageCount = computed(() => {
  const info = langInfoMap.value[locale.value];
  return info ? info.page_count.toLocaleString() : '…';
});

const version = computed(() => langInfoMap.value[locale.value]?.version ?? '');

const getRandom = async (setter: (word: string) => void) => {
  await fetch(`/api/random?lang=${encodeURI(locale.value)}`)
    .then((r) => r.json())
    .then((word) => {
      if (typeof word === 'string') setter(word);
    })
    .catch(console.error);
};
const getRandomFrom = () =>
  getRandom((w) => {
    wordFrom.value = w;
    clearFrom();
  });
const getRandomTo = () =>
  getRandom((w) => {
    wordTo.value = w;
    clearTo();
  });

const search = async () => {
  if (!wordFrom.value.trim() || !wordTo.value.trim()) return;
  await router.push({
    path: '/search',
    query: { lang: locale.value, wordFrom: wordFrom.value, wordTo: wordTo.value },
  });
};
</script>

<style scoped>
.home {
  padding: 56px 64px 40px;
}

.inner {
  max-width: 760px;
  margin: 0 auto;
}

/* Tagline */
.tagline {
  font-family: var(--f-mono);
  font-size: 11px;
  letter-spacing: 0.22em;
  text-transform: uppercase;
  color: var(--c-accent);
  margin-bottom: 18px;
}

/* Hero heading */
.hero {
  font-family: var(--f-head);
  font-weight: 400;
  font-size: 44px;
  line-height: 1.15;
  margin-bottom: 14px;
  letter-spacing: -0.01em;
}

html.lang-ja .hero {
  font-weight: 500;
  font-size: 42px;
  line-height: 1.35;
  letter-spacing: 0.02em;
}

.hero-em {
  color: var(--c-accent);
  font-style: italic;
}

html.lang-ja .hero-em {
  font-style: normal;
}

/* Subtitle */
.subtitle {
  font-family: var(--f-body);
  font-size: 17px;
  line-height: 1.6;
  color: var(--c-dim);
  margin-bottom: 40px;
  max-width: 560px;
}

html.lang-ja .subtitle {
  line-height: 1.9;
}

/* Field row (Field + random button) */
.field-row {
  display: flex;
  align-items: flex-end;
  gap: 12px;
}

/* Random button */
.random-btn {
  all: unset;
  cursor: pointer;
  padding: 6px 10px;
  border: 1px solid var(--c-rule);
  border-radius: 2px;
  font-family: var(--ui-font);
  font-size: 10px;
  letter-spacing: 0.16em;
  text-transform: var(--ui-tt);
  color: var(--c-dim);
  background: var(--c-bg);
  white-space: nowrap;
  flex-shrink: 0;
}

html.lang-ja .random-btn {
  font-size: 11px;
  letter-spacing: 0.2em;
}

.random-btn:hover {
  color: var(--c-ink);
  border-color: var(--c-ink);
}

/* Actions row inside card */
.card-actions {
  display: flex;
  justify-content: flex-end;
  padding-top: 8px;
}

/* Submit button */
.submit {
  all: unset;
  cursor: pointer;
  padding: 14px 28px;
  background: var(--c-ink);
  color: var(--c-paper);
  font-family: var(--cta-font);
  font-size: 12px;
  letter-spacing: 0.22em;
  text-transform: var(--ui-tt);
  font-weight: 500;
}

html.lang-ja .submit {
  font-size: 14px;
  letter-spacing: 0.32em;
}

.submit:hover {
  background: var(--c-accent);
}

/* Stats bar */
.stats {
  margin-top: 28px;
  font-family: var(--f-mono);
  font-size: 11px;
  color: var(--c-dim);
  display: flex;
  justify-content: center;
  gap: 16px;
  flex-wrap: wrap;
}

html.lang-ja .stats {
  font-family: var(--f-body);
}

.stats-sep {
  opacity: 0.4;
}

/* Mobile */
@media (max-width: 640px) {
  .home {
    padding: 28px 18px 32px;
  }

  .hero {
    font-size: 28px;
  }

  html.lang-ja .hero {
    font-size: 26px;
  }

  .subtitle {
    font-size: 15px;
    margin-bottom: 28px;
  }

  .card-actions {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .submit {
    text-align: center;
    width: 100%;
  }
}
</style>
