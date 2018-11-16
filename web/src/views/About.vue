<template>
  <div class="about">
    <header>
      <h1>PediaRouteについて</h1>
    </header>
    <article>
      <section>
        <h2>検索について</h2>
        <p>
          Wikipedia上には日本語版だけで数百万ものページがありますが、その中のどんな2つのページを選んでも最大6回のページ内のリンクをたどれば到達できる（らしい）というのが
          2011年8月ごろにネットで話題になりました。
          <a href="https://ja.wikipedia.org/wiki/%E5%85%AD%E6%AC%A1%E3%81%AE%E9%9A%94%E3%81%9F%E3%82%8A">六次の隔たり</a>のWikipedia版ですね。
        </p>
        <p>
          このサイトはWikipedia内の全ページのタイトルおよびページ間のリンクをオンメモリにもっており、スタートとゴールのページの双方向からリンクをたどれるか検索を行い、
          最短で見つかったリンクの経路を表示します。
        </p>
        <p>
          詳しいアルゴリズムに知りたい方は<a href="https://github.com/mtgto/pediaroute-go">ソースコード</a>を参照してください。
        </p>
      </section>
      <section>
        <h2>使用しているデータについて</h2>
        <p>
          このサイトでは
          <a
            href="https://ja.wikipedia.org/wiki/Wikipedia:%E3%83%87%E3%83%BC%E3%82%BF%E3%83%99%E3%83%BC%E3%82%B9%E3%83%80%E3%82%A6%E3%83%B3%E3%83%AD%E3%83%BC%E3%83%89"
          >Wikipediaのデータベース</a>
          からページ情報とページ間のリンク情報を取得し、利用しています。
        </p>
        <p>このサイトの情報は
          <a
            href="https://creativecommons.org/licenses/by-sa/3.0/deed.ja"
          >クリエイティブ・コモンズ 表示-継承 3.0 非移植ライセンス</a> の下で利用可能です。
        </p>
      </section>
      <aside id="history">
      <h2>更新履歴</h2>
      <dl>
        <dt>
          <time>2013/03/10 17:00</time>ごろ
        </dt>
        <dd>2013/03/07のデータで作り直し</dd>
        <dt>
          <time>2011/08/31 02:00</time>ごろ
        </dt>
        <dd>HTML書き直し。大文字小文字の違いしかないページ名は同じ扱いにした。</dd>
        <dt>
          <time>2011/08/30 04:30</time>ごろ
        </dt>
        <dd>ページが固まることがあったので改善。ランダムで選ぶときにひとつもリンクがないページは選ばないようにした。</dd>
        <dt>
          <time>2011/08/29 02:20</time>ごろ
        </dt>
        <dd>ランダムに２つのページを選択機能、検索結果のキャッシュ処理を追加、Twitterウィジェット追加、検索結果をつぶやくリンク追加</dd>
        <dt>
          <time>2011/08/28 02:21</time>ごろ
        </dt>
        <dd>公開</dd>
      </dl>
    </aside>
    </article>
    <aside>
      <h2>作者について</h2>
      <p>
        <a href="http://twitter.com/mtgto">@mtgto</a>
      </p>
    </aside>
  </div>
</template>
