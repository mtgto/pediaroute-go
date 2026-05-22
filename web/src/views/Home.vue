<template>
  <div>
    <header>
      <h1><span>P</span>edia <span>R</span>oute.com</h1>
      <p>{{ t('message.header') }}</p>
    </header>
    <article>
      <fieldset>
        <p>
          <label>
            <div class="word-input-wrapper">
              <input
                type="text"
                class="word"
                name="wordFrom"
                :value="wordFrom"
                @input="setWordFrom"
                @focus="getRelated(wordFrom, (w) => (relatedFrom = w))"
                @blur="relatedFrom = []"
              />
              <button class="random" :title="t('message.buttonRandom')" @click="getRandomFrom">
                <img src="../assets/baseline-shuffle-24px.svg" />
              </button>
              <ul v-if="showFromDropdown" class="related-dropdown">
                <li v-for="word in relatedFrom" :key="word" @mousedown.prevent @click="selectFrom(word)">{{ word }}</li>
              </ul>
            </div>
            {{ t('message.searchFrom') }}
          </label>
        </p>
        <p>
          <label>
            <div class="word-input-wrapper">
              <input
                type="text"
                class="word"
                name="wordTo"
                :value="wordTo"
                @input="setWordTo"
                @focus="getRelated(wordTo, (w) => (relatedTo = w))"
                @blur="relatedTo = []"
              />
              <button class="random" :title="t('message.buttonRandom')" @click="getRandomTo">
                <img src="../assets/baseline-shuffle-24px.svg" />
              </button>
              <ul v-if="showToDropdown" class="related-dropdown">
                <li v-for="word in relatedTo" :key="word" @mousedown.prevent @click="selectTo(word)">{{ word }}</li>
              </ul>
            </div>
            {{ t('message.searchTo') }}
          </label>
        </p>
      </fieldset>
      <div class="center">
        <input type="button" class="submit" :value="t('message.search')" @click="search" />
      </div>
      <div class="center">
        <p>
          <a v-if="locale === 'ja'" href="#" @click="chooseEnglish">Choose English</a>
          <a v-if="locale === 'en'" href="#" @click="chooseJapanese">日本語を選択</a>
        </p>
      </div>
    </article>
    <footer v-if="currentStats">
      <p>
        {{ locale === 'ja' ? '収録' : 'Index:' }} {{ currentStats.page_count.toLocaleString() }} {{ locale === 'ja' ? '記事' : 'articles' }} ·
        {{ locale === 'ja' ? 'リンク数' : 'Links:' }} {{ currentStats.link_count.toLocaleString() }} · {{ locale === 'ja' ? '版' : 'Version:' }}
        {{ currentStats.version || '...' }}
      </p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useMainStore } from '../store';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

interface LangInfo {
  page_count: number;
  link_count: number;
  version: string;
}

const langInfoMap = ref<Record<string, LangInfo>>({});

onMounted(async () => {
  await fetch('/api/info')
    .then((r) => r.json())
    .then((data: Record<string, LangInfo>) => {
      langInfoMap.value = data;
    })
    .catch(console.error);
});

const state = useMainStore();
const { t, locale } = useI18n();
const router = useRouter();

const wordFrom = computed(() => state.wordFrom);
const wordTo = computed(() => state.wordTo);
const currentStats = computed(() => langInfoMap.value[locale.value]);

const relatedFrom = ref<string[]>([]);
const relatedTo = ref<string[]>([]);

const showFromDropdown = computed(
  () => relatedFrom.value.length > 0 && !relatedFrom.value.some((c) => c.toLowerCase() === wordFrom.value.toLowerCase()),
);
const showToDropdown = computed(() => relatedTo.value.length > 0 && !relatedTo.value.some((c) => c.toLowerCase() === wordTo.value.toLowerCase()));

const relatedCache = new Map<string, string[]>();

const getRelated = async (word: string, setter: (words: string[]) => void) => {
  if ([...word].length < 2) {
    setter([]);
    return;
  }
  const cacheKey = `${locale.value}:${word}`;
  const cached = relatedCache.get(cacheKey);
  if (cached) {
    setter(cached);
    return;
  }
  await fetch(`/api/related?lang=${encodeURI(locale.value)}&title=${encodeURIComponent(word)}`)
    .then((r) => r.json())
    .then((words: string[]) => {
      relatedCache.set(cacheKey, words);
      setter(words);
    })
    .catch(console.error);
};

const chooseEnglish = () => {
  locale.value = 'en';
};
const chooseJapanese = () => {
  locale.value = 'ja';
};
const getRandom = async (setter: (word: string) => void) => {
  await fetch(`/api/random?lang=${encodeURI(locale.value)}`)
    .then((response) => response.json())
    .then((word) => {
      if (typeof word === 'string') setter(word);
    })
    .catch((error) => console.log(error));
};
const getRandomFrom = () =>
  getRandom((w) => {
    state.setWordFrom(w);
    relatedFrom.value = [];
  });
const getRandomTo = () =>
  getRandom((w) => {
    state.setWordTo(w);
    relatedTo.value = [];
  });
const setWordFrom = (e: Event) => {
  if (!(e.target instanceof HTMLInputElement)) return;
  state.setWordFrom(e.target.value);
  getRelated(e.target.value, (w) => (relatedFrom.value = w));
};
const setWordTo = (e: Event) => {
  if (!(e.target instanceof HTMLInputElement)) return;
  state.setWordTo(e.target.value);
  getRelated(e.target.value, (w) => (relatedTo.value = w));
};
const selectFrom = (word: string) => {
  state.setWordFrom(word);
  relatedFrom.value = [];
};
const selectTo = (word: string) => {
  state.setWordTo(word);
  relatedTo.value = [];
};
const search = async () => {
  await router.push({ path: '/search', query: { lang: locale.value, wordFrom: state.wordFrom, wordTo: state.wordTo } });
};
</script>

<style scoped>
.word-input-wrapper {
  position: relative;
  display: inline-block;
}

.related-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  min-width: 400px;
  background: #fff;
  border: 1px solid #ccc;
  list-style: none;
  z-index: 10;
  max-height: 300px;
  overflow-y: auto;
}

.related-dropdown li {
  padding: 6px 10px;
  font-size: 16pt;
  cursor: pointer;
}

.related-dropdown li:hover {
  background: #e8f5f0;
}
</style>
