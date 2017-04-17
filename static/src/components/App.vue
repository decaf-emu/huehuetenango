<template>
  <div id="app">
    <TopMenu />
    <SideMenu />
    <router-view class="view"></router-view>
  </div>
</template>

<script>
import TopMenu from './TopMenu.vue';
import SideMenu from './SideMenu.vue';

export default {
  components: { TopMenu, SideMenu },
};
</script>

<style>
body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif,
    "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
  margin: 0;
}
</style>
