<template>
  <div id="app">
    <TopMenu />
    <SideMenu />
    <router-view class="view"></router-view>
    <SearchResults />
  </div>
</template>

<script>
import TopMenu from './TopMenu.vue';
import SideMenu from './SideMenu.vue';
import SearchResults from './search/SearchResults.vue';

export default {
  components: { TopMenu, SideMenu, SearchResults },
};
</script>

<style>
body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif,
    "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
  margin: 0;
}
</style>
