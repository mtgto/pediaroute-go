<template>
  <div class="app">
    <header class="header">
      <RouterLink to="/" class="logo">
        <div class="logo__word">
          <span class="logo__accent">P</span>edia<span class="logo__accent logo__r">R</span>oute<span class="logo__com">.com</span>
        </div>
      </RouterLink>
      <div class="nav">
        <span class="nav__est">{{ t('nav.est') }}</span>
        <span class="nav__dot">·</span>
        <button :class="['lang', !isJa && 'lang--active']" @click="locale = 'en'">EN</button>
        <button :class="['lang', 'lang--jp', isJa && 'lang--active']" @click="locale = 'ja'">日本語</button>
      </div>
    </header>

    <main>
      <RouterView />
    </main>

    <footer class="footer">
      <div class="footer__nav">
        <RouterLink to="/">{{ t('nav.home') }}</RouterLink>
        <span class="sep"> · </span>
        <RouterLink to="/about">{{ t('nav.about') }}</RouterLink>
        <span class="sep"> · </span>
        <a href="https://github.com/mtgto/pediaroute-go" target="_blank">{{ t('nav.source') }}</a>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue';
import { RouterLink, RouterView } from 'vue-router';
import { useI18n } from 'vue-i18n';

const { t, locale } = useI18n();
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
.header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  border-bottom: 1px solid var(--c-rule);
  padding: 36px 64px 22px;
}

.logo {
  text-decoration: none;
  color: inherit;
  flex-shrink: 0;
}

.logo:hover {
  text-decoration: none;
}

.logo__word {
  font-family: var(--f-serif);
  font-weight: 700;
  font-size: 44px;
  letter-spacing: -0.01em;
  line-height: 1;
  display: flex;
  align-items: baseline;
}

.logo__accent {
  color: var(--c-accent);
}

.logo__r {
  margin-left: 0.06em;
}

.logo__com {
  font-family: var(--f-mono);
  font-weight: 400;
  font-size: 14px;
  color: var(--c-dim);
  margin-left: 10px;
  letter-spacing: 0.02em;
}

.logo__sub {
  font-family: 'Shippori Mincho', 'Hiragino Mincho ProN', serif;
  font-size: 12px;
  color: var(--c-dim);
  margin-top: 6px;
  letter-spacing: 0.2em;
}

.nav {
  font-family: var(--f-mono);
  font-size: 11px;
  color: var(--c-dim);
  display: flex;
  gap: 18px;
  align-items: center;
}

.nav__est {
  text-transform: uppercase;
  letter-spacing: 0.16em;
}

.nav__dot {
  opacity: 0.4;
}

.lang {
  all: unset;
  cursor: pointer;
  font-family: var(--f-mono);
  font-size: 11px;
  color: var(--c-dim);
  border-bottom: 1px solid transparent;
  padding-bottom: 1px;
}

.lang--jp {
  font-family: 'Shippori Mincho', 'Hiragino Mincho ProN', serif;
}

.lang--active {
  font-weight: 600;
  color: var(--c-ink);
  border-bottom-color: var(--c-ink);
}

/* ---- Main ---- */
main {
  flex: 1;
}

/* ---- Footer ---- */
.footer {
  font-family: var(--f-mono);
  font-size: 11px;
  color: var(--c-dim);
  padding: 24px 64px 32px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px solid var(--c-rule);
}

.footer__nav {
  display: flex;
  align-items: center;
}

.footer__nav a {
  color: var(--c-dim);
  font-family: var(--f-mono);
  font-size: 11px;
}

.footer__nav a:hover {
  color: var(--c-ink);
  text-decoration: none;
}

.sep {
  opacity: 0.4;
  padding: 0 4px;
}

/* ---- Mobile ---- */
@media (max-width: 640px) {
  .header {
    padding: 20px 18px 14px;
    align-items: center;
  }

  .logo__word {
    font-size: 22px;
  }

  .logo__com {
    font-size: 10px;
    margin-left: 6px;
  }

  .logo__sub {
    display: none;
  }

  .nav__est,
  .nav__dot {
    display: none;
  }

  .nav {
    gap: 12px;
  }

  .footer {
    padding: 20px 18px;
    flex-direction: column;
    gap: 8px;
    align-items: flex-start;
  }
}
</style>
