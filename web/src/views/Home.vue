<template>
  <div class="home">
    <div class="home__inner">
      <p class="home__tagline">
        {{ isJa ? 'ウィキペディア経路探索 ・ 2013年から' : 'The Wikipedia Route Finder · Since 2013' }}
      </p>

      <h1 class="home__hero">
        <template v-if="isJa">
          ウィキペディアの二つの記事を、<br />
          <em class="home__hero-em">６つのリンクで</em>結ぶ。
        </template>
        <template v-else>
          Find a route between any two<br />
          Wikipedia pages, in <em class="home__hero-em">6 links or fewer.</em>
        </template>
      </h1>

      <p class="home__subtitle">
        <template v-if="isJa">
          ２つの記事名を入力してください。またはサイコロを振って任意の記事を選んでも構いません。記事から記事へ、ハイパーリンクの最短経路をたどります。
        </template>
        <template v-else> Type two article titles below — or roll the dice. We'll trace the shortest chain of hyperlinks between them. </template>
      </p>

      <!-- Catalog card -->
      <div class="lib-card">
        <div class="lib-card__header">
          <span>{{ isJa ? '蔵 書 照 会' : 'Catalog Lookup' }}</span>
          <span class="lib-card__num">№ 0001</span>
        </div>
        <div class="lib-card__body">
          <!-- From field -->
          <div class="lib-field-row">
            <div class="lib-field">
              <div class="lib-field__label">{{ isJa ? '出 発 点' : 'From' }}</div>
              <input
                class="lib-field__input"
                type="text"
                :value="wordFrom"
                :placeholder="isJa ? '記事名を入力' : 'Article title'"
                @input="onInputFrom"
                @keydown.ctrl.enter.prevent="search"
                @keydown.meta.enter.prevent="search"
              />
            </div>
            <button class="lib-random-btn" :title="t('message.buttonRandom')" @click="getRandomFrom">↻ {{ isJa ? 'ランダム' : 'Random' }}</button>
          </div>

          <!-- To field -->
          <div class="lib-field-row">
            <div class="lib-field">
              <div class="lib-field__label">{{ isJa ? '到 着 点' : 'To' }}</div>
              <input
                class="lib-field__input"
                type="text"
                :value="wordTo"
                :placeholder="isJa ? '記事名を入力' : 'Article title'"
                @input="onInputTo"
                @keydown.ctrl.enter.prevent="search"
                @keydown.meta.enter.prevent="search"
              />
            </div>
            <button class="lib-random-btn" :title="t('message.buttonRandom')" @click="getRandomTo">↻ {{ isJa ? 'ランダム' : 'Random' }}</button>
          </div>

          <!-- Actions row -->
          <div class="lib-card__actions">
            <span class="lib-card__hint">
              <span class="lib-card__hint-mono">⌘ + Enter</span>
              {{ isJa ? 'で経路を探す' : 'to trace route' }}
            </span>
            <button class="lib-submit" @click="search">
              {{ isJa ? '経 路 を 探 す →' : 'Trace Route →' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Stats bar -->
      <div class="home__stats">
        <span><template v-if="isJa">収録</template><template v-else>Index:</template> 6,847,221 articles</span>
        <span class="home__stats-dot">·</span>
        <span><template v-if="isJa">平均</template><template v-else>Avg. trace:</template> 87 ms</span>
        <span class="home__stats-dot">·</span>
        <span><template v-if="isJa">最大</template><template v-else>Max hops:</template> 6</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useMainStore } from '../store';

const { t, locale } = useI18n();
const router = useRouter();
const store = useMainStore();

const isJa = computed(() => locale.value === 'ja');
const wordFrom = computed(() => store.wordFrom);
const wordTo = computed(() => store.wordTo);

const onInputFrom = (e: Event) => {
  if (e.target instanceof HTMLInputElement) store.setWordFrom(e.target.value);
};
const onInputTo = (e: Event) => {
  if (e.target instanceof HTMLInputElement) store.setWordTo(e.target.value);
};

const getRandom = async (setter: (word: string) => void) => {
  await fetch(`/api/random?lang=${encodeURI(locale.value)}`)
    .then((r) => r.json())
    .then((word) => {
      if (typeof word === 'string') setter(word);
    })
    .catch(console.error);
};
const getRandomFrom = () => getRandom(store.setWordFrom);
const getRandomTo = () => getRandom(store.setWordTo);

const search = async () => {
  if (!store.wordFrom.trim() || !store.wordTo.trim()) return;
  await router.push({
    path: '/search',
    query: { lang: locale.value, wordFrom: store.wordFrom, wordTo: store.wordTo },
  });
};
</script>

<style scoped>
.home {
  padding: 56px 64px 40px;
}

.home__inner {
  max-width: 760px;
  margin: 0 auto;
}

/* Tagline */
.home__tagline {
  font-family: var(--f-mono);
  font-size: 11px;
  letter-spacing: 0.22em;
  text-transform: uppercase;
  color: var(--c-accent);
  margin-bottom: 18px;
}

/* Hero heading */
.home__hero {
  font-family: var(--f-head);
  font-weight: 400;
  font-size: 44px;
  line-height: 1.15;
  margin-bottom: 14px;
  letter-spacing: -0.01em;
}

html.lang-ja .home__hero {
  font-weight: 500;
  font-size: 42px;
  line-height: 1.35;
  letter-spacing: 0.02em;
}

.home__hero-em {
  color: var(--c-accent);
  font-style: italic;
}

html.lang-ja .home__hero-em {
  font-style: normal;
}

/* Subtitle */
.home__subtitle {
  font-family: var(--f-body);
  font-size: 17px;
  line-height: 1.6;
  color: var(--c-dim);
  margin-bottom: 40px;
  max-width: 560px;
}

html.lang-ja .home__subtitle {
  line-height: 1.9;
}

/* Catalog card */
.lib-card {
  background: var(--c-paper);
  border: 1px solid var(--c-rule);
  border-radius: 2px;
  box-shadow:
    0 1px 0 rgba(28, 27, 24, 0.04),
    0 12px 24px -18px rgba(28, 27, 24, 0.18);
}

.lib-card__header {
  padding: 14px 24px 10px;
  border-bottom: 1px solid var(--c-rule);
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-family: var(--f-mono);
  font-size: 10px;
  color: var(--c-dim);
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

html.lang-ja .lib-card__header {
  font-family: var(--f-body);
  font-size: 11px;
  letter-spacing: 0.24em;
  text-transform: none;
}

.lib-card__num {
  font-family: var(--f-mono);
  letter-spacing: 0.14em;
}

.lib-card__body {
  padding: 28px 36px 32px;
  display: flex;
  flex-direction: column;
  gap: 22px;
}

/* Field row (input + random btn) */
.lib-field-row {
  display: flex;
  align-items: flex-end;
  gap: 12px;
}

.lib-field {
  flex: 1;
  border-bottom: 1px solid var(--c-ink);
  padding-bottom: 8px;
}

.lib-field__label {
  font-family: var(--f-mono);
  font-size: 10px;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--c-dim);
  margin-bottom: 6px;
}

html.lang-ja .lib-field__label {
  font-family: var(--f-body);
  font-size: 11px;
  letter-spacing: 0.3em;
  text-transform: none;
}

.lib-field__input {
  width: 100%;
  border: none;
  background: none;
  font-family: var(--f-head);
  font-size: 26px;
  color: var(--c-ink);
  line-height: 1.1;
  outline: none;
  padding: 0;
}

.lib-field__input::placeholder {
  color: var(--c-dim);
  font-style: italic;
}

html.lang-ja .lib-field__input {
  font-size: 24px;
  font-weight: 500;
}

.lib-random-btn {
  all: unset;
  cursor: pointer;
  padding: 6px 10px;
  border: 1px solid var(--c-rule);
  border-radius: 2px;
  font-family: var(--f-mono);
  font-size: 10px;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: var(--c-dim);
  background: var(--c-bg);
  white-space: nowrap;
  flex-shrink: 0;
}

html.lang-ja .lib-random-btn {
  font-family: var(--f-body);
  font-size: 11px;
  letter-spacing: 0.2em;
  text-transform: none;
}

.lib-random-btn:hover {
  color: var(--c-ink);
  border-color: var(--c-ink);
}

/* Actions row */
.lib-card__actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 8px;
}

.lib-card__hint {
  font-family: var(--f-body);
  font-size: 11px;
  color: var(--c-dim);
}

html.lang-ja .lib-card__hint {
  font-family: var(--f-body);
  font-size: 12px;
}

.lib-card__hint-mono {
  font-family: var(--f-mono);
}

.lib-submit {
  all: unset;
  cursor: pointer;
  padding: 14px 28px;
  background: var(--c-ink);
  color: var(--c-paper);
  font-family: var(--f-mono);
  font-size: 12px;
  letter-spacing: 0.22em;
  text-transform: uppercase;
  font-weight: 500;
}

html.lang-ja .lib-submit {
  font-family: var(--f-head);
  font-size: 14px;
  letter-spacing: 0.32em;
  text-transform: none;
}

.lib-submit:hover {
  background: var(--c-accent);
}

/* Stats bar */
.home__stats {
  margin-top: 28px;
  font-family: var(--f-mono);
  font-size: 11px;
  color: var(--c-dim);
  display: flex;
  justify-content: center;
  gap: 16px;
  flex-wrap: wrap;
}

html.lang-ja .home__stats {
  font-family: var(--f-body);
}

.home__stats-dot {
  opacity: 0.4;
}

/* Mobile */
@media (max-width: 640px) {
  .home {
    padding: 28px 18px 32px;
  }

  .home__hero {
    font-size: 28px;
  }

  html.lang-ja .home__hero {
    font-size: 26px;
  }

  .home__subtitle {
    font-size: 15px;
    margin-bottom: 28px;
  }

  .lib-card__body {
    padding: 18px 18px 20px;
    gap: 18px;
  }

  .lib-card__actions {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .lib-submit {
    text-align: center;
    width: 100%;
  }

  .lib-card__header {
    padding: 12px 18px 10px;
  }

  .lib-field__input {
    font-size: 20px;
  }
}
</style>
