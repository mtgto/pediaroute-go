<template>
  <div class="card">
    <div class="header">
      <slot name="header-title" />
      <span class="num">№ 0001</span>
    </div>
    <div class="body">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.card {
  background: var(--c-paper);
  border: 1px solid var(--c-rule);
  border-radius: 2px;
  box-shadow:
    0 1px 0 rgba(28, 27, 24, 0.04),
    0 12px 24px -18px rgba(28, 27, 24, 0.18);
}

.header {
  padding: 14px 24px 10px;
  border-bottom: 1px solid var(--c-rule);
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-family: var(--ui-font);
  font-size: 10px;
  color: var(--c-dim);
  letter-spacing: 0.14em;
  text-transform: var(--ui-tt);
}

html.lang-ja .header {
  font-size: 11px;
  letter-spacing: 0.24em;
}

.num {
  font-family: var(--f-mono);
  letter-spacing: 0.14em;
}

.body {
  padding: 28px 36px 32px;
  display: flex;
  flex-direction: column;
  gap: 22px;
}

@media (max-width: 640px) {
  .header {
    padding: 12px 18px 10px;
  }

  .body {
    padding: 18px 18px 20px;
    gap: 18px;
  }
}
</style>
