<template>
  <div>
    <header>
      <h1><span>P</span>edia <span>R</span>oute.com</h1>
      <p>{{ $t('message.header') }}</p>
    </header>
    <article>
      <fieldset>
        <p>
          <label>
            <input type="text" class="word" name="wordFrom" :value="wordFrom" @input="setWordFrom" />
            <button class="random" :title="$t('message.buttonRandom')" @click="getRandomFrom">
              <img src="../assets/baseline-shuffle-24px.svg" />
            </button>
            {{ $t('message.searchFrom') }}
          </label>
        </p>
        <p>
          <label>
            <input type="text" class="word" name="wordTo" :value="wordTo" @input="setWordTo" />
            <button class="random" :title="$t('message.buttonRandom')" @click="getRandomTo">
              <img src="../assets/baseline-shuffle-24px.svg" />
            </button>
            {{ $t('message.searchTo') }}
          </label>
        </p>
      </fieldset>
      <div class="center">
        <input type="button" class="submit" :value="$t('message.search')" @click="search" />
      </div>
      <div class="center">
        <p>
          <a v-if="$i18n.locale === 'ja'" href="#" @click="chooseEnglish">Choose English</a>
          <a v-if="$i18n.locale === 'en'" href="#" @click="chooseJapanese">日本語を選択</a>
        </p>
      </div>
    </article>
  </div>
</template>

<script lang="ts">
import { defineComponent, computed } from 'vue';
import { useMainStore } from '../store';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

export default defineComponent({
  setup() {
    const state = useMainStore();
    const i18n = useI18n();
    const router = useRouter();

    const wordFrom = computed((): string => state.wordFrom);
    const wordTo = computed((): string => state.wordTo);
    const chooseEnglish = () => {
      i18n.locale.value = 'en';
    };
    const chooseJapanese = () => {
      i18n.locale.value = 'ja';
    };
    const getRandomFrom = async () => {
      await fetch(`/api/random?lang=${encodeURI(i18n.locale.value)}`)
        .then((response) => response.json())
        .then((word) => {
          if (typeof word === 'string') {
            state.setWordFrom(word);
          }
        })
        .catch((error) => console.log(error));
    };
    const getRandomTo = async () => {
      await fetch(`/api/random?lang=${encodeURI(i18n.locale.value)}`)
        .then((response) => response.json())
        .then((word) => {
          if (typeof word === 'string') {
            state.setWordTo(word);
          }
        })
        .catch((error) => console.log(error));
    };
    const setWordFrom = (e: Event) => {
      if (e.target instanceof HTMLInputElement) {
        state.setWordFrom(e.target.value);
      }
    };
    const setWordTo = (e: Event) => {
      if (e.target instanceof HTMLInputElement) {
        state.setWordTo(e.target.value);
      }
    };
    const search = async () => {
      await router.push({ path: '/search', query: { lang: i18n.locale.value, wordFrom: state.wordFrom, wordTo: state.wordTo } });
    };
    return {
      wordFrom,
      wordTo,
      chooseEnglish,
      chooseJapanese,
      getRandomFrom,
      getRandomTo,
      setWordFrom,
      setWordTo,
      search,
    };
  },
});
</script>
