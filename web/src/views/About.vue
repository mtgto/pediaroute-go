<template>
  <div class="about">
    <div class="about__inner">
      <p class="about__colophon">{{ t('about.colophon') }}</p>
      <h1 class="about__heading">{{ t('about.heading') }}</h1>
      <p class="about__lead">{{ t('about.lead') }}</p>

      <div class="about__sections">
        <!-- Section 01: Route finding -->
        <div class="about__section">
          <div class="about__section-meta">
            <div class="about__section-num">§ 01</div>
            <div class="about__section-label">{{ t('about.findingRoute') }}</div>
          </div>
          <div class="about__section-body">
            <p>
              {{ t('about.first') }}<a :href="t('about.sixDegreeUrl')">{{ t('about.sixDegree') }}</a
              >{{ t('about.second') }}
            </p>
            <p>{{ t('about.third') }}</p>
            <p>
              {{ t('about.fourth') }}<a href="https://github.com/mtgto/pediaroute-go">{{ t('about.sourceCode') }}</a
              >{{ t('about.fifth') }}
            </p>
          </div>
        </div>

        <!-- Section 02: Data -->
        <div class="about__section">
          <div class="about__section-meta">
            <div class="about__section-num">§ 02</div>
            <div class="about__section-label">{{ t('about.data') }}</div>
          </div>
          <div class="about__section-body">
            <p>
              {{ t('about.sixth') }}<a :href="t('about.databaseUrl')">{{ t('about.database') }}</a
              >{{ t('about.seventh') }}
            </p>
            <p>
              {{ t('about.eighth') }}<a :href="t('about.commonsUrl')">{{ t('about.commons') }}</a
              >{{ t('about.ninth') }}
            </p>
          </div>
        </div>

        <!-- Section 03: Author -->
        <div class="about__section">
          <div class="about__section-meta">
            <div class="about__section-num">§ 03</div>
            <div class="about__section-label">{{ t('about.author') }}</div>
          </div>
          <div class="about__section-body">
            <p>
              <i18n-t keypath="about.authorBio">
                <template #handle><a href="https://twitter.com/mtgto">@mtgto</a></template>
              </i18n-t>
            </p>
            <p class="about__meta">
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

.about__inner {
  max-width: 760px;
  margin: 0 auto;
}

.about__colophon {
  font-family: var(--f-mono);
  font-size: 11px;
  letter-spacing: 0.22em;
  text-transform: uppercase;
  color: var(--c-accent);
  margin-bottom: 18px;
}

html.lang-ja .about__colophon {
  font-family: var(--f-body);
  font-size: 12px;
  letter-spacing: 0.4em;
  text-transform: none;
}

.about__heading {
  font-family: var(--f-head);
  font-weight: 400;
  font-size: 56px;
  line-height: 1;
  margin-bottom: 8px;
  letter-spacing: -0.015em;
}

html.lang-ja .about__heading {
  font-weight: 600;
  font-size: 54px;
  line-height: 1.15;
  letter-spacing: 0.04em;
}

.about__lead {
  font-family: var(--f-body);
  font-style: italic;
  font-size: 18px;
  color: var(--c-dim);
  margin-bottom: 48px;
}

html.lang-ja .about__lead {
  font-style: normal;
  font-size: 17px;
  line-height: 1.9;
}

/* Sections */
.about__sections {
  border-top: 1px solid var(--c-rule);
}

.about__section {
  display: grid;
  grid-template-columns: 120px 1fr;
  gap: 32px;
  padding: 28px 0;
  border-bottom: 1px solid var(--c-rule);
}

html.lang-ja .about__section {
  grid-template-columns: 140px 1fr;
}

.about__section-meta {
  flex-shrink: 0;
}

.about__section-num {
  font-family: var(--f-mono);
  font-size: 11px;
  color: var(--c-accent);
  letter-spacing: 0.16em;
  margin-bottom: 4px;
}

.about__section-label {
  font-family: var(--f-body);
  font-style: italic;
  font-size: 18px;
  color: var(--c-ink);
}

html.lang-ja .about__section-label {
  font-family: var(--f-head);
  font-style: normal;
  font-size: 17px;
  font-weight: 600;
  letter-spacing: 0.04em;
}

.about__section-body {
  font-family: var(--f-body);
  font-size: 17px;
  line-height: 1.6;
  color: var(--c-ink);
  display: flex;
  flex-direction: column;
  gap: 14px;
}

html.lang-ja .about__section-body {
  font-size: 16px;
  line-height: 1.95;
}

.about__meta {
  font-family: var(--f-mono);
  font-size: 12px;
  color: var(--c-dim);
}

/* Mobile */
@media (max-width: 640px) {
  .about {
    padding: 28px 18px 40px;
  }

  .about__heading {
    font-size: 36px;
  }

  html.lang-ja .about__heading {
    font-size: 32px;
  }

  .about__section {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .about__section-meta {
    display: flex;
    align-items: baseline;
    gap: 12px;
  }
}
</style>
