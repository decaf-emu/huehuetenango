<template>
  <div id="offcanvas-nav-primary" uk-offcanvas="overlay: true">
    <div class="uk-offcanvas-bar uk-flex uk-flex-column">
      <ul class="uk-nav uk-nav-primary uk-nav-center uk-margin-auto-vertical">
        <li class="uk-nav-header">Huehuetenango</li>
        <li>
          <router-link :to="{ name: 'titles' }">Titles</router-link>
        </li>
        <li>
          <router-link :to="{ name: 'title',  params: { titleId: '000500101000400A' }}">System Library</router-link>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
</script>

<style>
</style>
