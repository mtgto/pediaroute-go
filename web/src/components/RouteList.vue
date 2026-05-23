<template>
  <div>
    <div class="list">
      <template v-for="(word, i) in routes" :key="word">
        <div :class="['slip', i === routes.length - 1 && 'dest']">
          <div class="step">
            <template v-if="i === 0">{{ t('search.origin') }}</template>
            <template v-else-if="i === routes.length - 1">{{ t('search.destination') }}</template>
            <template v-else>
              <i18n-t keypath="search.step">
                <template #num
                  ><span class="step-n">{{ i }}</span></template
                >
              </i18n-t>
            </template>
          </div>
          <div class="title">
            <a :href="t('message.wikipediaUrl', { word })" target="_blank" rel="noopener">{{ word }}<span class="ext">↗</span></a>
          </div>
          <div class="num">{{ i + 1 }}</div>
        </div>
        <div v-if="i < routes.length - 1" class="connector">↓</div>
      </template>
    </div>

    <div class="actions">
      <Button :as="RouterLink" variant="primary" :to="{ path: '/search', query: { lang: locale, wordFrom: wordTo, wordTo: wordFrom } }">
        {{ t('search.reverseRoute') }}
      </Button>
      <Button :as="RouterLink" variant="outline" to="/">
        {{ t('search.newSearch') }}
      </Button>
      <Button as="a" variant="ghost" :href="shareUrl" target="_blank" rel="noopener">
        {{ t('message.share') }}
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router';
import { useI18n } from 'vue-i18n';
import Button from './Button.vue';

defineProps<{
  routes: readonly string[];
  wordFrom: string;
  wordTo: string;
  shareUrl: string;
}>();

const { t, locale } = useI18n();
</script>

<style scoped>
/* Route list */
.list {
  display: flex;
  flex-direction: column;
}

.slip {
  display: grid;
  grid-template-columns: 72px 1fr auto;
  align-items: center;
  gap: 24px;
  padding: 18px 28px;
  background: var(--c-paper);
  border: 1px solid var(--c-rule);
  border-radius: 2px;
}

.slip.dest {
  background: var(--c-accent);
  border: 1px solid var(--c-accent);
  border-left: 6px solid var(--c-ink);
  padding: 20px 28px 20px 22px;
  color: var(--c-paper);
}

.slip.dest .step {
  color: rgba(250, 245, 230, 0.85);
  font-weight: 600;
  border-right: 1px solid rgba(250, 245, 230, 0.18);
}

.slip.dest .title {
  font-size: 26px;
  font-weight: 700;
}

.slip.dest .title a {
  color: var(--c-paper);
}

.slip.dest .title a:hover {
  color: rgba(250, 245, 230, 0.75);
  text-decoration: none;
}

.slip.dest .ext {
  color: rgba(250, 245, 230, 0.55);
}

.slip.dest .num {
  color: rgba(250, 245, 230, 0.85);
  opacity: 0.95;
  font-size: 34px;
  font-weight: 700;
}

html.lang-ja .slip {
  grid-template-columns: 88px 1fr auto;
}

.step {
  font-family: var(--f-mono);
  font-size: 11px;
  color: var(--c-dim);
  letter-spacing: 0.1em;
  text-transform: uppercase;
  border-right: 1px solid var(--c-rule);
  padding-right: 16px;
  line-height: 1.4;
}

html.lang-ja .step {
  font-family: var(--f-body);
  font-size: 12px;
  letter-spacing: 0.18em;
  text-transform: none;
}

.step-n {
  font-family: var(--f-head);
  font-style: italic;
}

.title {
  font-family: var(--f-head);
  font-size: 22px;
  line-height: 1.2;
  font-weight: 500;
}

.title a {
  color: var(--c-ink);
}

.title a:hover {
  color: var(--c-accent);
  text-decoration: none;
}

.ext {
  margin-left: 6px;
  font-family: var(--f-mono);
  font-size: 11px;
  color: var(--c-dim);
  font-weight: 400;
}

.num {
  font-family: var(--f-head);
  font-style: italic;
  font-size: 32px;
  color: var(--c-accent);
  opacity: 0.35;
}

.connector {
  height: 16px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-family: var(--f-mono);
  color: var(--c-accent);
  font-size: 12px;
}

/* Actions */
.actions {
  margin-top: 36px;
  display: flex;
  gap: 14px;
  flex-wrap: wrap;
}

/* Mobile */
@media (max-width: 640px) {
  .slip {
    grid-template-columns: 56px 1fr auto;
    gap: 14px;
    padding: 14px 16px;
  }

  .slip.dest {
    padding: 14px 16px 14px 10px;
  }

  html.lang-ja .slip {
    grid-template-columns: 72px 1fr auto;
  }

  .title {
    font-size: 17px;
  }

  .slip.dest .title {
    font-size: 17px;
  }

  .num {
    font-size: 22px;
  }

  .slip.dest .num {
    font-size: 22px;
  }

  .actions {
    flex-direction: column;
  }
}
</style>
