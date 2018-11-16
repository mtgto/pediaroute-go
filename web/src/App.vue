<template>
  <div id="app">
    <router-view/>
    <footer>
      <nav>
        <router-link to="/">Home</router-link> |
        <router-link to="/about">About</router-link>
      </nav>
    </footer>
  </div>
</template>

<style>
* {
    margin : 0;
    padding : 0;
}

body {
    width: 980px;
    margin: 0 auto;
}

header {
    float: left;
    width: 800px;
}

header > h1{
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

article fieldset{
    border: none;
}
article input[type=text]{
    margin: 5px 0;
    width: 400px;
    height: 44px;
    font-size: 22pt;
}
article input[type=submit]{
    margin: 5px;
    float: right;
    width: 200px;
    height: 30px;
    border-radius: 5px;
}
article input[type=button]{
    margin: 5px;
    float: right;
    width: 200px;
    height: 30px;
    border-radius: 5px;
}

article ol{
    width: 800px;
    text-align: center;
    list-style-position: inside;
}

aside ul{
    width: 800px;
    text-align: center;
    list-style: none;
}

article ol li{
    font-size: 18pt;
    padding-bottom: 32px;
    background: url(../public/arrow.png) no-repeat center bottom;
}

article ol li:last-child{
    background: none;
}

aside {
    float: left;
    width: 800px;
}

aside#history h2{
    font-size: 14pt;
}
aside#history dl{
    margin: 10px;
}
aside#history dt, aside#history dd{
    margin: 5px;
}
aside#history dd{
    margin-left: 5px;
    text-indent: 10px;
}

section h2, aside h2 {
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
</style>
