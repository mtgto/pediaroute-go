<template>
  <div class="field">
    <div class="label">{{ label }}</div>
    <input class="input" v-bind="$attrs" v-model="model" @focus="onFocus" @blur="onBlur" @keydown="onKeydown" />
    <ul v-if="showDropdown" class="dropdown" role="listbox">
      <li
        v-for="(word, i) in suggestions"
        :key="word"
        class="dropdown-item"
        :class="{ 'dropdown-item--active': i === activeIndex }"
        role="option"
        @mousedown.prevent
        @click="onSelect(word)"
      >
        {{ word }}
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

const model = defineModel<string>({ required: true });

const props = defineProps<{
  label: string;
  suggestions?: string[];
}>();

defineOptions({ inheritAttrs: false });

const isFocused = ref(false);
const activeIndex = ref(-1);

watch(
  () => props.suggestions,
  () => {
    activeIndex.value = -1;
  },
);

const showDropdown = computed(() => {
  if (!isFocused.value || !props.suggestions?.length) return false;
  return !props.suggestions.some((s) => s.toLowerCase() === model.value.toLowerCase());
});

const onFocus = () => {
  isFocused.value = true;
};
const onBlur = () => {
  isFocused.value = false;
  activeIndex.value = -1;
};

const onKeydown = (e: KeyboardEvent) => {
  if (!showDropdown.value) return;
  const list = props.suggestions!;
  if (e.key === 'ArrowDown') {
    e.preventDefault();
    activeIndex.value = (activeIndex.value + 1) % list.length;
  } else if (e.key === 'ArrowUp') {
    e.preventDefault();
    activeIndex.value = (activeIndex.value - 1 + list.length) % list.length;
  } else if (e.key === 'Enter' && activeIndex.value >= 0) {
    e.preventDefault();
    model.value = list[activeIndex.value]!;
    isFocused.value = false;
  } else if (e.key === 'Escape') {
    isFocused.value = false;
  }
};

const onSelect = (word: string) => {
  model.value = word;
  isFocused.value = false;
};
</script>

<style scoped>
.field {
  flex: 1;
  border-bottom: 1px solid var(--c-ink);
  padding-bottom: 8px;
  position: relative;
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

.dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  background: var(--c-paper);
  border: 1px solid var(--c-rule);
  border-top: none;
  list-style: none;
  z-index: 20;
  max-height: 280px;
  overflow-y: auto;
  box-shadow: 0 6px 16px -8px rgba(28, 27, 24, 0.22);
}

.dropdown-item {
  padding: 8px 12px;
  font-family: var(--f-head);
  font-size: 20px;
  color: var(--c-ink);
  cursor: pointer;
  line-height: 1.2;
  border-bottom: 1px solid var(--c-rule);
}

.dropdown-item:last-child {
  border-bottom: none;
}

.dropdown-item:hover,
.dropdown-item--active {
  background: var(--c-accent);
  color: var(--c-paper);
}

html.lang-ja .dropdown-item {
  font-size: 18px;
  font-weight: 500;
}

@media (max-width: 640px) {
  .input {
    font-size: 20px;
  }

  .dropdown-item {
    font-size: 16px;
  }
}
</style>
