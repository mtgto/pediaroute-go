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
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useMainStore } from '../store';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

const state = useMainStore();
const { t, locale } = useI18n();
const router = useRouter();

const wordFrom = computed(() => state.wordFrom);
const wordTo = computed(() => state.wordTo);

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
