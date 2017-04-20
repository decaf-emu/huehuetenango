<template>
  <div class="uk-overlay uk-overlay-primary uk-position-cover">
    <div class="uk-position-center" uk-spinner></div>
  </div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import queryString from 'query-string';

export default {
  computed: {
    ...mapGetters(['isLoggedIn']),
  },

  methods: {
    ...mapActions(['processAuth']),
  },

  watch: {
    isLoggedIn(isLoggedIn) {
      if (isLoggedIn) {
        this.$router.push({ name: 'titles' });
      }
    },
  },

  beforeMount() {
    const { state, code } = queryString.parse(window.location.search);
    this.processAuth({ state, code });
  },
};
</script>

<style>
</style>
