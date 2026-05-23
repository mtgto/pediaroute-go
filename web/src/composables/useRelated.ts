import { ref, onUnmounted } from 'vue';
import type { Ref } from 'vue';

const cache = new Map<string, string[]>();

export function useRelated(locale: Ref<string>) {
  const suggestions = ref<string[] | undefined>(undefined);
  let timer: ReturnType<typeof setTimeout> | null = null;

  const fetch = (word: string) => {
    if (timer) clearTimeout(timer);
    if ([...word].length < 2) {
      suggestions.value = undefined;
      return;
    }
    suggestions.value = undefined;
    const key = `${locale.value}:${word}`;
    const cached = cache.get(key);
    if (cached !== undefined) {
      suggestions.value = cached;
      return;
    }
    timer = setTimeout(() => {
      globalThis
        .fetch(`/api/related?lang=${encodeURIComponent(locale.value)}&title=${encodeURIComponent(word)}`)
        .then((r) => r.json())
        .then((words: string[]) => {
          cache.set(key, words);
          suggestions.value = words;
        })
        .catch(console.error);
    }, 200);
  };

  const clear = () => {
    suggestions.value = undefined;
  };

  onUnmounted(() => {
    if (timer) {
      clearTimeout(timer);
    }
  });

  return { suggestions, fetch, clear };
}
