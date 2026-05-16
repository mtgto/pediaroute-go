<template>
  <div class="field">
    <div class="label">{{ label }}</div>
    <input class="input" v-bind="$attrs" />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  label: string;
}>();

defineOptions({ inheritAttrs: false });
</script>

<style scoped>
.field {
  flex: 1;
  border-bottom: 1px solid var(--c-ink);
  padding-bottom: 8px;
}

.label {
  font-family: var(--ui-font);
  font-size: 10px;
  letter-spacing: 0.18em;
  text-transform: var(--ui-tt);
  color: var(--c-dim);
  margin-bottom: 6px;
}

html.lang-ja .label {
  font-size: 11px;
  letter-spacing: 0.3em;
}

.input {
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

.input::placeholder {
  color: var(--c-dim);
  font-style: italic;
}

html.lang-ja .input {
  font-size: 24px;
  font-weight: 500;
}

@media (max-width: 640px) {
  .input {
    font-size: 20px;
  }
}
</style>
