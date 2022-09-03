<template>
  <div id="app">
    <router-view></router-view>
    <footer>
      <nav>
        <router-link to="/">Home</router-link>&nbsp;|
        <router-link to="/about">About</router-link>
      </nav>
    </footer>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
  // type inference enabled
});
</script>

<style>
* {
  margin: 0;
  padding: 0;
}

body {
  width: 980px;
  margin: 0 auto;
}

div.center {
  text-align: center;
}

a {
  display: inline-block;
  padding: 0.1em;
  color: #5fa990;
}

header {
  float: left;
  width: 800px;
}

header > h1 {
  font-size: 24pt;
}

header > h1 > span {
  font-size: 1.2em;
  color: #8fd9c0;
}

header > p {
  font-size: 11pt;
}

article {
  margin: 10px 0;
  float: left;
  width: 800px;
}

article fieldset {
  border: none;
}

article button.random {
  width: 40px;
  height: 40px;
  margin: 5px;
}

article input[type='text'] {
  margin: 5px 0;
  width: 400px;
  height: 44px;
  font-size: 22pt;
}

article input.submit {
  font-size: 1.2em;
  margin: 15px;
  width: 200px;
  height: 40px;
  border-radius: 5px;
}

article ol {
  width: 800px;
  text-align: center;
  list-style-position: inside;
}

aside ul {
  width: 800px;
  text-align: center;
  list-style: none;
}

article ol li:first-of-type {
  list-style-type: none;
}

article ol li {
  font-size: 18pt;
  padding-bottom: 32px;
  background: url(/arrow.png) no-repeat center bottom;
}

article ol li:last-child {
  background: none;
}

aside {
  float: left;
  width: 800px;
}

section h2,
aside h2 {
  font-size: 14pt;
  margin: 10px 0px;
}

footer {
  margin-top: 25px;
  padding-top: 25px;
  width: 800px;
  text-align: center;
  clear: both;
}

@media screen and (max-width: 480px) {
  body {
    width: auto;
    padding: 20px;
  }

  article input[type='text'] {
    padding: auto;
    width: auto;
  }

  header,
  footer,
  aside,
  article,
  article ol,
  aside ul {
    width: 100%;
  }
}
</style>
