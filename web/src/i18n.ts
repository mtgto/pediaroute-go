import Vue from 'vue';
import VueI18n from 'vue-i18n';

Vue.use(VueI18n);

const detectUserLanguage = (fallbackLanguage: string = 'en'): string => {
  const language: string = navigator.language;
  if (language) {
    return language.split('-')[0];
  }
  return fallbackLanguage;
};

export const i18n: VueI18n = new VueI18n({
  locale: detectUserLanguage(),
  fallbackLocale: 'en',
  messages: {
    en: {
      message: {
        header: 'Find a route between any page in Wikipedia within 6 links',
        search: 'Search',
        searchFrom: 'to',
        searchTo: '',
        buttonRandom: 'Select a random page',
        searching: 'Searching route from "{wordFrom}" to "{wordTo}"',
        searchResult: 'Search result of route from "{wordFrom}" to "{wordTo}" ({second} sec)',
        searchInReverse: 'Search a route from "{wordTo}" to "{wordFrom}"',
        tweet: 'Tweet (open new window)',
        tweetFind: 'https://twitter.com/home?status=From {wordFrom} to {wordTo} can be reached in {length} links in Wikipedia. {link} {hashTag}',
        tweetNotFound: 'https://twitter.com/home?status=From {wordFrom} to {wordTo} can not be reached within 6 links in Wikipedia. {link} {hashTag}',
        wikipediaUrl: 'https://en.wikipedia.org/wiki/{word}',
      },
      error: {
        notFoundFrom: 'The page "{from}" does not exists.',
        notFoundTo: 'The page "{to}" does not exists.',
        notFoundRoute: "I can't find a route within 6 links...",
      },
      about: {
        header: 'About PediaRoute',
        sixDegree: 'Six degrees of separation',
        sixDegreeUrl: 'https://en.wikipedia.org/wiki/Six_degrees_of_separation',
        first: 'There are millions of pages on Wikipedia. ' + 'Do you hear about it is possible to arrive from any page to any page within 6 clicks?',
        second: ' of Wikipedia.',
        third:
          'This site uses all titles and links between any two pages in Wikipedia and searches whether the shortest route can be traced from both the start and goal pages.',
        fourth: 'For detail, see ',
        sourceCode: 'Source code',
        fifth: '.',
        sixth: 'This site uses ',
        database: 'Database of Wikipedia',
        databaseUrl: 'https://en.wikipedia.org/wiki/Wikipedia:Database_download',
        seventh: '.',
        eighth: 'You can use the data of this site under ',
        commons: 'CC BY-SA 3.0',
        commonsUrl: 'https://creativecommons.org/licenses/by-sa/3.0/deed',
        nineth: '.',
        findingRoute: 'Route finding',
        data: 'Data using on this site',
        author: 'Author',
      },
    },
    ja: {
      message: {
        header: '日本語版Wikipediaで6回リンク辿ればいけるか調べる',
        search: '検索',
        searchFrom: 'から',
        searchTo: 'へのルートを',
        buttonRandom: 'ページをランダムに選ぶ',
        searching: '「{wordFrom}」から「{wordTo}」へのルートを検索中',
        searchResult: '「{wordFrom}」から「{wordTo}」へのリンクの検索結果 (実行時間 {second} sec)',
        searchInReverse: '「{wordTo}」から「{wordFrom}」を検索する',
        tweet: '結果をTwitterにつぶやく (別ウィンドウで開きます)',
        tweetFind: 'https://twitter.com/home?status=「{wordFrom}」から「{wordTo}」へはWikipediaで{length}リンクで行けるよ！ {link} {hashTag}',
        tweetNotFound: 'https://twitter.com/home?status=「{wordFrom}」から「{wordTo}」へはWikipediaで6回のリンクじゃいけないみたい… {link} {hashTag}',
        wikipediaUrl: 'https://ja.wikipedia.org/wiki/{word}',
      },
      error: {
        notFoundFrom: '"{from}"というページがないみたい',
        notFoundTo: '"{to}"というページがないみたい',
        notFoundRoute: '6回のリンクじゃ見つからなかった…ごめんね！',
      },
      about: {
        header: 'PediaRouteについて',
        sixDegree: '六次の隔たり',
        sixDegreeUrl: 'https://ja.wikipedia.org/wiki/%E5%85%AD%E6%AC%A1%E3%81%AE%E9%9A%94%E3%81%9F%E3%82%8A',
        first:
          'Wikipedia上には日本語版だけで数百万ものページがありますが、' +
          'その中のどんな2つのページを選んでも最大6回のページ内のリンクをたどれば到達できる（らしい）というのが2011年8月ごろにネットで話題になりました。',
        second: 'のWikipedia版ですね。',
        third:
          'このサイトはWikipedia内の全ページのタイトルおよびページ間のリンクを使い、スタートとゴールのページの双方向からリンクをたどれるか検索を行い、最短で見つかったリンクの経路を表示します。',
        fourth: '詳しいアルゴリズムに知りたい方は',
        sourceCode: 'ソースコード',
        fifth: 'を参照してください。',
        sixth: 'このサイトでは',
        database: 'Wikipediaのデータベース',
        databaseUrl:
          'https://ja.wikipedia.org/wiki/Wikipedia:%E3%83%87%E3%83%BC%E3%82%BF%E3%83%99%E3%83%BC%E3%82%B9%E3%83%80%E3%82%A6%E3%83%B3%E3%83%AD%E3%83%BC%E3%83%89',
        seventh: 'からページ情報とページ間のリンク情報を取得し、利用しています。',
        eighth: 'このサイトの情報は',
        commons: 'クリエイティブ・コモンズ 表示-継承 3.0 非移植ライセンス',
        commonsUrl: 'https://creativecommons.org/licenses/by-sa/3.0/deed.ja',
        nineth: 'の下で利用可能です。',
        findingRoute: '経路検索について',
        data: '使用しているデータについて',
        author: '作者について',
      },
    },
  },
});
