<template>
  <div class="notice">
    <div class="header">
      <slot name="header-title" />
      <span class="num">{{ num }}</span>
    </div>
    <p :class="['body', bodyLarge && 'body--large']">
      <slot name="body" />
    </p>
    <p v-if="$slots.note" class="note">
      <slot name="note" />
    </p>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    num?: string;
    bodyLarge?: boolean;
  }>(),
  { num: '№ 0001 · v', bodyLarge: false },
);
</script>

<style scoped>
.notice {
  background: var(--c-paper);
  border: 1px solid var(--c-rule);
  padding: 28px 36px;
  border-radius: 2px;
  margin-bottom: 36px;
  box-shadow:
    0 1px 0 rgba(28, 27, 24, 0.04),
    0 12px 24px -18px rgba(28, 27, 24, 0.18);
}

.header {
  font-family: var(--ui-font);
  font-size: 10px;
  letter-spacing: 0.18em;
  text-transform: var(--ui-tt);
  color: var(--c-dim);
  margin-bottom: 16px;
  display: flex;
  justify-content: space-between;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--c-rule);
}

html.lang-ja .header {
  font-size: 11px;
  letter-spacing: 0.32em;
}

.num {
  font-family: var(--f-mono);
  letter-spacing: 0.12em;
}

.body {
  font-family: var(--f-body);
  font-size: 17px;
  line-height: 1.6;
  color: var(--c-ink);
}

.body--large {
  font-size: 21px;
  line-height: 1.45;
  margin-bottom: 18px;
}

html.lang-ja .body--large {
  font-family: var(--f-head);
  font-size: 22px;
  line-height: 1.7;
  font-weight: 500;
}

.note {
  font-family: var(--f-body);
  font-size: 16px;
  line-height: 1.65;
  color: var(--c-dim);
  max-width: 640px;
}

html.lang-ja .note {
  font-size: 14px;
  line-height: 1.9;
}

@media (max-width: 640px) {
  .notice {
    padding: 20px 18px;
  }
}
</style>
