<template>
  <div class="about">
    <div class="inner">
      <p class="colophon">{{ t('about.colophon') }}</p>
      <h1 class="heading">{{ t('about.heading') }}</h1>
      <p class="lead">{{ t('about.lead') }}</p>

      <div class="sections">
        <!-- Section 01: Route finding -->
        <div class="section">
          <div class="section-meta">
            <div class="section-num">§ 01</div>
            <div class="section-label">{{ t('about.findingRoute') }}</div>
          </div>
          <div class="section-body">
            <p>
              {{ t('about.sixDegreeIntro') }} <a :href="t('about.sixDegreeUrl')">{{ t('about.sixDegree') }}</a
              >{{ t('about.sixDegreeSuffix') }}
            </p>
            <p>{{ t('about.routeDescription') }}</p>
          </div>
        </div>

        <!-- Section 02: Data -->
        <div class="section">
          <div class="section-meta">
            <div class="section-num">§ 02</div>
            <div class="section-label">{{ t('about.data') }}</div>
          </div>
          <div class="section-body">
            <p>
              {{ t('about.databasePrefix') }}<a :href="t('about.databaseUrl')">{{ t('about.database') }}</a
              >{{ t('about.databaseSuffix') }}
            </p>
            <p>
              {{ t('about.licensePrefix') }}<a :href="t('about.commonsUrl')">{{ t('about.commons') }}</a
              >{{ t('about.licenseSuffix') }}
            </p>
          </div>
        </div>

        <!-- Section 03: Author -->
        <div class="section">
          <div class="section-meta">
            <div class="section-num">§ 03</div>
            <div class="section-label">{{ t('about.author') }}</div>
          </div>
          <div class="section-body">
            <p>
              <i18n-t keypath="about.authorBio">
                <template #handle><a href="https://x.com/mtgto">@mtgto</a></template>
              </i18n-t>
            </p>
            <p class="meta">
              <a href="https://github.com/mtgto/pediaroute-go">{{ t('about.sourceCode') }}</a>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>

<style scoped>
.about {
  padding: 56px 64px 40px;
}

.inner {
  max-width: 760px;
  margin: 0 auto;
}

.colophon {
  font-family: var(--f-mono);
  font-size: 11px;
  letter-spacing: 0.22em;
  text-transform: uppercase;
  color: var(--c-accent);
  margin-bottom: 18px;
}

html.lang-ja .colophon {
  font-family: var(--f-body);
  font-size: 12px;
  letter-spacing: 0.4em;
  text-transform: none;
}

.heading {
  font-family: var(--f-head);
  font-weight: 400;
  font-size: 56px;
  line-height: 1;
  margin-bottom: 8px;
  letter-spacing: -0.015em;
}

html.lang-ja .heading {
  font-weight: 600;
  font-size: 54px;
  line-height: 1.15;
  letter-spacing: 0.04em;
}

.lead {
  font-family: var(--f-body);
  font-style: italic;
  font-size: 18px;
  color: var(--c-dim);
  margin-bottom: 48px;
}

html.lang-ja .lead {
  font-style: normal;
  font-size: 17px;
  line-height: 1.9;
}

/* Sections */
.sections {
  border-top: 1px solid var(--c-rule);
}

.section {
  display: grid;
  grid-template-columns: 120px 1fr;
  gap: 32px;
  padding: 28px 0;
  border-bottom: 1px solid var(--c-rule);
}

html.lang-ja .section {
  grid-template-columns: 140px 1fr;
}

.section-meta {
  flex-shrink: 0;
}

.section-num {
  font-family: var(--f-mono);
  font-size: 11px;
  color: var(--c-accent);
  letter-spacing: 0.16em;
  margin-bottom: 4px;
}

.section-label {
  font-family: var(--f-body);
  font-style: italic;
  font-size: 18px;
  color: var(--c-ink);
}

html.lang-ja .section-label {
  font-family: var(--f-head);
  font-style: normal;
  font-size: 17px;
  font-weight: 600;
  letter-spacing: 0.04em;
}

.section-body {
  font-family: var(--f-body);
  font-size: 17px;
  line-height: 1.6;
  color: var(--c-ink);
  display: flex;
  flex-direction: column;
  gap: 14px;
}

html.lang-ja .section-body {
  font-size: 16px;
  line-height: 1.95;
}

.meta {
  font-family: var(--f-mono);
  font-size: 12px;
  color: var(--c-dim);
}

/* Mobile */
@media (max-width: 640px) {
  .about {
    padding: 28px 18px 40px;
  }

  .heading {
    font-size: 36px;
  }

  html.lang-ja .heading {
    font-size: 32px;
  }

  .section {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .section-meta {
    display: flex;
    align-items: baseline;
    gap: 12px;
  }
}
</style>
