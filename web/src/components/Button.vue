<template>
  <component :is="as" :class="['btn', variant && `btn--${variant}`]" v-bind="$attrs">
    <slot />
  </component>
</template>

<script setup lang="ts">
import type { Component } from 'vue';

withDefaults(
  defineProps<{
    as?: string | Component;
    variant?: 'primary' | 'outline' | 'ghost';
  }>(),
  { as: 'button' },
);

defineOptions({ inheritAttrs: false });
</script>

<style scoped>
.btn {
  display: inline-block;
  padding: 12px 20px;
  font-family: var(--cta-font);
  font-size: 11px;
  letter-spacing: 0.18em;
  text-transform: var(--ui-tt);
  cursor: pointer;
  text-decoration: none;
  border: none;
  background: none;
}

html.lang-ja .btn {
  font-size: 13px;
  letter-spacing: 0.32em;
}

.btn--primary {
  background: var(--c-ink);
  color: var(--c-paper);
}

.btn--primary:hover {
  background: var(--c-accent);
  color: var(--c-paper);
  text-decoration: none;
}

.btn--outline {
  border: 1px solid var(--c-ink);
  color: var(--c-ink);
}

.btn--outline:hover {
  border-color: var(--c-accent);
  color: var(--c-accent);
  text-decoration: none;
}

.btn--ghost {
  color: var(--c-dim);
  border-bottom: 1px solid var(--c-rule);
}

.btn--ghost:hover {
  color: var(--c-ink);
  text-decoration: none;
}

@media (max-width: 640px) {
  .btn {
    text-align: center;
  }
}
</style>
