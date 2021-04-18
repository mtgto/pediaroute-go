<template>
  <div class="about">
    <header>
      <h1>{{ $t('about.header') }}</h1>
    </header>
    <article>
      <section>
        <h2>{{ $t('about.findingRoute') }}</h2>
        <p>
          {{ $t('about.first') }}
          <a :href="$t('about.sixDegreeUrl')">{{ $t('about.sixDegree') }}</a>
          {{ $t('about.second') }}
        </p>
        <p>{{ $t('about.third') }}</p>
        <p>
          {{ $t('about.fourth') }}
          <a href="https://github.com/mtgto/pediaroute-go">{{ $t('about.sourceCode') }}</a>
          {{ $t('about.fifth') }}
        </p>
      </section>
      <section>
        <h2>{{ $t('about.data') }}</h2>
        <p>
          {{ $t('about.sixth') }}
          <a :href="$t('about.databaseUrl')">{{ $t('about.database') }}</a>
          {{ $t('about.seventh') }}
        </p>
        <p>
          {{ $t('about.eighth') }}
          <a :href="$t('about.commonsUrl')">{{ $t('about.commons') }}</a>
          {{ $t('about.ninth') }}
        </p>
      </section>
    </article>
    <aside>
      <h2>{{ $t('about.author') }}</h2>
      <p>
        <a href="http://twitter.com/mtgto">@mtgto</a>
      </p>
    </aside>
  </div>
</template>
