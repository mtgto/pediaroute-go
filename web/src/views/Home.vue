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
            <input type="text" class="word" name="wordFrom" :value="wordFrom" @input="setWordFrom" />
            <button class="random" :title="t('message.buttonRandom')" @click="getRandomFrom">
              <img src="../assets/baseline-shuffle-24px.svg" />
            </button>
            {{ t('message.searchFrom') }}
          </label>
        </p>
        <p>
          <label>
            <input type="text" class="word" name="wordTo" :value="wordTo" @input="setWordTo" />
            <button class="random" :title="t('message.buttonRandom')" @click="getRandomTo">
              <img src="../assets/baseline-shuffle-24px.svg" />
            </button>
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
const getRandomFrom = () => getRandom(state.setWordFrom);
const getRandomTo = () => getRandom(state.setWordTo);
const setWordFrom = (e: Event) => {
  if (e.target instanceof HTMLInputElement) state.setWordFrom(e.target.value);
};
const setWordTo = (e: Event) => {
  if (e.target instanceof HTMLInputElement) state.setWordTo(e.target.value);
};
const search = async () => {
  await router.push({ path: '/search', query: { lang: locale.value, wordFrom: state.wordFrom, wordTo: state.wordTo } });
};
</script>
