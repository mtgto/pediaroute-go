<template>
  <div class="app">
    <header class="lib-header">
      <RouterLink to="/" class="lib-logo">
        <div class="lib-logo__word">
          <span class="lib-logo__accent">P</span>edia<span class="lib-logo__accent lib-logo__r">R</span>oute<span class="lib-logo__com">.com</span>
        </div>
        <div v-if="isJa" class="lib-logo__sub">ペディアルート</div>
      </RouterLink>
      <div class="lib-nav">
        <span class="lib-nav__est">{{ isJa ? '2013年創設' : 'Est. 2013' }}</span>
        <span class="lib-nav__dot">·</span>
        <button :class="['lib-lang', !isJa && 'lib-lang--active']" @click="locale = 'en'">EN</button>
        <button :class="['lib-lang', 'lib-lang--jp', isJa && 'lib-lang--active']" @click="locale = 'ja'">日本語</button>
      </div>
    </header>

    <main>
      <RouterView />
    </main>

    <footer class="lib-footer">
      <div class="lib-footer__nav">
        <RouterLink to="/">{{ isJa ? 'トップ' : 'Home' }}</RouterLink>
        <span class="lib-sep"> · </span>
        <RouterLink to="/about">{{ isJa ? 'このサイトについて' : 'About' }}</RouterLink>
        <span class="lib-sep"> · </span>
        <a href="https://github.com/mtgto/pediaroute-go" target="_blank">{{ isJa ? 'ソースコード' : 'Source' }}</a>
      </div>
      <span class="lib-footer__credit">Wikipedia · CC BY-SA 3.0</span>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue';
import { RouterLink, RouterView } from 'vue-router';
import { useI18n } from 'vue-i18n';

const { locale } = useI18n();
const isJa = computed(() => locale.value === 'ja');

watch(
  locale,
  (lang) => {
    document.documentElement.setAttribute('lang', lang);
    document.documentElement.className = `lang-${lang}`;
  },
  { immediate: true },
);
</script>

<style scoped>
.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

/* ---- Header ---- */
.lib-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  border-bottom: 1px solid var(--c-rule);
  padding: 36px 64px 22px;
}

.lib-logo {
  text-decoration: none;
  color: inherit;
  flex-shrink: 0;
}

.lib-logo:hover {
  text-decoration: none;
}

.lib-logo__word {
  font-family: var(--f-serif);
  font-weight: 700;
  font-size: 44px;
  letter-spacing: -0.01em;
  line-height: 1;
  display: flex;
  align-items: baseline;
}

.lib-logo__accent {
  color: var(--c-accent);
}

.lib-logo__r {
  margin-left: 0.06em;
}

.lib-logo__com {
  font-family: var(--f-mono);
  font-weight: 400;
  font-size: 14px;
  color: var(--c-dim);
  margin-left: 10px;
  letter-spacing: 0.02em;
}

.lib-logo__sub {
  font-family: 'Shippori Mincho', 'Hiragino Mincho ProN', serif;
  font-size: 12px;
  color: var(--c-dim);
  margin-top: 6px;
  letter-spacing: 0.2em;
}

.lib-nav {
  font-family: var(--f-mono);
  font-size: 11px;
  color: var(--c-dim);
  display: flex;
  gap: 18px;
  align-items: center;
}

.lib-nav__est {
  text-transform: uppercase;
  letter-spacing: 0.16em;
}

.lib-nav__dot {
  opacity: 0.4;
}

.lib-lang {
  all: unset;
  cursor: pointer;
  font-family: var(--f-mono);
  font-size: 11px;
  color: var(--c-dim);
  border-bottom: 1px solid transparent;
  padding-bottom: 1px;
}

.lib-lang--jp {
  font-family: 'Shippori Mincho', 'Hiragino Mincho ProN', serif;
}

.lib-lang--active {
  font-weight: 600;
  color: var(--c-ink);
  border-bottom-color: var(--c-ink);
}

/* ---- Main ---- */
main {
  flex: 1;
}

/* ---- Footer ---- */
.lib-footer {
  font-family: var(--f-mono);
  font-size: 11px;
  color: var(--c-dim);
  padding: 24px 64px 32px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px solid var(--c-rule);
}

.lib-footer__nav {
  display: flex;
  align-items: center;
}

.lib-footer__nav a {
  color: var(--c-dim);
  font-family: var(--f-mono);
  font-size: 11px;
}

.lib-footer__nav a:hover {
  color: var(--c-ink);
  text-decoration: none;
}

.lib-sep {
  opacity: 0.4;
  padding: 0 4px;
}

.lib-footer__credit {
  opacity: 0.7;
}

/* ---- Mobile ---- */
@media (max-width: 640px) {
  .lib-header {
    padding: 20px 18px 14px;
    align-items: center;
  }

  .lib-logo__word {
    font-size: 22px;
  }

  .lib-logo__com {
    font-size: 10px;
    margin-left: 6px;
  }

  .lib-logo__sub {
    display: none;
  }

  .lib-nav__est,
  .lib-nav__dot {
    display: none;
  }

  .lib-nav {
    gap: 12px;
  }

  .lib-footer {
    padding: 20px 18px;
    flex-direction: column;
    gap: 8px;
    align-items: flex-start;
  }
}
</style>
