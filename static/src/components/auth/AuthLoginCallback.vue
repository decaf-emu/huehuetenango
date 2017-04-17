<template>
  <div class="uk-overlay uk-overlay-primary uk-position-cover">
    <div class="uk-position-center" uk-spinner></div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import queryString from 'query-string';

export default {
  beforeMount() {
    const { state, code } = queryString.parse(window.location.search);
    this.$store.dispatch('processAuth', { state, code });
  },
  computed: mapGetters({
    isLoggedIn: 'isLoggedIn',
  }),
  watch: {
    isLoggedIn(isLoggedIn) {
      if (isLoggedIn) {
        this.$router.push({ name: 'titles' });
      }
    },
  },
};
</script>

<style>
</style>
