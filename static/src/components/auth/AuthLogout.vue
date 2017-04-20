<template>
  <div class="uk-container uk-container-expand">
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';

export default {
  computed: {
    ...mapGetters(['isLoggedIn']),
  },

  watch: {
    ...mapActions(['logout']),

    isLoggedIn(isLoggedIn) {
      if (!isLoggedIn) {
        this.$router.push({ name: 'titles' });
      }
    },
  },

  beforeMount() {
    this.logout();
  },
};
</script>

<style>
</style>
