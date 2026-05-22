<template>
  <div :class="['stamp', found ? 'found' : 'notfound']">
    <div class="word">
      {{ found ? t('stamp.found') : t('stamp.notFound') }}
    </div>
    <div class="meta">{{ meta }}</div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';

defineProps<{
  found: boolean;
  meta: string;
}>();

const { t } = useI18n();
</script>

<style scoped>
.stamp {
  border: 1.5px solid var(--c-accent);
  padding: 8px 14px;
  color: var(--c-accent);
  font-family: var(--ui-font);
  font-size: 11px;
  letter-spacing: 0.18em;
  text-transform: var(--ui-tt);
  border-radius: 3px;
  background: rgba(138, 51, 36, 0.04);
  text-align: center;
  flex-shrink: 0;
}

html.lang-ja .stamp {
  font-size: 14px;
  letter-spacing: 0.4em;
}

.stamp.found {
  transform: rotate(-2deg);
}

.stamp.notfound {
  transform: rotate(-3deg);
}

.word {
  font-weight: 600;
  font-size: 13px;
  letter-spacing: 0.16em;
}

html.lang-ja .word {
  font-size: 14px;
  letter-spacing: 0.5em;
}

.meta {
  font-size: 9px;
  opacity: 0.8;
  margin-top: 1px;
  font-family: var(--f-mono);
  letter-spacing: 0.08em;
}
</style>
