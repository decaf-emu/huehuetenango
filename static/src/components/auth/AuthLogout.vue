<template>
  <div class="uk-container uk-container-expand">
  </div>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  beforeMount() {
    this.$store.dispatch('logout');
  },
  computed: mapGetters({
    isLoggedIn: 'isLoggedIn',
  }),
  watch: {
    isLoggedIn(isLoggedIn) {
      console.log(isLoggedIn);
      if (!isLoggedIn) {
        this.$router.push({ name: 'titles' });
      }
    },
  },
};
</script>

<style>
</style>
